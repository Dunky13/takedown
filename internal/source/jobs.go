package source

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/Dunky13/takedown/internal/auth"
	downloadstate "github.com/Dunky13/takedown/internal/state"
	"google.golang.org/api/dataportability/v1"
	"google.golang.org/api/googleapi"
)

// Job is a local, credential-free receipt. State is IN_PROGRESS, COMPLETE,
// FAILED, CANCELLED, or EXPIRED. Archived retains the terminal state and all
// grant/cadence guards. No signed download URLs are exposed here.
type Job struct {
	ID, State, AccessType, Message                 string
	Resources                                      []string
	CreatedAt, LastPoll, ExpiresAt, RetryNotBefore time.Time
	RetryCount                                     int
	Archived                                       bool
}

type JobSnapshot struct {
	Jobs   []Job
	Intent *PendingAction
}

// PendingAction describes a durable mutation whose outcome is unresolved.
// Kind is create, retry, cancel, or reset. Recovery is always explicit.
type PendingAction struct {
	Kind, JobID string
	Resources   []string
	CreatedAt   time.Time
	Message     string
}

type AccessInfo struct {
	OneTimeResources, TimeBasedResources []string
}

type resetConfirmationKey struct{}

// WithResetConfirmation records the caller's explicit confirmation that reset
// revokes consent and permits Google to garbage-collect previous archives.
// Call only after the user confirms. It cannot bypass the durable completion
// and download checks in ResetAuthorization.
func WithResetConfirmation(ctx context.Context) context.Context {
	return context.WithValue(ctx, resetConfirmationKey{}, true)
}

// Jobs reads the local journal without logging in or contacting Google.
func (m *Manager) Jobs(ctx context.Context) (JobSnapshot, error) {
	_, state, _, release, err := m.jobState(ctx, false)
	if err != nil {
		return JobSnapshot{}, err
	}
	defer release()
	return m.jobSnapshot(state), nil
}

// RefreshJobs reads at most 100 due jobs per call. Each job has a durable
// five-minute status cadence. This never retries a mutation or resolves intent.
// Oldest LastPoll first, with ID ties, gives fair progress across restarts even
// when every call is spaced far enough apart for the previous batch to be due.
func (m *Manager) RefreshJobs(ctx context.Context) (JobSnapshot, error) {
	ctx, state, name, release, err := m.jobState(ctx, true)
	if err != nil {
		return JobSnapshot{}, err
	}
	defer release()
	due := make([]int, 0, len(state.Jobs))
	for i := range state.Jobs {
		j := &state.Jobs[i]
		if j.Archived || j.Generation != state.Generation || !selectedJob(*j, m.opts.Resources) {
			continue
		}
		if !m.now().Before(j.CreatedAt.Add(portabilityDataLife)) {
			if j.State != "EXPIRED" && j.State != "COMPLETE" {
				j.State = "EXPIRED"
				if err = m.persist(name, state); err != nil {
					return m.jobSnapshot(state), err
				}
			}
			continue
		}
		if j.State == "FAILED" || j.State == "CANCELLED" || j.State == "EXPIRED" || portabilityCooldown(*j, m.now()) != nil {
			continue
		}
		due = append(due, i)
	}
	// Sort indices, preserving durable history order. LastPoll is persisted
	// before every attempt, so fairness needs no separate in-memory cursor.
	sort.Slice(due, func(a, b int) bool {
		left, right := state.Jobs[due[a]], state.Jobs[due[b]]
		if left.LastPoll.Equal(right.LastPoll) {
			return left.ID < right.ID
		}
		return left.LastPoll.Before(right.LastPoll)
	})
	if len(due) > 100 {
		due = due[:100]
	}
	for _, i := range due {
		j := &state.Jobs[i]
		if err = m.pollPortability(ctx, state, name, i); err != nil {
			if isJobFailure(err) {
				j.ReadIssue = true
				if saveErr := m.persist(name, state); saveErr != nil {
					return m.jobSnapshot(state), saveErr
				}
				continue
			}
			return m.jobSnapshot(state), err
		}
	}
	return m.jobSnapshot(state), nil
}

func (m *Manager) jobState(ctx context.Context, online bool) (context.Context, *portabilityState, string, func(), error) {
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return nil, nil, "", nil, err
	}
	fail := func(err error) (context.Context, *portabilityState, string, func(), error) {
		release()
		return nil, nil, "", nil, err
	}
	if m.opts.Provider != "portability" {
		return fail(errors.New("job management is available only for the Data Portability source"))
	}
	if m.opts.Session == nil {
		return fail(errors.New("configure a Google OAuth profile before reading recorded jobs"))
	}
	if online {
		if err = m.configured(); err != nil {
			return fail(err)
		}
	}
	state, name, unlock, err := m.openState(ctx)
	if err != nil {
		return fail(err)
	}
	return ctx, state, name, func() { unlock(); release() }, nil
}

func (m *Manager) jobSnapshot(s *portabilityState) JobSnapshot {
	out := JobSnapshot{Jobs: make([]Job, 0, len(s.Jobs))}
	for _, j := range s.Jobs {
		status := j.State
		if status == "" {
			status = "IN_PROGRESS"
		}
		if !m.now().Before(j.CreatedAt.Add(portabilityDataLife)) && j.Generation == s.Generation {
			status = "EXPIRED"
		}
		message := "Job is pending; refresh respects the five-minute polling cadence"
		switch status {
		case "COMPLETE":
			message = "Archive is ready; signed links last six hours and data lasts fourteen days"
		case "FAILED":
			message = "Google reported failure; retry explicitly, at most three attempts in this retry chain"
		case "CANCELLED":
			message = "Job was cancelled; no replacement was created"
		case "EXPIRED":
			message = "The conservative fourteen-day local data lifetime elapsed; history and grant holds remain"
		}
		if j.SupersededBy != "" {
			message = "A retry successor is recorded; use that job instead of retrying this predecessor"
		}
		if j.Archived {
			message = "Archived locally; archive files and safety receipts are retained"
		}
		if j.ReadIssue {
			message = "This job's last status or immutable archive metadata could not be accepted; its parts remain held while unrelated jobs continue"
		}
		out.Jobs = append(out.Jobs, Job{ID: j.ID, State: status, AccessType: j.AccessType, Message: message, Resources: append([]string(nil), j.Resources...), CreatedAt: j.CreatedAt, LastPoll: j.LastPoll, ExpiresAt: j.CreatedAt.Add(portabilityDataLife), RetryNotBefore: j.RetryNotBefore, RetryCount: j.RetryCount, Archived: j.Archived})
	}
	if s.Intent != nil {
		i := s.Intent
		message := "Creation outcome is unknown. Google cannot enumerate jobs. Supply an evidence-backed ID for this exact request, or use explicit expiry recovery only for a historically recorded time-based grant after fourteen days plus five minutes. No automatic replay is permitted."
		switch intentKind(*i) {
		case "retry":
			message = "Retry outcome is unknown. Reconcile an evidence-backed successor ID from this exact retry, not its predecessor. Expiry recovery requires historically recorded time-based access and fourteen days plus five minutes. No replay is permitted."
		case "cancel":
			message = "Cancellation outcome is unknown. Reconcile the original recorded job ID after Google reports terminal status. No cancellation replay is permitted."
		case "reset":
			message = "Authorization reset outcome is unknown. A login, grant check, job status or file deletion does not prove reset succeeded. The durable safety hold remains; this reset must not be replayed."
		}
		out.Intent = &PendingAction{Kind: intentKind(*i), JobID: i.JobID, Resources: append([]string(nil), i.Resources...), CreatedAt: i.CreatedAt, Message: message}
	}
	return out
}

// CheckAccess checks the current token's grant type. It does not establish that
// a one-time resource is unused, nor clear a historical hold or log the user in.
func (m *Manager) CheckAccess(ctx context.Context) (AccessInfo, error) {
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return AccessInfo{}, err
	}
	defer release()
	if m.opts.Provider != "portability" {
		return AccessInfo{}, errors.New("grant checks are available only for the Data Portability source")
	}
	if err = m.configured(); err != nil {
		return AccessInfo{}, err
	}
	return m.checkAccess(ctx)
}

func (m *Manager) checkAccess(ctx context.Context) (AccessInfo, error) {
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	// accessType.check is a read expressed as POST. Even this POST is not replayed.
	r, err := m.portability.AccessType.Check(&dataportability.CheckAccessTypeRequest{}).Context(callCtx).Do()
	if err != nil {
		return AccessInfo{}, apiError(ctx, "check Data Portability grant", err)
	}
	one, err := normalizeResources(r.OneTimeResources)
	if err != nil {
		return AccessInfo{}, errors.New("Google returned invalid one-time grant resources")
	}
	timed, err := normalizeResources(r.TimeBasedResources)
	if err != nil {
		return AccessInfo{}, errors.New("Google returned invalid time-based grant resources")
	}
	if overlapping(one, timed) {
		return AccessInfo{}, errors.New("Google returned conflicting grant types for the same resources")
	}
	return AccessInfo{OneTimeResources: one, TimeBasedResources: timed}, nil
}

func containsAll(have, want []string) bool {
	set := map[string]bool{}
	for _, r := range have {
		set[r] = true
	}
	for _, r := range want {
		if !set[r] {
			return false
		}
	}
	return true
}
func accessFor(a AccessInfo, resources []string) (string, error) {
	all := append(append([]string(nil), a.OneTimeResources...), a.TimeBasedResources...)
	if !containsAll(all, resources) {
		return "", errors.New("current Google consent does not cover every recorded resource; authorize the exact resource scopes first")
	}
	if overlapping(a.OneTimeResources, resources) {
		return "ACCESS_TYPE_ONE_TIME", nil
	}
	return "ACCESS_TYPE_TIME_BASED", nil
}

func managedJob(s *portabilityState, id string, resources []string) (int, error) {
	if !validDriveID.MatchString(id) {
		return -1, errors.New("invalid Data Portability job ID; supply the recorded Google job ID")
	}
	for i, j := range s.Jobs {
		if j.ID == id {
			if j.Archived || j.Generation != s.Generation {
				return -1, errors.New("job is archived or belongs to a completed authorization generation")
			}
			if !selectedJob(j, resources) {
				return -1, errors.New("recorded job includes resources outside the selected OAuth scopes")
			}
			return i, nil
		}
	}
	return -1, errors.New("job ID is not recorded in this Google profile")
}
func noPending(s *portabilityState) error {
	if s.Intent != nil {
		return errors.New("a previous mutation has an unknown outcome; explicit reconciliation is required before another mutation")
	}
	return nil
}

// RetryJob retries only a recorded FAILED job, never as a watch side effect.
// It returns the new raw Google job ID. Three attempts cover the whole retry
// lineage, including attempts whose remote outcome is unknown.
func (m *Manager) RetryJob(ctx context.Context, id string) (string, error) {
	ctx, s, name, release, err := m.jobState(ctx, true)
	if err != nil {
		return "", err
	}
	defer release()
	if err = noPending(s); err != nil {
		return "", err
	}
	i, err := managedJob(s, id, m.opts.Resources)
	if err != nil {
		return "", err
	}
	j := s.Jobs[i]
	if j.State != "FAILED" || j.SupersededBy != "" {
		return "", errors.New("only a failed job without a recorded successor can be retried")
	}
	if !m.now().Before(j.CreatedAt.Add(portabilityDataLife)) {
		return "", errors.New("expired jobs cannot be retried")
	}
	if m.now().Before(j.RetryNotBefore) {
		return "", errors.New("Google backoff is still active; wait before explicitly retrying this job")
	}
	if activeJobs(s) >= maxJobs {
		return "", errors.New("active history limit reached; archive terminal job history first")
	}
	if len(s.Jobs) >= maxDriveEntries {
		return "", errors.New("retained job safety receipt limit reached; history cannot be discarded to bypass safety guards")
	}
	root := j.RetryRoot
	if root == "" {
		root = j.ID
	}
	count := j.RetryCount
	for _, other := range s.Jobs {
		if other.ID == root || other.RetryRoot == root {
			if other.RetryCount > count {
				count = other.RetryCount
			}
		}
	}
	if count >= 3 {
		return "", errors.New("the three-attempt retry limit for this job lineage has been reached")
	}
	count++
	for k := range s.Jobs {
		if s.Jobs[k].ID == root || s.Jobs[k].RetryRoot == root {
			s.Jobs[k].RetryCount = count
		}
	}
	s.Intent = &creationIntent{Kind: "retry", JobID: id, Resources: append([]string(nil), j.Resources...), CreatedAt: m.now().UTC(), AccessType: j.AccessType, RetryCount: count}
	if err = m.persist(name, s); err != nil {
		return "", err
	}
	callCtx, cancel := m.creationContext(ctx, s.Intent.CreatedAt)
	if callCtx.Err() != nil {
		cancel()
		return "", errors.New("retry intent persistence consumed the request margin; no mutation was sent; the durable intent is retained for explicit reconciliation")
	}
	r, err := m.portability.ArchiveJobs.Retry("archiveJobs/"+id, &dataportability.RetryPortabilityArchiveRequest{}).Context(callCtx).Do()
	cancel()
	if err != nil {
		return "", m.mutationFailure(ctx, s, name, "retry Data Portability job", err)
	}
	if !validDriveID.MatchString(r.ArchiveJobId) || hasJob(s, r.ArchiveJobId) {
		return "", errors.New("Google returned invalid or duplicate retry metadata; the durable retry intent remains held")
	}
	s.Jobs[i].SupersededBy = r.ArchiveJobId
	s.Jobs = append(s.Jobs, portabilityJob{ID: r.ArchiveJobId, Resources: append([]string(nil), j.Resources...), CreatedAt: s.Intent.CreatedAt, AccessType: j.AccessType, State: "IN_PROGRESS", LastPoll: s.Intent.CreatedAt, RetryCount: count, RetryRoot: root, Generation: s.Generation})
	s.Intent = nil
	if err = m.persist(name, s); err != nil {
		return "", fmt.Errorf("Google accepted retry job %s but its receipt could not be saved; the durable retry intent requires explicit reconciliation", r.ArchiveJobId)
	}
	return r.ArchiveJobId, nil
}
func hasJob(s *portabilityState, id string) bool {
	for _, j := range s.Jobs {
		if j.ID == id {
			return true
		}
	}
	return false
}

// CancelJob cancels only a live time-based job. A lost response retains an intent
// which must be resolved with ReconcileIntent(ctx, the same recorded job ID).
func (m *Manager) CancelJob(ctx context.Context, id string) error {
	ctx, s, name, release, err := m.jobState(ctx, true)
	if err != nil {
		return err
	}
	defer release()
	if err = noPending(s); err != nil {
		return err
	}
	i, err := managedJob(s, id, m.opts.Resources)
	if err != nil {
		return err
	}
	j := &s.Jobs[i]
	if j.State != "IN_PROGRESS" || j.AccessType != "ACCESS_TYPE_TIME_BASED" || !m.now().Before(j.CreatedAt.Add(portabilityDataLife)) {
		return errors.New("only an in-progress, unexpired time-based job can be cancelled")
	}
	if m.now().Before(j.RetryNotBefore) {
		return errors.New("Google backoff is still active; wait before cancelling")
	}
	s.Intent = &creationIntent{Kind: "cancel", JobID: id, Resources: append([]string(nil), j.Resources...), CreatedAt: m.now().UTC(), AccessType: j.AccessType}
	if err = m.persist(name, s); err != nil {
		return err
	}
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	_, err = m.portability.ArchiveJobs.Cancel("archiveJobs/"+id, &dataportability.CancelPortabilityArchiveRequest{}).Context(callCtx).Do()
	cancel()
	if err != nil {
		return m.mutationFailure(ctx, s, name, "cancel Data Portability job", err)
	}
	j.State = "CANCELLED"
	s.Intent = nil
	delete(m.links, id)
	return m.persist(name, s)
}

func (m *Manager) mutationFailure(ctx context.Context, s *portabilityState, name, operation string, err error) error {
	var g *googleapi.Error
	if errors.As(err, &g) && g.Code >= 400 && g.Code < 500 && g.Code != 408 {
		// A definitive rejection is safe to acknowledge, but never retry this POST.
		safe := apiError(ctx, operation, err)
		if s.Intent != nil && s.Intent.JobID != "" {
			for i := range s.Jobs {
				if s.Jobs[i].ID == s.Intent.JobID {
					if delay := RetryDelay(safe); delay > 0 {
						s.Jobs[i].RetryNotBefore = m.now().Add(delay).UTC()
					}
				}
			}
		}
		s.Intent = nil
		if saveErr := m.persist(name, s); saveErr != nil {
			return saveErr
		}
		return &sourceError{message: safe.Error(), authorization: errors.Is(safe, auth.ErrAuthorizationRequired)}
	}
	if ctx.Err() != nil {
		return fmt.Errorf("%w; mutation outcome is unknown and its durable intent blocks replay", ctx.Err())
	}
	return errors.New(operation + ": outcome is unknown; durable intent retained; do not replay the mutation")
}

// ReconcileIntent adopts an evidence-backed Google job ID. The caller must
// establish that it belongs to this exact pending request and account: Google's
// status API does not return resources, grant identity, or creation time and has
// no enumeration endpoint. We retain the intent's original time/resources.
// Legacy creation intents with no recorded access type are conservatively
// classified as one-time, even if the current token has time-based access.
func (m *Manager) ReconcileIntent(ctx context.Context, knownJobID string) error {
	ctx, s, name, release, err := m.jobState(ctx, true)
	if err != nil {
		return err
	}
	defer release()
	if s.Intent == nil {
		return errors.New("there is no pending mutation to reconcile")
	}
	intent := s.Intent
	kind := intentKind(*intent)
	if kind == "reset" {
		return errors.New("authorization reset outcome is unknown; neither signing in nor a job ID proves reset succeeded; the safety hold is retained")
	}
	if !validDriveID.MatchString(knownJobID) {
		return errors.New("provide an evidence-backed Google job ID; Google has no job enumeration API and this app cannot invent an ID")
	}
	if !containsAll(m.opts.Resources, intent.Resources) {
		return errors.New("authorize all resources from the durable pending intent before reconciling")
	}
	if kind == "cancel" {
		if knownJobID != intent.JobID {
			return errors.New("cancellation reconciliation requires the original recorded job ID")
		}
		for i := range s.Jobs {
			j := &s.Jobs[i]
			if j.ID != knownJobID {
				continue
			}
			if j.State != "COMPLETE" && j.State != "FAILED" && j.State != "CANCELLED" {
				if err = m.pollPortability(ctx, s, name, i); err != nil {
					return err
				}
			}
			if j.State == "IN_PROGRESS" || j.State == "" || j.State == "EXPIRED" {
				return errors.New("cancellation is not proven terminal; retain the intent and reconcile status after the cadence interval")
			}
			s.Intent = nil
			return m.persist(name, s)
		}
		return errors.New("the pending cancellation has no recorded job; restore the journal")
	}
	if kind != "create" && kind != "retry" {
		return errors.New("this pending action cannot be reconciled by job ID")
	}
	if hasJob(s, knownJobID) {
		return errors.New("the supplied ID is already recorded; it cannot establish a new unknown creation or retry outcome")
	}
	if activeJobs(s) >= maxJobs {
		return errors.New("active job history limit reached; archive terminal history before adopting a job")
	}
	if len(s.Jobs) >= maxDriveEntries {
		return errors.New("retained job safety receipt limit reached; safety history must be preserved")
	}
	if err = portabilityCooldown(portabilityJob{LastPoll: intent.LastPoll, RetryNotBefore: intent.RetryNotBefore}, m.now()); err != nil {
		return err
	}
	// Record the read slot before either API request. Failures and changing the
	// candidate ID cannot bypass the minimum status polling cadence.
	intent.LastPoll = m.now().UTC()
	if err = m.persist(name, s); err != nil {
		return err
	}
	access, err := m.checkAccess(ctx)
	if err != nil {
		return m.intentReadFailure(s, name, err)
	}
	if _, err = accessFor(access, intent.Resources); err != nil {
		return err
	}
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	r, err := m.portability.ArchiveJobs.GetPortabilityArchiveState("archiveJobs/" + knownJobID + "/portabilityArchiveState").Context(callCtx).Do()
	cancel()
	if err != nil {
		return m.intentReadFailure(s, name, apiError(ctx, "reconcile Data Portability job", err))
	}
	if r.Name != "archiveJobs/"+knownJobID+"/portabilityArchiveState" {
		return errors.New("Google did not return the exact known job identity; the pending intent remains held")
	}
	switch r.State {
	case "IN_PROGRESS", "COMPLETE", "FAILED", "CANCELLED":
	default:
		return errors.New("Google returned an unknown state; the pending intent remains held")
	}
	accessType := intent.AccessType
	if accessType == "" {
		accessType = "ACCESS_TYPE_ONE_TIME"
	}
	j := portabilityJob{ID: knownJobID, Resources: append([]string(nil), intent.Resources...), CreatedAt: intent.CreatedAt, AccessType: accessType, LastPoll: intent.LastPoll, State: "IN_PROGRESS", RetryCount: intent.RetryCount, Generation: s.Generation}
	parent := -1
	if kind == "retry" {
		for i, old := range s.Jobs {
			if old.ID == intent.JobID {
				parent = i
				j.RetryRoot = old.RetryRoot
				if j.RetryRoot == "" {
					j.RetryRoot = old.ID
				}
				break
			}
		}
		if parent < 0 {
			return errors.New("pending retry has no predecessor; restore the durable journal")
		}
	}
	// A prospective copy prevents invalid metadata from publishing a partly
	// adopted job in the durable journal.
	next := *s
	next.Jobs = append(append([]portabilityJob(nil), s.Jobs...), j)
	if parent >= 0 {
		next.Jobs[parent].SupersededBy = knownJobID
	}
	next.Intent = nil
	return m.acceptPortability(ctx, &next, name, len(next.Jobs)-1, r, cachedLinks{})
}

func (m *Manager) intentReadFailure(s *portabilityState, name string, err error) error {
	if delay := RetryDelay(err); delay > 0 && s.Intent != nil {
		s.Intent.RetryNotBefore = m.now().Add(delay).UTC()
		if saveErr := m.persist(name, s); saveErr != nil {
			return saveErr
		}
	}
	return err
}

// ResolveExpiredIntent is deliberately narrow. Google's methods guide documents
// an internal reset 14 days after the first initiate request. After 14 days plus
// a five-minute request margin, explicitly retire only a creation or retry intent
// with historically recorded time-based access and exclusively current time-based
// scope evidence.
// A tombstone keeps original resources, timestamps, and any one-time ambiguity;
// neither a checkbox, login, grant check, nor file deletion clears such holds.
// One-time/unknown grants and unknown cancel/reset outcomes remain held.
// https://developers.google.com/data-portability/user-guide/methods
func (m *Manager) ResolveExpiredIntent(ctx context.Context) error {
	ctx, s, name, release, err := m.jobState(ctx, true)
	if err != nil {
		return err
	}
	defer release()
	if s.Intent == nil {
		return errors.New("there is no pending mutation to reconcile")
	}
	i := *s.Intent
	if intentKind(i) != "create" && intentKind(i) != "retry" {
		return errors.New("only creation or retry intents can use policy-checked expiry; cancellation and reset outcomes must be reconciled separately")
	}
	if m.now().Before(i.CreatedAt.Add(portabilityIntentExpiry)) {
		return errors.New("the pending request has not reached the conservative fourteen-day policy expiry plus five-minute request margin; retain the hold or supply an evidence-backed job ID")
	}
	if i.AccessType != "ACCESS_TYPE_TIME_BASED" {
		return errors.New("expiry requires a time-based grant recorded before the original request; one-time and historically unknown grants retain their hold and require an evidence-backed known job ID")
	}
	if !containsAll(m.opts.Resources, i.Resources) {
		return errors.New("select every resource from the durable pending intent before reconciling expiry")
	}
	if err = portabilityCooldown(portabilityJob{LastPoll: i.LastPoll, RetryNotBefore: i.RetryNotBefore}, m.now()); err != nil {
		return err
	}
	s.Intent.LastPoll = m.now().UTC()
	if err = m.persist(name, s); err != nil {
		return err
	}
	access, err := m.checkAccess(ctx)
	if err != nil {
		return m.intentReadFailure(s, name, err)
	}
	if !containsAll(access.TimeBasedResources, i.Resources) || overlapping(access.OneTimeResources, i.Resources) {
		return errors.New("expiry is insufficient evidence for a one-time or unknown grant; retain the hold and reconcile an evidence-backed known job ID")
	}
	if len(s.ExpiredIntents) >= maxJobs {
		return errors.New("expired-intent safety receipt limit reached; safety receipts cannot be discarded")
	}
	s.ExpiredIntents = append(s.ExpiredIntents, i)
	s.Intent = nil
	return m.persist(name, s)
}

// ArchiveJob hides terminal history from discovery and the active 1000-job
// limit. It deletes no archives, parts, grant/cadence receipts, or retry lineage.
func (m *Manager) ArchiveJob(ctx context.Context, id string) error {
	_, s, name, release, err := m.jobState(ctx, false)
	if err != nil {
		return err
	}
	defer release()
	for i := range s.Jobs {
		j := &s.Jobs[i]
		if j.ID != id {
			continue
		}
		if s.Intent != nil && s.Intent.JobID == id {
			return errors.New("reconcile the pending action for this job before archiving it")
		}
		if j.Archived {
			return nil
		}
		if j.State != "COMPLETE" && j.State != "FAILED" && j.State != "CANCELLED" && j.State != "EXPIRED" && m.now().Before(j.CreatedAt.Add(portabilityDataLife)) {
			return errors.New("only terminal or expired jobs can be archived")
		}
		j.Archived = true
		delete(m.links, id)
		return m.persist(name, s)
	}
	return errors.New("job ID is not recorded in this Google profile")
}

// ResetAuthorization implements Google's destructive reset semantics, never a
// new-job convenience. Only a successful remote reset advances the generation.
// Unknown outcomes retain their mutation intent indefinitely; a new login is
// not proof of reset. Local files and prior-generation receipts remain intact.
// https://developers.google.com/data-portability/reference/rest/v1/authorization/reset
func (m *Manager) ResetAuthorization(ctx context.Context) error {
	if ctx == nil || ctx.Value(resetConfirmationKey{}) != true {
		return errors.New("explicit user confirmation is required; call WithResetConfirmation only after explaining consent revocation and removal of previous remote archives")
	}
	ctx, s, name, release, err := m.jobState(ctx, true)
	if err != nil {
		return err
	}
	defer release()
	if err = noPending(s); err != nil {
		return err
	}
	if len(s.ExpiredIntents) > 0 {
		return errors.New("expired unknown creation safety receipts prevent proving all Google jobs completed; reset is blocked")
	}
	if s.NeedsConsent {
		return errors.New("authorization was already reset; obtain fresh consent before starting another export")
	}
	count := 0
	for _, j := range s.Jobs {
		if j.Generation == s.Generation {
			count++
			if j.State != "COMPLETE" || len(j.Parts) == 0 {
				return errors.New("authorization reset requires every known job, including archived history, to be COMPLETE with immutable archive metadata")
			}
		}
	}
	if count == 0 {
		return errors.New("no known completed jobs establish a safe reset; reset is not a new-job setup action")
	}
	if err = m.verifyDownloaded(s); err != nil {
		return err
	}
	s.Intent = &creationIntent{Kind: "reset", Resources: []string{}, CreatedAt: m.now().UTC()}
	if err = m.persist(name, s); err != nil {
		return err
	}
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	_, err = m.portability.Authorization.Reset(&dataportability.ResetAuthorizationRequest{}).Context(callCtx).Do()
	cancel()
	if err != nil {
		return m.mutationFailure(ctx, s, name, "reset Data Portability authorization", err)
	}
	for i := range s.Jobs {
		if s.Jobs[i].Generation == s.Generation {
			s.Jobs[i].Archived = true
		}
	}
	s.Generation++
	s.NeedsConsent = true
	s.Intent = nil
	m.links = map[string]cachedLinks{}
	if err = m.persist(name, s); err != nil {
		return errors.New("Google accepted authorization reset but its receipt could not be saved; the durable reset intent remains held and must not be replayed")
	}
	return nil
}

func (m *Manager) verifyDownloaded(s *portabilityState) error {
	records, err := downloadstate.ReadSnapshot(filepath.Join(m.opts.DataDir, "takedown-state.json"))
	if err != nil {
		return errors.New("cannot verify durable download completions; authorization reset is blocked")
	}
	byKey := map[string]downloadstate.Record{}
	for _, r := range records {
		byKey[r.Part.Key()] = r
	}
	for _, j := range s.Jobs {
		if j.Generation != s.Generation {
			continue
		}
		export := jobExport(j, make([]string, len(j.Parts)))
		for _, p := range export.Parts {
			r, ok := byKey[p.Key()]
			digest, digestErr := hex.DecodeString(r.SHA256)
			// URL and RelativePath describe transport/local placement, not the
			// immutable provider part identity. Never compare whole Part values.
			if !ok || !r.Complete || r.Offset != p.Size || digestErr != nil || len(digest) != 32 || r.Part.ExportID != p.ExportID || r.Part.Index != p.Index || r.Part.Filename != p.Filename || r.Part.Size != p.Size || r.Part.MD5 != p.MD5 {
				return errors.New("every archive part needs an exact immutable identity, full byte offset, and durable SHA256 completion receipt before authorization reset")
			}
		}
	}
	return nil
}
