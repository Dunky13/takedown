// Package downloader transfers and verifies archives before publishing them.
package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/space"
	"github.com/Dunky13/takedown/internal/state"
	"github.com/Dunky13/takedown/internal/storage"
)

type LinkSource interface {
	RefreshLink(context.Context, string, int) (string, error)
}

type Options struct {
	Backend       storage.StorageBackend
	Destination   string
	State         *state.Store
	Guard         *space.Guard
	Bus           *events.Bus
	Client        *http.Client
	Links         LinkSource
	Concurrency   int
	PartByPart    bool
	CheckInterval int64
	MaxRetries    int
}

var (
	ErrIntegrity        = errors.New("archive integrity verification failed")
	ErrProtocol         = errors.New("invalid archive HTTP response")
	ErrIdentityChanged  = errors.New("remote archive identity changed; partial archive preserved")
	ErrReauthentication = errors.New("archive link requires authentication")
	ErrAlreadyRunning   = errors.New("download engine is already running")
)

type Engine struct {
	opts   Options
	read   storage.ReadBackend
	client *http.Client
	runMu  sync.Mutex
	// ioMu is a write barrier shared by every worker. Pauses close and checkpoint
	// all writers while holding it, before another worker can write any bytes.
	ioMu        sync.Mutex
	sessions    map[string]*session
	manualPause bool
	diskPause   bool
	pauseErr    error
	changed     chan struct{}
	spacePoll   time.Duration
	linkMu      sync.Mutex
	links       map[string]string
}

type session struct {
	checkpointErr   error
	record          state.Record
	name            string
	proveComplete   bool
	writer          io.WriteCloser
	reserved        bool
	sinceCheck      int64
	sinceSpaceCheck int64
	started         time.Time
	startOffset     int64
}

func New(opts Options) (*Engine, error) {
	if opts.Backend == nil || opts.State == nil || strings.TrimSpace(opts.Destination) == "" {
		return nil, errors.New("backend, destination, and state are required")
	}
	r, ok := opts.Backend.(storage.ReadBackend)
	if !ok {
		return nil, errors.New("destination must support reads for archive integrity verification")
	}
	if opts.Concurrency < 0 || opts.Concurrency > 128 || opts.CheckInterval < 0 || opts.MaxRetries < 0 || opts.MaxRetries > 20 {
		return nil, errors.New("invalid download concurrency, check interval, or retry limit")
	}
	if opts.Concurrency == 0 {
		opts.Concurrency = 2
	}
	if opts.PartByPart {
		opts.Concurrency = 1
	}
	if opts.CheckInterval == 0 {
		opts.CheckInterval = 50_000_000
	}
	if opts.CheckInterval > 50_000_000 {
		opts.CheckInterval = 50_000_000
	}
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 3
	}
	e := &Engine{opts: opts, read: r, sessions: make(map[string]*session), changed: make(chan struct{}), spacePoll: 2 * time.Second, links: make(map[string]string)}
	e.client = archiveClient(opts.Client)
	return e, nil
}

func (e *Engine) publish(kind events.Kind, key, message string, offset, total int64) {
	if e.opts.Bus != nil {
		e.opts.Bus.Publish(events.Event{Kind: kind, PartKey: key, Message: message, Bytes: offset, Total: total})
	}
}

func (e *Engine) publishPart(kind events.Kind, part model.Part, message string, offset int64) {
	if e.opts.Bus != nil {
		e.opts.Bus.Publish(events.Event{Kind: kind, PartKey: part.Key(), Filename: part.Filename, Message: message, Bytes: offset, Total: part.Size})
	}
}

func (e *Engine) notifyLocked() { close(e.changed); e.changed = make(chan struct{}) }

// Pause returns after all active writers have been closed and checkpointed.
// Download remains blocked until Resume or its context is canceled.
func (e *Engine) Pause() {
	e.ioMu.Lock()
	e.manualPause = true
	e.pauseErr = errors.Join(e.pauseErr, e.flushAllLocked())
	err := e.pauseErr
	e.notifyLocked()
	e.ioMu.Unlock()
	if err != nil {
		e.publish(events.Error, "", "Pause could not flush download checkpoints", 0, 0)
		return
	}
	e.publish(events.Status, "", "Paused; partial archives checkpointed", 0, 0)
}

// Resume wakes manual and disk-space waits. A disk pause still requires a
// successful free-space check; Resume cannot override the storage guard.
func (e *Engine) Resume() {
	e.ioMu.Lock()
	e.manualPause = false
	e.notifyLocked()
	e.ioMu.Unlock()
}

func (e *Engine) lowSpaceLocked() error {
	first := !e.diskPause
	e.diskPause = true
	if err := e.flushAllLocked(); err != nil {
		e.pauseErr = errors.Join(e.pauseErr, err)
		return err
	}
	if first {
		e.publish(events.SpaceWarning, "", "Low destination space; all workers paused and checkpointed", 0, 0)
	}
	e.notifyLocked()
	return nil
}

func (e *Engine) checkReservationsLocked(ctx context.Context) error {
	if e.opts.Guard == nil {
		return nil
	}
	// Update every existing reservation even if an early check finds a stale
	// total too large. Otherwise a paused worker's old total can prevent recovery.
	for key, s := range e.sessions {
		if !s.reserved {
			continue
		}
		err := e.opts.Guard.Check(ctx, key, s.record.Part.Size-s.record.Offset)
		if err != nil && !errors.Is(err, space.ErrInsufficientSpace) {
			return err
		}
	}
	for key, s := range e.sessions {
		remaining := s.record.Part.Size - s.record.Offset
		var err error
		if s.reserved {
			err = e.opts.Guard.Check(ctx, key, remaining)
		} else {
			err = e.opts.Guard.Reserve(ctx, key, remaining)
			if err == nil {
				s.reserved = true
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) waitReady(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		e.ioMu.Lock()
		if e.pauseErr != nil {
			err := e.pauseErr
			e.ioMu.Unlock()
			return err
		}
		if e.diskPause && !e.manualPause {
			err := e.checkReservationsLocked(ctx)
			if err == nil {
				e.diskPause = false
				for _, s := range e.sessions {
					s.sinceSpaceCheck = 0
				}
				e.publish(events.Status, "", "Destination space recovered; downloads resumed", 0, 0)
			} else if !errors.Is(err, space.ErrInsufficientSpace) {
				e.ioMu.Unlock()
				return err
			}
		}
		paused, change := e.manualPause || e.diskPause, e.changed
		e.ioMu.Unlock()
		if !paused {
			return nil
		}
		timer := time.NewTimer(e.spacePoll)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-change:
			timer.Stop()
		case <-timer.C:
		}
	}
}

// Download permits one run at a time and leaves backend/store ownership with
// the caller. Cancellation joins every worker before returning.
func (e *Engine) Download(ctx context.Context, parts []model.Part) error {
	if !e.runMu.TryLock() {
		return ErrAlreadyRunning
	}
	defer e.runMu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	seenKeys, seenNames := map[string]bool{}, map[string]bool{}
	for _, p := range parts {
		if err := p.Validate(); err != nil {
			return err
		}
		name := strings.ToLower(p.StoragePath())
		if seenKeys[p.Key()] || seenNames[name] {
			return errors.New("duplicate part or destination filename")
		}
		seenKeys[p.Key()], seenNames[name] = true, true
	}
	e.ioMu.Lock()
	e.pauseErr = nil
	e.ioMu.Unlock()
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	jobs := make(chan model.Part)
	var wg sync.WaitGroup
	var first error
	var errOnce sync.Once
	workers := min(e.opts.Concurrency, len(parts))
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				if err := e.downloadPart(runCtx, p); err != nil {
					errOnce.Do(func() { first = fmt.Errorf("archive %s: %w", p.Filename, err); cancel() })
					return
				}
			}
		}()
	}
feed:
	for _, p := range parts {
		select {
		case jobs <- p:
		case <-runCtx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()
	if first != nil {
		e.publish(events.Error, "", first.Error(), 0, 0)
		return first
	}
	return ctx.Err()
}

func (e *Engine) prepare(ctx context.Context, part model.Part) (*session, bool, error) {
	r, exists := e.opts.State.Get(part.Key())
	bound := exists && r.Destination == e.opts.Destination
	if !exists || r.Destination != e.opts.Destination {
		r = state.Record{Part: part, Destination: e.opts.Destination}
	} else if r.Part.Filename != part.Filename || r.Part.RelativePath != part.RelativePath || r.Part.Size != part.Size || !strings.EqualFold(r.Part.MD5, part.MD5) {
		return nil, false, ErrIdentityChanged
	}
	trustedLength := r.Offset
	r.Part = part
	if info, err := e.opts.Backend.Stat(ctx, part.StoragePath()); err == nil {
		if !info.Mode().IsRegular() || info.Size() != part.Size {
			return nil, false, fmt.Errorf("%w: existing archive has an unexpected size or type", ErrIntegrity)
		}
		if (!bound || trustedLength != part.Size || r.SHA256 == "") && part.MD5 == "" {
			// A filename and valid ZIP do not establish export identity. Prove
			// preexisting files against the source before registering completion.
			r.Offset, r.Complete = part.Size, false
			s := &session{record: r, name: part.StoragePath(), proveComplete: true}
			if err = e.proveExisting(ctx, s); err != nil {
				return nil, false, err
			}
			r.ETag, r.LastModified = s.record.ETag, s.record.LastModified
		}
		var digest string
		if digest, err = e.verifyArchiveDigest(ctx, part.StoragePath(), part, r.SHA256); err != nil {
			return nil, false, err
		}
		r.SHA256 = digest
		r.Offset, r.Complete = part.Size, true
		if err = e.opts.State.Put(r); err != nil {
			return nil, false, err
		}
		e.publishPart(events.PartCompleted, part, "Existing archive verified", r.Offset)
		return nil, true, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, false, err
	}
	r.Complete = false
	if info, err := e.opts.Backend.Stat(ctx, part.StoragePath()+".part"); err == nil {
		if !info.Mode().IsRegular() || info.Size() < 0 || info.Size() > part.Size {
			return nil, false, fmt.Errorf("%w: invalid partial archive size or type", ErrIntegrity)
		}
		// Filesystem length is authoritative after a crash. Every retained prefix
		// is checked against a strong validator or replayed response before append.
		r.Offset = info.Size()
	} else if errors.Is(err, os.ErrNotExist) {
		r.Offset, r.ETag, r.LastModified = 0, "", ""
	} else {
		return nil, false, err
	}
	if r.Offset != part.Size {
		r.SHA256 = ""
	}
	s := &session{record: r, name: part.StoragePath() + ".part", proveComplete: part.MD5 == "" && r.Offset == part.Size && (!bound || trustedLength != part.Size || r.SHA256 == ""), started: time.Now(), startOffset: r.Offset}
	if err := e.saveCheckpoint(s); err != nil {
		return nil, false, err
	}
	e.ioMu.Lock()
	e.sessions[part.Key()] = s
	var err error
	if e.opts.Guard != nil {
		err = e.opts.Guard.Reserve(ctx, part.Key(), part.Size-r.Offset)
		if err == nil {
			s.reserved = true
		}
	}
	if errors.Is(err, space.ErrInsufficientSpace) {
		err = e.lowSpaceLocked()
	}
	if err != nil {
		delete(e.sessions, part.Key())
		if e.opts.Guard != nil {
			e.opts.Guard.Release(part.Key())
		}
	}
	e.ioMu.Unlock()
	return s, false, err
}

func (e *Engine) downloadPart(ctx context.Context, part model.Part) (retErr error) {
	if err := e.waitReady(ctx); err != nil {
		return err
	}
	s, done, err := e.prepare(ctx, part)
	if err != nil || done {
		return err
	}
	defer func() {
		e.ioMu.Lock()
		retErr = errors.Join(retErr, e.closeCheckpointLocked(s))
		delete(e.sessions, part.Key())
		if e.opts.Guard != nil {
			e.opts.Guard.Release(part.Key())
		}
		e.ioMu.Unlock()
	}()
	link := part.URL
	if link == "" && (s.record.Offset < part.Size || s.proveComplete) {
		link, err = e.refresh(ctx, part, "")
		if err != nil {
			return err
		}
	}
	for attempt := 0; s.record.Offset < part.Size || s.proveComplete; attempt++ {
		if err = e.waitReady(ctx); err != nil {
			return err
		}
		err = e.transfer(ctx, s, link)
		e.ioMu.Lock()
		if err == nil {
			s.proveComplete = false
		}
		checkpointErr := e.closeCheckpointLocked(s)
		e.ioMu.Unlock()
		if checkpointErr != nil {
			return errors.Join(err, checkpointErr)
		}
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if attempt >= e.opts.MaxRetries {
			return err
		}
		if errors.Is(err, ErrReauthentication) {
			link, err = e.refresh(ctx, part, link)
			if err != nil {
				return err
			}
		} else if !retryable(err) {
			return err
		}
		if err = retryWait(ctx, attempt, err); err != nil {
			return err
		}
	}
	if err = e.waitReady(ctx); err != nil {
		return err
	}
	digest, err := e.verifyArchiveDigest(ctx, part.StoragePath()+".part", part, s.record.SHA256)
	if err != nil {
		return err
	}
	s.record.SHA256 = digest
	if err = e.waitReady(ctx); err != nil {
		return err
	}
	for {
		if err = e.waitReady(ctx); err != nil {
			return err
		}
		e.ioMu.Lock()
		if !e.manualPause && !e.diskPause {
			break
		}
		e.ioMu.Unlock()
	}
	defer e.ioMu.Unlock()
	if err = ctx.Err(); err != nil {
		return err
	}
	if err = e.opts.Backend.Rename(ctx, part.StoragePath()+".part", part.StoragePath()); err != nil {
		return err
	}
	s.record.Complete = true
	if err = e.opts.State.Put(s.record); err != nil {
		return err
	}
	e.publishPart(events.PartCompleted, part, "Archive verified and completed", s.record.Offset)
	return nil
}

func (e *Engine) closeCheckpointLocked(s *session) error {
	if s.checkpointErr != nil {
		return s.checkpointErr
	}
	var err error
	if s.writer != nil {
		if f, ok := s.writer.(interface{ Flush() error }); ok {
			err = f.Flush()
		}
		if f, ok := s.writer.(interface{ Sync() error }); ok {
			err = errors.Join(err, f.Sync())
		}
		err = errors.Join(err, s.writer.Close())
		s.writer = nil
	}
	if err != nil {
		s.checkpointErr = fmt.Errorf("flush archive checkpoint: %w", err)
		return s.checkpointErr
	}
	s.checkpointErr = e.saveCheckpoint(s)
	return s.checkpointErr
}

// Orphan complete files must not acquire a trusted full-length checkpoint
// until a source response has proved every byte. Keep that distinction across
// cancellation and authorization failures, without persisting download links.
func (e *Engine) saveCheckpoint(s *session) error {
	r := s.record
	if s.proveComplete {
		r.Offset, r.Complete, r.ETag, r.LastModified, r.SHA256 = 0, false, "", "", ""
	}
	return e.opts.State.Put(r)
}

func (e *Engine) proveExisting(ctx context.Context, s *session) error {
	link := s.record.Part.URL
	var err error
	if link == "" {
		link, err = e.refresh(ctx, s.record.Part, "")
		if err != nil {
			return err
		}
	}
	for attempt := 0; ; attempt++ {
		if err = e.waitReady(ctx); err != nil {
			return err
		}
		err = e.transfer(ctx, s, link)
		if err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if attempt >= e.opts.MaxRetries {
			return err
		}
		if errors.Is(err, ErrReauthentication) {
			link, err = e.refresh(ctx, s.record.Part, link)
			if err != nil {
				return err
			}
		} else if !retryable(err) {
			return err
		}
		if err = retryWait(ctx, attempt, err); err != nil {
			return err
		}
	}
}

func (e *Engine) flushAllLocked() error {
	var err error
	for _, s := range e.sessions {
		err = errors.Join(err, e.closeCheckpointLocked(s))
	}
	return err
}

func (e *Engine) write(ctx context.Context, s *session, data []byte) error {
	for {
		if err := e.waitReady(ctx); err != nil {
			return err
		}
		e.ioMu.Lock()
		if e.manualPause || e.diskPause {
			e.ioMu.Unlock()
			continue
		}
		if err := ctx.Err(); err != nil {
			e.ioMu.Unlock()
			return err
		}
		if e.opts.Guard != nil && s.sinceSpaceCheck >= e.opts.CheckInterval {
			// Other workers' reservations retain their previous (larger)
			// remaining size. That overestimates demand between checks and
			// prevents overcommit without a StatVFS/SMB request on every chunk.
			err := e.opts.Guard.Check(ctx, s.record.Part.Key(), s.record.Part.Size-s.record.Offset)
			if errors.Is(err, space.ErrInsufficientSpace) {
				err = e.lowSpaceLocked()
				e.ioMu.Unlock()
				if err != nil {
					return err
				}
				continue
			}
			if err != nil {
				e.ioMu.Unlock()
				return err
			}
			s.sinceSpaceCheck = 0
		}
		if s.writer == nil {
			w, err := e.opts.Backend.OpenWriteResume(ctx, s.record.Part.StoragePath()+".part", s.record.Offset)
			if err != nil {
				e.ioMu.Unlock()
				return err
			}
			s.writer = w
		}
		chunk := data
		if e.opts.Guard != nil {
			chunk = data[:min(int64(len(data)), e.opts.CheckInterval-s.sinceSpaceCheck)]
		}
		n, err := s.writer.Write(chunk)
		if n < 0 || n > len(chunk) {
			e.ioMu.Unlock()
			return errors.New("destination returned an invalid write count")
		}
		s.record.Offset += int64(n)
		// Exact local length alone does not prove a clean HTTP EOF. Keep a
		// complete file untrusted until the response has passed the excess check.
		if s.record.Offset == s.record.Part.Size {
			s.proveComplete = true
		}
		s.sinceCheck += int64(n)
		s.sinceSpaceCheck += int64(n)
		if errors.Is(err, syscall.ENOSPC) && e.opts.Guard != nil {
			pauseErr := e.lowSpaceLocked()
			e.ioMu.Unlock()
			if pauseErr != nil {
				return pauseErr
			}
			data = data[n:]
			if len(data) == 0 {
				return e.waitReady(ctx)
			}
			continue
		}
		if err == nil && n != len(chunk) {
			err = io.ErrShortWrite
		}
		if err == nil && s.sinceCheck >= e.opts.CheckInterval {
			err = e.closeCheckpointLocked(s)
			s.sinceCheck = 0
		}
		offset := s.record.Offset
		e.ioMu.Unlock()
		if e.opts.Bus != nil {
			seconds := time.Since(s.started).Seconds()
			e.opts.Bus.Publish(events.Event{Kind: events.Progress, PartKey: s.record.Part.Key(), Filename: s.record.Part.Filename, Bytes: offset, Total: s.record.Part.Size, BytesPerSecond: float64(offset-s.startOffset) / seconds})
		}
		if err != nil {
			return err
		}
		data = data[n:]
		if len(data) == 0 {
			return nil
		}
	}
}
