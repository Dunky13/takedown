package source

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/Dunky13/takedown/internal/model"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

const (
	driveFolder     = "application/vnd.google-apps.folder"
	driveShortcut   = "application/vnd.google-apps.shortcut"
	driveFields     = "id,name,mimeType,size,md5Checksum,version,createdTime,trashed,parents,driveId,capabilities(canDownload,canListChildren)"
	maxDriveEntries = 10000
	maxDriveFolders = 1000
	maxDrivePages   = 1000
	maxDriveDepth   = 64
)

var validDriveID = regexp.MustCompile(`^[A-Za-z0-9_-]{1,200}$`)

func folderID(value string) (string, error) {
	if validDriveID.MatchString(value) && value != "root" {
		return value, nil
	}
	u, err := url.Parse(value)
	if err == nil && u.Scheme == "https" && u.Host == "drive.google.com" && u.User == nil && u.Fragment == "" {
		p := strings.Split(strings.Trim(u.Path, "/"), "/")
		// The standard folder URL can optionally include the account selector.
		if len(p) == 5 && p[0] == "drive" && p[1] == "u" && regexp.MustCompile(`^[0-9]+$`).MatchString(p[2]) {
			p = append([]string{p[0]}, p[3:]...)
		}
		if len(p) == 3 && p[0] == "drive" && p[1] == "folders" && validDriveID.MatchString(p[2]) && p[2] != "root" {
			return p[2], nil
		}
	}
	return "", errors.New("configure an explicit Google Drive folder ID or https://drive.google.com/drive/folders/ID URL")
}

type driveIdentity struct {
	FileID   string `json:"file"`
	FolderID string `json:"folder"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Version  int64  `json:"version"`
	MD5      string `json:"md5"`
}

func (d driveIdentity) exportID() string {
	b, _ := json.Marshal(d)
	return "drive-" + base64.RawURLEncoding.EncodeToString(b)
}

func parseDriveIdentity(id string) (driveIdentity, error) {
	var d driveIdentity
	if !strings.HasPrefix(id, "drive-") || len(id) > 4096 {
		return d, errors.New("invalid Drive archive identity")
	}
	b, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(id, "drive-"))
	if err != nil || json.Unmarshal(b, &d) != nil || !validDriveID.MatchString(d.FileID) || !validDriveID.MatchString(d.FolderID) || d.Version < 0 || d.exportID() != id {
		return d, errors.New("invalid Drive archive identity")
	}
	if err = validateArchive(d.Name, d.Size, d.MD5); err != nil {
		return d, err
	}
	return d, nil
}

func archiveName(name string) bool {
	s := strings.ToLower(name)
	return strings.HasSuffix(s, ".zip") || strings.HasSuffix(s, ".tgz") || strings.HasSuffix(s, ".tar.gz")
}

func validateArchive(name string, size int64, md5 string) error {
	p := model.Part{ExportID: "validate", Index: 1, Filename: name, Size: size, MD5: md5}
	if err := p.Validate(); err != nil {
		return errors.New("archive metadata requires a safe ZIP/TGZ filename, exact positive size, and valid checksum")
	}
	if len(name) > 255 || strings.HasSuffix(name, " ") || strings.HasSuffix(name, ".") || strings.ContainsAny(name, `<>"|?*`) || strings.IndexFunc(name, unicode.IsControl) >= 0 {
		return errors.New("unsafe archive filename")
	}
	stem := strings.ToUpper(strings.SplitN(name, ".", 2)[0])
	if stem == "CON" || stem == "PRN" || stem == "AUX" || stem == "NUL" || regexp.MustCompile(`^(COM|LPT)[0-9]$`).MatchString(stem) {
		return errors.New("unsafe archive filename")
	}
	return nil
}

func driveFileIdentity(f *drive.File, root string) (driveIdentity, error) {
	if f == nil || !validDriveID.MatchString(f.Id) || f.Trashed || f.MimeType == "" || strings.HasPrefix(f.MimeType, "application/vnd.google-apps.") {
		return driveIdentity{}, errors.New("Drive returned missing, trashed, or unsupported archive metadata")
	}
	if f.Capabilities == nil || !f.Capabilities.CanDownload {
		return driveIdentity{}, errors.New("Drive archive cannot be downloaded; check file permissions and download restrictions")
	}
	if err := validateArchive(f.Name, f.Size, f.Md5Checksum); err != nil {
		return driveIdentity{}, err
	}
	checksum, err := hex.DecodeString(f.Md5Checksum)
	if err != nil || len(checksum) != 16 {
		return driveIdentity{}, errors.New("Drive archive is missing a valid MD5 checksum; refusing an unbound download")
	}
	if f.Version < 0 {
		return driveIdentity{}, errors.New("Drive archive has an invalid version")
	}
	return driveIdentity{f.Id, root, f.Name, f.Size, f.Version, strings.ToLower(f.Md5Checksum)}, nil
}

func (m *Manager) getDrive(ctx context.Context, id string) (*drive.File, error) {
	return readAPI(ctx, "read Drive metadata", func(ctx context.Context) (*drive.File, error) {
		return m.drive.Files.Get(id).SupportsAllDrives(true).Fields(googleapi.Field(driveFields)).Context(ctx).Do()
	})
}

func checkFolder(f *drive.File, id string) error {
	if f == nil || f.Id != id || f.Trashed || f.MimeType != driveFolder {
		return errors.New("configured Drive folder is missing, trashed, or is not a real folder; shortcuts are unsupported")
	}
	if f.Capabilities == nil || !f.Capabilities.CanListChildren {
		return errors.New("Drive folder cannot be listed; grant read access to the selected folder and its children")
	}
	return nil
}

func (m *Manager) listDrive(ctx context.Context) ([]model.Export, error) {
	rootID := m.opts.DriveFolderID
	if rootID == "" {
		return nil, errors.New("configure the explicit Takeout Drive folder ID before listing archives")
	}
	root, err := m.getDrive(ctx, rootID)
	if err != nil {
		return nil, err
	}
	if err = checkFolder(root, rootID); err != nil {
		return nil, err
	}
	type folder struct {
		id    string
		depth int
	}
	queue := []folder{{rootID, 0}}
	seen := map[string]bool{rootID: true}
	pages, entries, folders := 0, 0, 1
	result := []model.Export{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		token := ""
		tokens := map[string]bool{}
		for {
			if err = ctx.Err(); err != nil {
				return nil, err
			}
			pages++
			if pages > maxDrivePages {
				return nil, errors.New("Drive page limit reached; select a smaller archive subtree")
			}
			page, err := readAPI(ctx, "list Drive folder", func(ctx context.Context) (*drive.FileList, error) {
				call := m.drive.Files.List().Q("'" + current.id + "' in parents and trashed = false").Spaces("drive").PageSize(1000).PageToken(token).SupportsAllDrives(true).IncludeItemsFromAllDrives(true).Fields(googleapi.Field("nextPageToken,incompleteSearch,files(" + driveFields + ")")).Context(ctx)
				if root.DriveId != "" {
					call = call.Corpora("drive").DriveId(root.DriveId)
				} else {
					call = call.Corpora("user")
				}
				return call.Do()
			})
			if err != nil {
				return nil, err
			}
			if page.IncompleteSearch {
				return nil, errors.New("Drive reported an incomplete search; no partial archive listing is accepted")
			}
			for _, f := range page.Files {
				entries++
				if entries > maxDriveEntries {
					return nil, errors.New("Drive entry limit reached; select a smaller archive subtree")
				}
				if f == nil || !validDriveID.MatchString(f.Id) || f.Trashed || f.MimeType == "" || len(f.Parents) != 1 || f.Parents[0] != current.id {
					return nil, errors.New("Drive returned ambiguous or incomplete subtree metadata")
				}
				if seen[f.Id] {
					return nil, errors.New("Drive returned a repeated file or folder; refusing an ambiguous subtree")
				}
				seen[f.Id] = true
				if f.MimeType == driveShortcut {
					return nil, errors.New("Drive subtree contains a shortcut; select a real archive folder without shortcuts")
				}
				if f.MimeType == driveFolder {
					if err = checkFolder(f, f.Id); err != nil {
						return nil, err
					}
					folders++
					if folders > maxDriveFolders || current.depth >= maxDriveDepth {
						return nil, errors.New("Drive folder/depth limit reached; select a smaller archive subtree")
					}
					queue = append(queue, folder{f.Id, current.depth + 1})
					continue
				}
				if !archiveName(f.Name) {
					continue
				}
				d, err := driveFileIdentity(f, rootID)
				if err != nil {
					return nil, err
				}
				created, err := time.Parse(time.RFC3339Nano, f.CreatedTime)
				if err != nil || created.IsZero() {
					return nil, errors.New("Drive archive is missing a valid creation timestamp")
				}
				id := d.exportID()
				result = append(result, model.Export{ID: id, CreatedAt: created, Parts: []model.Part{{ExportID: id, Index: 1, Filename: d.Name, Size: d.Size, MD5: d.MD5, URL: driveMediaURL(d.FileID)}}})
			}
			next := page.NextPageToken
			if next == "" {
				break
			}
			if next == token || tokens[next] {
				return nil, errors.New("Drive pagination loop detected; no partial archive listing is accepted")
			}
			tokens[next] = true
			token = next
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].ID < result[j].ID
		}
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	m.emit("Listed Drive archive revisions in the configured subtree; each file is a separate export, without inferring unseen Takeout parts")
	return result, nil
}

// files.get alt=media is stable, not a short-lived signed link. The identity
// includes version/name/size/checksum, RefreshLink verifies it, and Part.MD5
// binds the bytes accepted by the downloader even if Drive changes mid-transfer.
func driveMediaURL(id string) string {
	return "https://www.googleapis.com/drive/v3/files/" + url.PathEscape(id) + "?alt=media&supportsAllDrives=true"
}

func (m *Manager) refreshDrive(ctx context.Context, id string, index int) (string, error) {
	if index != 1 {
		return "", errors.New("each Drive archive revision has exactly one part")
	}
	d, err := parseDriveIdentity(id)
	if err != nil {
		return "", err
	}
	if d.FolderID != m.opts.DriveFolderID {
		return "", errors.New("Drive archive belongs to a different configured subtree; list archives again")
	}
	f, err := m.getDrive(ctx, d.FileID)
	if err != nil {
		return "", err
	}
	current, err := driveFileIdentity(f, d.FolderID)
	if err != nil {
		return "", err
	}
	if current != d {
		return "", errors.New("Drive archive revision changed (name, size, version, or checksum); list archives again instead of resuming stale bytes")
	}
	// Recheck custody: a file moved outside the explicitly selected subtree must
	// not remain downloadable through a previously issued archive identity.
	parents := f.Parents
	seen := map[string]bool{f.Id: true}
	for depth := 0; depth <= maxDriveDepth; depth++ {
		if len(parents) != 1 || seen[parents[0]] {
			return "", errors.New("Drive archive ancestry is ambiguous or no longer within the selected subtree")
		}
		parent := parents[0]
		seen[parent] = true
		folder, err := m.getDrive(ctx, parent)
		if err != nil {
			return "", err
		}
		if err = checkFolder(folder, parent); err != nil {
			return "", err
		}
		if parent == d.FolderID {
			return driveMediaURL(d.FileID), nil
		}
		parents = folder.Parents
	}
	return "", errors.New("Drive archive ancestry exceeds the configured traversal limit")
}
