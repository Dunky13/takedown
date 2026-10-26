package source

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

const maxJobs = 1000

type creationIntent struct {
	Kind           string    `json:"kind,omitempty"`
	JobID          string    `json:"job_id,omitempty"`
	Resources      []string  `json:"resources"`
	CreatedAt      time.Time `json:"created_at"`
	AccessType     string    `json:"access_type,omitempty"`
	RetryCount     int       `json:"retry_count,omitempty"`
	LastPoll       time.Time `json:"last_poll,omitempty"`
	RetryNotBefore time.Time `json:"retry_not_before,omitempty"`
}
type portabilityJob struct {
	ID             string            `json:"id"`
	Resources      []string          `json:"resources"`
	CreatedAt      time.Time         `json:"created_at"`
	AccessType     string            `json:"access_type"`
	LastPoll       time.Time         `json:"last_poll,omitempty"`
	RetryNotBefore time.Time         `json:"retry_not_before,omitempty"`
	State          string            `json:"state,omitempty"`
	Parts          []archiveMetadata `json:"parts,omitempty"`
	RetryCount     int               `json:"retry_count,omitempty"`
	RetryRoot      string            `json:"retry_root,omitempty"`
	SupersededBy   string            `json:"superseded_by,omitempty"`
	Archived       bool              `json:"archived,omitempty"`
	Generation     uint64            `json:"generation,omitempty"`
	ReadIssue      bool              `json:"read_issue,omitempty"`
}
type portabilityState struct {
	Version        int              `json:"version"`
	Intent         *creationIntent  `json:"intent,omitempty"`
	Jobs           []portabilityJob `json:"jobs"`
	ExpiredIntents []creationIntent `json:"expired_intents,omitempty"`
	Generation     uint64           `json:"generation,omitempty"`
	NeedsConsent   bool             `json:"needs_consent,omitempty"`
}

func (m *Manager) openState(ctx context.Context) (*portabilityState, string, func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, "", nil, err
	}
	if m.opts.DataDir == "" {
		return nil, "", nil, errors.New("configure a data directory before using Data Portability jobs")
	}
	profile, err := m.opts.Session.ProfileID(ctx)
	if err != nil {
		return nil, "", nil, apiError(ctx, "determine Google OAuth profile", err)
	}
	digest, err := hex.DecodeString(profile)
	if err != nil || len(digest) != 32 || len(profile) != 64 {
		return nil, "", nil, errors.New("Google OAuth profile ID must be a 64-character SHA256 hex digest")
	}
	dir := filepath.Join(m.opts.DataDir, "source", strings.ToLower(profile))
	for _, path := range []string{filepath.Dir(dir), dir} {
		if err := os.MkdirAll(path, 0700); err != nil {
			return nil, "", nil, errors.New("cannot create private archive source state directory")
		}
		if info, err := os.Lstat(path); err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil, "", nil, errors.New("archive source state directory must be a real directory")
		}
	}
	unlock, err := lockState(filepath.Join(dir, "portability.lock"))
	if err != nil {
		return nil, "", nil, err
	}
	fail := func(err error) (*portabilityState, string, func(), error) { unlock(); return nil, "", nil, err }
	name := filepath.Join(dir, "portability.json")
	info, err := os.Lstat(name)
	if errors.Is(err, os.ErrNotExist) {
		return &portabilityState{Version: 1, Jobs: []portabilityJob{}}, name, unlock, nil
	}
	if err != nil || !info.Mode().IsRegular() || info.Size() > 4<<20 {
		return fail(errors.New("archive job state must be a regular file smaller than 4 MiB"))
	}
	f, err := os.Open(name)
	if err != nil {
		return fail(errors.New("cannot read recorded archive jobs"))
	}
	defer f.Close()
	dec := json.NewDecoder(io.LimitReader(f, (4<<20)+1))
	dec.DisallowUnknownFields()
	var state portabilityState
	if err = dec.Decode(&state); err != nil {
		return fail(errors.New("recorded archive jobs are invalid; restore the state file before creating another export"))
	}
	var extra any
	if dec.Decode(&extra) != io.EOF {
		return fail(errors.New("recorded archive jobs contain trailing data"))
	}
	if err = validateState(&state); err != nil {
		return fail(err)
	}
	return &state, name, unlock, nil
}

func validateState(s *portabilityState) error {
	bad := errors.New("recorded archive job state is invalid; restore the state file before creating another export")
	if s.Version != 1 || activeJobs(s) > maxJobs || len(s.Jobs) > maxDriveEntries || len(s.ExpiredIntents) > maxJobs {
		return bad
	}
	if s.Intent != nil {
		if !validIntent(*s.Intent) {
			return bad
		}
	}
	for _, intent := range s.ExpiredIntents {
		if !validIntent(intent) || intent.AccessType != "ACCESS_TYPE_TIME_BASED" || intentKind(intent) != "create" && intentKind(intent) != "retry" {
			return bad
		}
	}
	seen := map[string]bool{}
	for _, j := range s.Jobs {
		rs, err := normalizeResources(j.Resources)
		if err != nil || len(rs) == 0 || !reflect.DeepEqual(rs, j.Resources) || !validDriveID.MatchString(j.ID) || seen[j.ID] || j.CreatedAt.IsZero() {
			return bad
		}
		seen[j.ID] = true
		if j.AccessType != "ACCESS_TYPE_ONE_TIME" && j.AccessType != "ACCESS_TYPE_TIME_BASED" {
			return bad
		}
		if j.RetryCount < 0 || j.RetryCount > 3 || j.Generation > s.Generation || j.RetryRoot != "" && !validDriveID.MatchString(j.RetryRoot) || j.SupersededBy != "" && !validDriveID.MatchString(j.SupersededBy) {
			return bad
		}
		switch j.State {
		case "", "IN_PROGRESS", "COMPLETE", "FAILED", "CANCELLED", "EXPIRED":
		default:
			return bad
		}
		if len(j.Parts) > maxDriveEntries {
			return bad
		}
		objects, names := map[string]bool{}, map[string]bool{}
		for _, p := range j.Parts {
			digest, digestErr := hex.DecodeString(p.Object)
			folded := strings.ToLower(p.Filename)
			if digestErr != nil || len(digest) != 32 || objects[p.Object] || names[folded] || validateArchive(p.Filename, p.Size, p.MD5) != nil {
				return bad
			}
			objects[p.Object] = true
			names[folded] = true
		}
	}
	return nil
}

func activeJobs(s *portabilityState) int {
	n := 0
	for _, j := range s.Jobs {
		if !j.Archived {
			n++
		}
	}
	return n
}

func intentKind(i creationIntent) string {
	if i.Kind == "" {
		return "create"
	}
	return i.Kind
}

func validIntent(i creationIntent) bool {
	rs, err := normalizeResources(i.Resources)
	if err != nil || !reflect.DeepEqual(rs, i.Resources) || i.CreatedAt.IsZero() || i.RetryCount < 0 || i.RetryCount > 3 {
		return false
	}
	if i.AccessType != "" && i.AccessType != "ACCESS_TYPE_ONE_TIME" && i.AccessType != "ACCESS_TYPE_TIME_BASED" {
		return false
	}
	switch intentKind(i) {
	case "create":
		return len(rs) > 0 && i.JobID == ""
	case "retry", "cancel":
		return len(rs) > 0 && validDriveID.MatchString(i.JobID)
	case "reset":
		return i.JobID == ""
	default:
		return false
	}
}

// Persist intent before POST, and job ID before acknowledging creation. Rename
// plus file/directory sync gives a durable single-writer receipt without tokens
// or signed URLs. A surviving intent is deliberately a non-retryable hold.
func saveState(name string, state *portabilityState) error {
	if err := validateState(state); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return errors.New("cannot encode archive job state")
	}
	if len(data) > 4<<20 {
		return errors.New("archive job journal exceeds 4 MiB; retain all safety receipts; further mutations are held")
	}
	f, err := os.CreateTemp(filepath.Dir(name), ".portability-*")
	if err != nil {
		return errors.New("cannot create atomic archive job state")
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	err = f.Chmod(0600)
	if err == nil {
		_, err = f.Write(data)
	}
	if err == nil {
		err = f.Sync()
	}
	err = errors.Join(err, f.Close())
	if err != nil {
		return errors.New("cannot durably write archive job state")
	}
	return replaceState(tmp, name)
}
