package app

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/Dunky13/takedown/internal/downloader"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/source"
	"github.com/Dunky13/takedown/internal/storage"
)

type WatchOptions struct {
	DownloadOptions DownloadOptions
	Interval        time.Duration
	ExportIDs       []string
}

type watchedArchive struct {
	part     model.Part
	modified time.Time
}

// Watch discovers and downloads new archives until canceled. It refreshes OAuth
// tokens through the account client, but never creates exports or opens a login
// window. Revoked authorization returns an actionable error to the operator.
func (a *App) Watch(ctx context.Context, opts WatchOptions) error {
	interval := opts.Interval
	if interval == 0 {
		interval = a.Settings.PollInterval()
	}
	minimum := time.Minute
	if a.Settings.EffectiveSource() == "portability" {
		minimum = 5 * time.Minute
	}
	if interval < minimum || interval > 24*time.Hour {
		return errors.New("watch interval must be at least one minute for Drive or five minutes for Data Portability, and at most 24 hours")
	}
	if err := a.beginJob(); err != nil {
		return err
	}
	defer a.endJob()
	known := make(map[string]watchedArchive)
	failures := 0
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		delay := interval
		if err := a.watchPass(ctx, opts, known); err != nil {
			if source.IsJobFailure(err) {
				// This job's immutable data is held, but a different recorded job
				// may become ready later. Retain the normal polling cadence.
				failures = 0
				a.Bus.Publish(events.Event{Kind: events.Status, Message: "An export job is held for review; other jobs remain eligible and watching will continue"})
			} else {
				if !source.IsTemporary(err) && !downloader.IsTemporary(err) {
					return err
				}
				failures++
				delay = max(delay, min(time.Minute*time.Duration(1<<min(failures, 6)), time.Hour), source.RetryDelay(err), downloader.RetryDelay(err))
				a.Bus.Publish(events.Event{Kind: events.Error, Message: "Temporary Google or network failure; partial files retained; discovery will retry after backoff"})
			}
		} else {
			failures = 0
			a.Bus.Publish(events.Event{Kind: events.Status, Message: "Watching for new ready archives; no export will be created automatically"})
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func (a *App) watchPass(ctx context.Context, opts WatchOptions, known map[string]watchedArchive) (retErr error) {
	exports, err := a.Account.List(ctx)
	if err != nil {
		return err
	}
	wanted := make(map[string]bool)
	for _, id := range opts.ExportIDs {
		wanted[id] = true
	}
	var parts []model.Part
	keys := map[string]bool{}
	for _, export := range exports {
		if len(wanted) > 0 && !wanted[export.ID] {
			continue
		}
		for _, p := range export.Parts {
			if err = p.Validate(); err != nil {
				return err
			}
			if keys[p.Key()] {
				return errors.New("source returned duplicate archive identities")
			}
			keys[p.Key()] = true
			parts = append(parts, p)
		}
	}
	if len(parts) == 0 {
		return nil
	}
	destination := opts.DownloadOptions.Destination
	if destination == "" {
		destination = a.Settings.Destination
	}
	opts.DownloadOptions.StoragePassword, err = a.storagePassword(destination, opts.DownloadOptions.StoragePassword)
	if err != nil {
		return err
	}
	parts, err = a.previewLayout(parts, destination)
	if err != nil {
		return err
	}
	backend, err := storage.Open(ctx, destination, storage.Options{Password: opts.DownloadOptions.StoragePassword, PrivateKeyPath: a.Settings.PrivateKeyPath, KnownHostsPath: a.Settings.KnownHostsPath, AllowUnknownSpace: a.Settings.AllowUnknownSpace})
	if err != nil {
		return err
	}
	var pending []model.Part
	for _, part := range parts {
		old, ok := known[part.Key()]
		if !ok || old.part.Filename != part.Filename || old.part.Size != part.Size || !strings.EqualFold(old.part.MD5, part.MD5) {
			pending = append(pending, part)
			continue
		}
		info, statErr := backend.Stat(ctx, part.StoragePath())
		if statErr == nil && info.Mode().IsRegular() && info.Size() == part.Size && info.ModTime().Equal(old.modified) {
			continue
		}
		if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return errors.Join(statErr, backend.Close())
		}
		pending = append(pending, part)
	}
	if err = backend.Close(); err != nil {
		return err
	}
	if len(pending) == 0 {
		return nil
	}
	if err = a.download(ctx, pending, opts.DownloadOptions); err != nil {
		return err
	}
	backend, err = storage.Open(ctx, destination, storage.Options{Password: opts.DownloadOptions.StoragePassword, PrivateKeyPath: a.Settings.PrivateKeyPath, KnownHostsPath: a.Settings.KnownHostsPath, AllowUnknownSpace: a.Settings.AllowUnknownSpace})
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, backend.Close()) }()
	for _, part := range pending {
		info, statErr := backend.Stat(ctx, part.StoragePath())
		if statErr != nil {
			return statErr
		}
		known[part.Key()] = watchedArchive{part: part, modified: info.ModTime()}
	}
	return nil
}
