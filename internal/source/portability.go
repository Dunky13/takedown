package source

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Dunky13/takedown/internal/auth"
	"github.com/Dunky13/takedown/internal/model"
	"google.golang.org/api/dataportability/v1"
	"google.golang.org/api/googleapi"
)

const (
	portabilityPollInterval = 5 * time.Minute
	portabilityLinkLife     = 6 * time.Hour
	portabilityDataLife     = 14 * 24 * time.Hour
	// The intent precedes a bounded 30-second POST. Allow five extra minutes
	// beyond the fourteen-day policy window, rather than expiring at send time.
	portabilityIntentExpiry = portabilityDataLife + 5*time.Minute
)

// The complete send window, including a slow durable write, must fit inside
// the five-minute margin used by explicit expired-intent recovery.
func (m *Manager) creationContext(ctx context.Context, intentAt time.Time) (context.Context, context.CancelFunc) {
	remaining := intentAt.Add(portabilityIntentExpiry - portabilityDataLife).Sub(m.now())
	if remaining > 30*time.Second {
		remaining = 30 * time.Second
	}
	return context.WithTimeout(ctx, remaining)
}

type cachedLinks struct {
	URLs []string
	At   time.Time
	// List hands out every URL. A legacy RefreshLink call hands out one part;
	// siblings may claim the rest of that generation without another API poll.
	delivered []bool
	rejected  []bool
}

func (c cachedLinks) fresh(now time.Time) bool {
	return len(c.URLs) > 0 && now.Before(c.At.Add(portabilityLinkLife-5*time.Minute))
}
func (c cachedLinks) partRejected(index int) bool {
	return index < len(c.rejected) && c.rejected[index]
}
func (c cachedLinks) allAvailable(now time.Time) bool {
	if !c.fresh(now) {
		return false
	}
	for _, rejected := range c.rejected {
		if rejected {
			return false
		}
	}
	return true
}

func (m *Manager) deliverPortability(jobID string, index int) string {
	c := m.links[jobID]
	if len(c.delivered) != len(c.URLs) {
		c.delivered = make([]bool, len(c.URLs))
	}
	c.delivered[index] = true
	m.links[jobID] = c
	return c.URLs[index]
}

func (m *Manager) rejectPortability(jobID string, index int) {
	c := m.links[jobID]
	if len(c.rejected) != len(c.URLs) {
		c.rejected = make([]bool, len(c.URLs))
	}
	c.rejected[index] = true
	m.links[jobID] = c
}

func portabilityCooldown(j portabilityJob, now time.Time) error {
	until := j.LastPoll.Add(portabilityPollInterval)
	if j.RetryNotBefore.After(until) {
		until = j.RetryNotBefore
	}
	if !now.Before(until) {
		return nil
	}
	return &sourceError{message: "Data Portability archive refresh is cooling down; retry after the remaining polling or provider backoff interval", temporary: true, delay: until.Sub(now)}
}

func selectedJob(j portabilityJob, selected []string) bool {
	set := map[string]bool{}
	for _, r := range selected {
		set[r] = true
	}
	for _, r := range j.Resources {
		if !set[r] {
			return false
		}
	}
	return true
}
func overlapping(a, b []string) bool {
	set := map[string]bool{}
	for _, r := range a {
		set[r] = true
	}
	for _, r := range b {
		if set[r] {
			return true
		}
	}
	return false
}

func (m *Manager) createPortability(ctx context.Context) (string, error) {
	state, name, unlock, err := m.openState(ctx)
	if err != nil {
		return "", err
	}
	defer unlock()
	if state.Intent != nil {
		return "", errors.New("a previous Google mutation has an unknown outcome; reconcile its durable intent before another mutation; automatic retries are disabled")
	}
	if activeJobs(state) >= maxJobs {
		return "", errors.New("archive job history limit reached; archive the recorded job history before creating more exports")
	}
	if len(state.Jobs) >= maxDriveEntries {
		return "", errors.New("retained job safety receipt limit reached; history must not be discarded to bypass grant or cadence guards")
	}
	now := m.now().UTC()
	for _, j := range state.Jobs {
		if !overlapping(j.Resources, m.opts.Resources) {
			continue
		}
		if now.Before(j.CreatedAt.Add(24 * time.Hour)) {
			return "", errors.New("these Data Portability resources were requested less than 24 hours ago; wait before explicitly creating another export")
		}
		if j.AccessType == "ACCESS_TYPE_ONE_TIME" && j.Generation == state.Generation {
			return "", errors.New("these Data Portability resources used a one-time grant; do not repeat the export under that grant; obtain and verify a new grant and reconcile the retained job history first")
		}
	}
	// This evidence belongs to the original attempt, before its POST. A later
	// login or accessType.check cannot establish a lost request's grant type.
	access, checkErr := m.checkAccess(ctx)
	if checkErr != nil {
		return "", checkErr
	}
	accessType, checkErr := accessFor(access, m.opts.Resources)
	if checkErr != nil {
		return "", checkErr
	}
	if len(state.ExpiredIntents) > 0 {
		for _, intent := range state.ExpiredIntents {
			if overlapping(intent.Resources, m.opts.Resources) && !containsAll(access.TimeBasedResources, intent.Resources) {
				return "", errors.New("an expired unknown creation retains a one-time grant hold; recover an evidence-backed job ID before exporting these resources")
			}
		}
	}
	if err = ctx.Err(); err != nil {
		return "", err
	}
	// Start the policy/cadence clock after the preflight read, immediately
	// before recording the intent and sending the mutation.
	now = m.now().UTC()
	state.Intent = &creationIntent{Kind: "create", Resources: append([]string(nil), m.opts.Resources...), CreatedAt: now, AccessType: accessType}
	if err = m.persist(name, state); err != nil {
		return "", err
	}
	callCtx, cancel := m.creationContext(ctx, now)
	if callCtx.Err() != nil {
		cancel()
		return "", errors.New("creation intent persistence consumed the request margin; no mutation was sent; the durable intent is retained for explicit reconciliation")
	}
	// No retry of POST: the API does not provide an idempotency key or a job
	// enumeration endpoint to resolve a response lost after remote acceptance.
	response, err := m.portability.PortabilityArchive.Initiate(&dataportability.InitiatePortabilityArchiveRequest{Resources: append([]string(nil), m.opts.Resources...)}).Context(callCtx).Do()
	cancel()
	if err != nil {
		var g *googleapi.Error
		// A structured rejection establishes that this request did not start a
		// job. Network errors, 408, 5xx, or malformed success retain the hold.
		if errors.As(err, &g) && g.Code >= 400 && g.Code < 500 && g.Code != 408 {
			state.Intent = nil
			if saveErr := m.persist(name, state); saveErr != nil {
				return "", saveErr
			}
			safe := apiError(ctx, "create Data Portability export", err)
			return "", &sourceError{message: safe.Error(), authorization: errors.Is(safe, auth.ErrAuthorizationRequired)}
		}
		if ctx.Err() != nil {
			return "", fmt.Errorf("%w; export creation outcome is unknown and the persisted intent blocks duplicate requests", ctx.Err())
		}
		return "", errors.New("Google export creation outcome is unknown; the persisted intent blocks duplicate requests; reconcile it with Google before trying again")
	}
	if !validDriveID.MatchString(response.ArchiveJobId) || (response.AccessType != "ACCESS_TYPE_ONE_TIME" && response.AccessType != "ACCESS_TYPE_TIME_BASED") {
		return "", errors.New("Google returned incomplete export creation metadata; the persisted intent blocks duplicate requests until reconciled")
	}
	for _, j := range state.Jobs {
		if j.ID == response.ArchiveJobId {
			return "", errors.New("Google returned a duplicate archive job ID; the persisted intent is retained for reconciliation")
		}
	}
	state.Jobs = append(state.Jobs, portabilityJob{ID: response.ArchiveJobId, Resources: append([]string(nil), m.opts.Resources...), CreatedAt: now, AccessType: response.AccessType, LastPoll: now, State: "IN_PROGRESS", Generation: state.Generation})
	state.Intent = nil
	state.NeedsConsent = false
	if err = m.persist(name, state); err != nil {
		return "", fmt.Errorf("Google accepted archive job %s but its receipt could not be saved; reconcile the retained creation intent before another export", response.ArchiveJobId)
	}
	m.emit("Data Portability export job recorded; listing will poll this job without creating further exports")
	return "Data Portability archive job " + response.ArchiveJobId + " created and recorded. Check exports after at least five minutes; processing may take days.", nil
}

func (m *Manager) listPortability(ctx context.Context) ([]model.Export, error) {
	state, name, unlock, err := m.openState(ctx)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if state.Intent != nil {
		m.emit("A previous Google mutation has an unknown outcome; existing recorded jobs remain readable, but further mutations are held")
	}
	exports := []model.Export{}
	var scopedFailure error
	for i := range state.Jobs {
		j := &state.Jobs[i]
		if j.Archived || j.Generation != state.Generation || !selectedJob(*j, m.opts.Resources) {
			continue
		}
		if !m.now().Before(j.CreatedAt.Add(portabilityDataLife)) {
			if j.State != "EXPIRED" && j.State != "COMPLETE" {
				j.State = "EXPIRED"
				if err = m.persist(name, state); err != nil {
					return nil, err
				}
			}
			m.emit("A recorded Data Portability job has reached the fourteen-day data lifetime")
			continue
		}
		if j.State == "FAILED" || j.State == "CANCELLED" || j.State == "EXPIRED" {
			m.emit("A recorded Data Portability job is terminal; other ready jobs remain readable; use explicit job management for recovery")
			continue
		}
		cached := m.links[j.ID]
		// Completed URLs only live in memory. A restart observes persisted
		// LastPoll before fetching new signatures, preserving the API cadence.
		cacheValid := cached.allAvailable(m.now())
		if !cacheValid && !m.now().Before(j.LastPoll.Add(portabilityPollInterval)) && !m.now().Before(j.RetryNotBefore) {
			if err = m.pollPortability(ctx, state, name, i); err != nil {
				if isJobFailure(err) {
					j.ReadIssue = true
					if saveErr := m.persist(name, state); saveErr != nil {
						return nil, saveErr
					}
					scopedFailure = err
					m.emit("One Data Portability job has unusable metadata; its parts remain held while unrelated ready jobs continue")
					continue
				}
				return nil, err
			}
			cached = m.links[j.ID]
			cacheValid = cached.allAvailable(m.now())
		}
		if j.State != "COMPLETE" || !cacheValid {
			m.emit("Data Portability job is pending or within the five-minute polling cooldown")
			continue
		}
		if len(cached.URLs) != len(j.Parts) {
			return nil, errors.New("recorded Data Portability archive part count differs from current links")
		}
		cached.delivered = make([]bool, len(cached.URLs))
		for i := range cached.delivered {
			cached.delivered[i] = true
		}
		m.links[j.ID] = cached
		exports = append(exports, jobExport(*j, cached.URLs))
	}
	if len(exports) == 0 && scopedFailure != nil {
		return nil, scopedFailure
	}
	m.emit("Listed only Data Portability jobs explicitly created and recorded by this app; pre-existing web Takeout exports are not enumerated")
	return exports, nil
}

func jobExport(j portabilityJob, urls []string) model.Export {
	id := "portability-" + j.ID
	export := model.Export{ID: id, CreatedAt: j.CreatedAt, Parts: make([]model.Part, 0, len(j.Parts))}
	for i, p := range j.Parts {
		export.Parts = append(export.Parts, model.Part{ExportID: id, Index: i + 1, Filename: p.Filename, Size: p.Size, MD5: p.MD5, URL: urls[i]})
	}
	return export
}

func (m *Manager) pollPortability(ctx context.Context, state *portabilityState, name string, index int) error {
	j := &state.Jobs[index]
	if err := portabilityCooldown(*j, m.now()); err != nil {
		return err
	}
	if !m.now().Before(j.CreatedAt.Add(portabilityDataLife)) {
		return errors.New("Data Portability job data reached its fourteen-day lifetime; explicitly create an export under a valid grant")
	}
	j.LastPoll = m.now().UTC()
	j.ReadIssue = false
	j.RetryNotBefore = time.Time{}
	previous := m.links[j.ID]
	delete(m.links, j.ID)
	// Persist even failed attempts so errors/restarts cannot turn a watch loop
	// into polling faster than Google's documented five-minute minimum.
	if err := m.persist(name, state); err != nil {
		return err
	}
	// One API attempt per cadence slot. Retrying immediately would violate the
	// minimum polling cadence. Retry-After is carried to the actionable error.
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	response, err := m.portability.ArchiveJobs.GetPortabilityArchiveState("archiveJobs/" + j.ID + "/portabilityArchiveState").Context(callCtx).Do()
	cancel()
	if err != nil {
		safe := apiError(ctx, "read Data Portability archive state", err)
		var g *googleapi.Error
		if ctx.Err() == nil && errors.As(err, &g) && (g.Code == 404 || g.Code == 410) {
			// Even an otherwise job-local missing-resource response may carry
			// provider backoff. Preserve it as a global safety failure.
			if delay := retryAfter(g.Header.Get("Retry-After"), time.Now()); delay > 0 {
				j.RetryNotBefore = m.now().Add(delay).UTC()
				if saveErr := m.persist(name, state); saveErr != nil {
					return saveErr
				}
				return safe
			}
			return &jobReadError{safe}
		}
		return m.portabilityReadError(state, name, j, safe)
	}
	return m.acceptPortability(ctx, state, name, index, response, previous)
}

func (m *Manager) acceptPortability(ctx context.Context, state *portabilityState, name string, index int, response *dataportability.PortabilityArchiveState, previous cachedLinks) error {
	j := &state.Jobs[index]
	var err error
	if response.Name != "" && response.Name != "archiveJobs/"+j.ID+"/portabilityArchiveState" {
		return jobFailure("Google returned state for a different archive job")
	}
	switch response.State {
	case "IN_PROGRESS":
		if j.State == "COMPLETE" || len(j.Parts) > 0 {
			return jobFailure("completed Data Portability job unexpectedly returned to in-progress state")
		}
		if len(response.Urls) > 0 {
			return jobFailure("in-progress Data Portability job unexpectedly included archive URLs")
		}
		j.State = "IN_PROGRESS"
		return m.persist(name, state)
	case "FAILED", "CANCELLED":
		j.State = response.State
		if err = m.persist(name, state); err != nil {
			return err
		}
		return nil
	case "COMPLETE":
	default:
		return jobFailure("Google returned an unknown Data Portability job state")
	}
	if len(response.Urls) == 0 || len(response.Urls) > maxDriveEntries {
		return jobFailure("completed Data Portability job has missing URLs or exceeds the archive part limit")
	}
	parts := make([]archiveMetadata, 0, len(response.Urls))
	objects, names := map[string]bool{}, map[string]bool{}
	for _, raw := range response.Urls {
		meta, err := m.probeSigned(ctx, raw)
		if err != nil {
			var safety *sourceError
			if ctx.Err() == nil && !errors.As(err, &safety) {
				return &jobReadError{err}
			}
			return m.portabilityReadError(state, name, j, err)
		}
		folded := strings.ToLower(meta.Filename)
		if objects[meta.Object] || names[folded] {
			return jobFailure("Data Portability job returned duplicate objects or ambiguous archive filenames")
		}
		objects[meta.Object] = true
		names[folded] = true
		parts = append(parts, meta)
	}
	// The API array defines one-based part order. Re-signing must preserve both
	// count and each object's order/name/exact size, including after a restart.
	if len(j.Parts) > 0 && !reflect.DeepEqual(j.Parts, parts) {
		return jobFailure("Data Portability archive part identity, order, name, or size changed; stale identities cannot receive new bytes")
	}
	j.State = "COMPLETE"
	j.Parts = parts
	if err = m.persist(name, state); err != nil {
		return err
	}
	// Metadata probing can take time. Count the six-hour signature lifetime
	// from the request slot, never from completion of all metadata probes.
	cached := cachedLinks{URLs: append([]string(nil), response.Urls...), At: j.LastPoll, delivered: make([]bool, len(response.Urls)), rejected: make([]bool, len(response.Urls))}
	for i, raw := range cached.URLs {
		if i < len(previous.URLs) && previous.URLs[i] == raw && previous.partRejected(i) {
			cached.rejected[i] = true
		}
	}
	m.links[j.ID] = cached
	return nil
}

func (m *Manager) portabilityReadError(state *portabilityState, name string, j *portabilityJob, err error) error {
	if delay := RetryDelay(err); IsTemporary(err) && delay > 0 {
		j.RetryNotBefore = m.now().Add(delay).UTC()
		if saveErr := m.persist(name, state); saveErr != nil {
			return saveErr
		}
	}
	return err
}

func (m *Manager) refreshPortability(ctx context.Context, id string, index int, rejectedURL string) (string, error) {
	jobID := strings.TrimPrefix(id, "portability-")
	if jobID == id || !validDriveID.MatchString(jobID) {
		return "", errors.New("invalid Data Portability export identity")
	}
	state, name, unlock, err := m.openState(ctx)
	if err != nil {
		return "", err
	}
	defer unlock()
	for i := range state.Jobs {
		j := &state.Jobs[i]
		if j.ID != jobID {
			continue
		}
		if j.Archived || j.Generation != state.Generation {
			return "", errors.New("archive job is archived or belongs to an earlier authorization generation")
		}
		if !selectedJob(*j, m.opts.Resources) {
			return "", errors.New("recorded job includes resource groups outside the selected OAuth scopes")
		}
		if index > len(j.Parts) || j.State != "COMPLETE" {
			return "", errors.New("requested Data Portability archive part has no recorded immutable metadata; list completed jobs first")
		}
		if !m.now().Before(j.CreatedAt.Add(portabilityDataLife)) {
			return "", errors.New("Data Portability job data reached its fourteen-day lifetime; explicitly create an export under a valid grant")
		}
		part := index - 1
		cached := m.links[jobID]
		if len(cached.URLs) == len(j.Parts) {
			if rejectedURL != "" && cached.URLs[part] == rejectedURL {
				m.rejectPortability(jobID, part)
				cached = m.links[jobID]
			}
			if cached.fresh(m.now()) && !cached.partRejected(part) {
				legacyUnclaimed := rejectedURL == "" && (part >= len(cached.delivered) || !cached.delivered[part])
				if legacyUnclaimed || rejectedURL != "" && cached.URLs[part] != rejectedURL {
					return m.deliverPortability(jobID, part), nil
				}
			}
		}
		if err = m.pollPortability(ctx, state, name, i); err != nil {
			return "", err
		}
		links := m.links[jobID]
		if len(links.URLs) != len(j.Parts) {
			return "", errors.New("Data Portability job did not return the same archive part count")
		}
		if rejectedURL != "" && links.URLs[part] == rejectedURL {
			m.rejectPortability(jobID, part)
			links = m.links[jobID]
		}
		if links.partRejected(part) {
			if err := portabilityCooldown(*j, m.now()); err != nil {
				return "", err
			}
			return "", &sourceError{message: "Google returned a rejected archive signature again; retry the refresh later", temporary: true}
		}
		return m.deliverPortability(jobID, part), nil
	}
	return "", errors.New("Data Portability export is not a job recorded by this app; pre-existing web Takeout exports are unsupported")
}
