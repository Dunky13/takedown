package source

import (
	"context"
	"errors"
	"sort"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type Folder struct{ ID, Name string }

// Resources returns an independent, sorted catalog for explicit scope choices.
func Resources() []string {
	resources := make([]string, 0, len(resourceScopes))
	for r := range resourceScopes {
		resources = append(resources, r)
	}
	sort.Strings(resources)
	return resources
}

// BrowseFolders lists direct real-folder children only, following explicit user
// navigation. Empty parentID means My Drive root. It never discovers archives or
// changes the configured download subtree, and never follows shortcuts.
func (m *Manager) BrowseFolders(ctx context.Context, parentID string) ([]Folder, error) {
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	if m.opts.Provider != "drive" {
		return nil, errors.New("folder browsing is available only for the Drive source")
	}
	if err = m.configured(); err != nil {
		return nil, err
	}
	if parentID == "" {
		parentID = "root"
	}
	if !validDriveID.MatchString(parentID) {
		return nil, errors.New("invalid Drive parent folder ID")
	}
	parent, err := m.getDrive(ctx, parentID)
	if err != nil {
		return nil, err
	}
	// The API resolves the root alias to its actual file ID.
	if parentID == "root" && parent != nil {
		parentID = parent.Id
	}
	if !validDriveID.MatchString(parentID) {
		return nil, errors.New("Drive returned an invalid root folder ID")
	}
	if err = checkFolder(parent, parentID); err != nil {
		return nil, err
	}
	result := []Folder{}
	seen, tokens := map[string]bool{parentID: true}, map[string]bool{}
	token := ""
	for pageCount := 0; pageCount < maxDrivePages; pageCount++ {
		page, err := readAPI(ctx, "browse Drive folders", func(ctx context.Context) (*drive.FileList, error) {
			call := m.drive.Files.List().Q("'" + parentID + "' in parents and trashed = false and mimeType = '" + driveFolder + "'").Spaces("drive").PageSize(1000).PageToken(token).SupportsAllDrives(true).IncludeItemsFromAllDrives(true).Fields(googleapi.Field("nextPageToken,incompleteSearch,files(" + driveFields + ")")).Context(ctx)
			if parent.DriveId != "" {
				call = call.Corpora("drive").DriveId(parent.DriveId)
			} else {
				call = call.Corpora("user")
			}
			return call.Do()
		})
		if err != nil {
			return nil, err
		}
		if page.IncompleteSearch {
			return nil, errors.New("Drive folder search is incomplete; no partial folder choices are accepted")
		}
		for _, f := range page.Files {
			if f == nil || !validDriveID.MatchString(f.Id) || seen[f.Id] || len(f.Parents) != 1 || f.Parents[0] != parentID || f.Name == "" || len(f.Name) > 1024 || strings.ContainsAny(f.Name, "\x00\r\n") {
				return nil, errors.New("Drive returned ambiguous folder metadata")
			}
			if err = checkFolder(f, f.Id); err != nil {
				return nil, err
			}
			seen[f.Id] = true
			result = append(result, Folder{ID: f.Id, Name: f.Name})
			if len(result) > maxDriveFolders {
				return nil, errors.New("Drive folder choice limit reached; browse a smaller parent folder")
			}
		}
		if page.NextPageToken == "" {
			sort.Slice(result, func(i, j int) bool {
				if result[i].Name == result[j].Name {
					return result[i].ID < result[j].ID
				}
				return result[i].Name < result[j].Name
			})
			return result, nil
		}
		if page.NextPageToken == token || tokens[page.NextPageToken] {
			return nil, errors.New("Drive folder pagination loop detected")
		}
		token = page.NextPageToken
		tokens[token] = true
	}
	return nil, errors.New("Drive folder page limit reached")
}
