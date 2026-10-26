package source

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"google.golang.org/api/drive/v3"
)

func folderFile(id, parent string) *drive.File {
	f := &drive.File{Id: id, Name: id, MimeType: driveFolder, Capabilities: &drive.FileCapabilities{CanListChildren: true}}
	if parent != "" {
		f.Parents = []string{parent}
	}
	return f
}
func archiveFile(id, name, parent string) *drive.File {
	return &drive.File{Id: id, Name: name, MimeType: "application/zip", Parents: []string{parent}, Size: 1234, Version: 7, Md5Checksum: "0123456789abcdef0123456789abcdef", CreatedTime: "2026-10-01T12:00:00Z", Capabilities: &drive.FileCapabilities{CanDownload: true}}
}

func TestDriveAllPagesRecursiveNoMediaProbe(t *testing.T) {
	var requests []string
	m := testManager(t, Options{DriveFolderID: "root-folder"}, func(r *http.Request) (*http.Response, error) {
		requests = append(requests, r.Method+" "+r.URL.Path+" "+r.URL.RawQuery)
		if r.Method != http.MethodGet || r.URL.Query().Get("alt") == "media" {
			t.Fatal("unexpected Drive body request")
		}
		if r.URL.Path == "/drive/v3/files/root-folder" {
			return responseJSON(r, 200, folderFile("root-folder", "")), nil
		}
		if r.URL.Query().Get("pageSize") != "1000" || !strings.Contains(r.URL.Query().Get("fields"), "incompleteSearch") {
			t.Fatal("missing enumeration guarantees")
		}
		switch r.URL.Query().Get("q") {
		case "'root-folder' in parents and trashed = false":
			switch r.URL.Query().Get("pageToken") {
			case "":
				return responseJSON(r, 200, &drive.FileList{Files: []*drive.File{archiveFile("a", "takeout-001.zip", "root-folder")}, NextPageToken: "page2"}), nil
			case "page2":
				return responseJSON(r, 200, &drive.FileList{NextPageToken: "page3"}), nil // empty intermediate page is not EOF
			case "page3":
				return responseJSON(r, 200, &drive.FileList{Files: []*drive.File{folderFile("sub", "root-folder"), archiveFile("b", "takeout-003.zip", "root-folder")}}), nil
			}
		case "'sub' in parents and trashed = false":
			return responseJSON(r, 200, &drive.FileList{Files: []*drive.File{archiveFile("c", "other.tgz", "sub")}}), nil
		}
		t.Fatalf("unexpected request %s", r.URL)
		return nil, nil
	})
	exports, err := m.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(exports) != 3 || len(requests) != 5 {
		t.Fatalf("exports=%d requests=%d", len(exports), len(requests))
	}
	for _, e := range exports {
		if len(e.Parts) != 1 || e.Parts[0].Index != 1 || e.Parts[0].Size != 1234 || e.Parts[0].MD5 == "" {
			t.Fatal("invented Takeout cardinality", e)
		}
		if err = e.Parts[0].Validate(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDriveRejectsPartialAmbiguousOrInaccessibleResults(t *testing.T) {
	for _, tc := range []struct {
		name   string
		file   func() *drive.File
		page   *drive.FileList
		status int
		loop   bool
	}{
		{name: "incomplete search", page: &drive.FileList{IncompleteSearch: true}},
		{name: "shortcut", file: func() *drive.File {
			f := archiveFile("bad", "link.zip", "root-folder")
			f.MimeType = driveShortcut
			return f
		}},
		{name: "no download permission", file: func() *drive.File {
			f := archiveFile("bad", "bad.zip", "root-folder")
			f.Capabilities.CanDownload = false
			return f
		}},
		{name: "missing capability", file: func() *drive.File { f := archiveFile("bad", "bad.zip", "root-folder"); f.Capabilities = nil; return f }},
		{name: "unknown exact size", file: func() *drive.File { f := archiveFile("bad", "bad.zip", "root-folder"); f.Size = 0; return f }},
		{name: "checksum absent", file: func() *drive.File { f := archiveFile("bad", "bad.zip", "root-folder"); f.Md5Checksum = ""; return f }},
		{name: "unsafe name", file: func() *drive.File { return archiveFile("bad", "../bad.zip", "root-folder") }},
		{name: "trash", file: func() *drive.File { f := archiveFile("bad", "bad.zip", "root-folder"); f.Trashed = true; return f }},
		{name: "outside subtree", file: func() *drive.File { return archiveFile("bad", "bad.zip", "outside") }},
		{name: "ambiguous parents", file: func() *drive.File {
			f := archiveFile("bad", "bad.zip", "root-folder")
			f.Parents = append(f.Parents, "outside")
			return f
		}},
		{name: "duplicate id", file: func() *drive.File { return archiveFile("good", "other.zip", "root-folder") }},
		{name: "cannot list folder", file: func() *drive.File {
			f := folderFile("sub", "root-folder")
			f.Capabilities.CanListChildren = false
			return f
		}},
		{name: "missing timestamp", file: func() *drive.File { f := archiveFile("bad", "bad.zip", "root-folder"); f.CreatedTime = ""; return f }},
		{name: "permissions", status: 403}, {name: "pagination loop", loop: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := testManager(t, Options{DriveFolderID: "root-folder"}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/drive/v3/files/root-folder" {
					return responseJSON(r, 200, folderFile("root-folder", "")), nil
				}
				if tc.status != 0 {
					return responseJSON(r, tc.status, map[string]any{"error": map[string]any{"code": tc.status, "message": "private secret URL"}}), nil
				}
				if r.URL.Query().Get("pageToken") == "" {
					return responseJSON(r, 200, &drive.FileList{Files: []*drive.File{archiveFile("good", "good.zip", "root-folder")}, NextPageToken: "next"}), nil
				}
				page := tc.page
				if page == nil {
					page = &drive.FileList{}
				}
				if tc.file != nil {
					page = &drive.FileList{Files: []*drive.File{tc.file()}}
				}
				if tc.loop {
					page = &drive.FileList{NextPageToken: "next"}
				}
				return responseJSON(r, 200, page), nil
			})
			exports, err := m.List(context.Background())
			if err == nil || exports != nil {
				t.Fatalf("accepted partial results: %v %v", exports, err)
			}
			if strings.Contains(err.Error(), "private secret") {
				t.Fatal("unsanitized error")
			}
		})
	}
}

func TestDriveRefreshPinsMetadataAndSubtreeAcrossRestart(t *testing.T) {
	for _, change := range []string{"none", "version", "name", "size", "checksum", "moved", "trashed"} {
		t.Run(change, func(t *testing.T) {
			file := archiveFile("a", "takeout.zip", "root-folder")
			identity, err := driveFileIdentity(file, "root-folder")
			if err != nil {
				t.Fatal(err)
			}
			switch change {
			case "version":
				file.Version++
			case "name":
				file.Name = "new.zip"
			case "size":
				file.Size++
			case "checksum":
				file.Md5Checksum = strings.Repeat("a", 32)
			case "moved":
				file.Parents = []string{"outside"}
			case "trashed":
				file.Trashed = true
			}
			m := testManager(t, Options{DriveFolderID: "root-folder"}, func(r *http.Request) (*http.Response, error) {
				switch r.URL.Path {
				case "/drive/v3/files/a":
					return responseJSON(r, 200, file), nil
				case "/drive/v3/files/root-folder":
					return responseJSON(r, 200, folderFile("root-folder", "")), nil
				case "/drive/v3/files/outside":
					return responseJSON(r, 200, folderFile("outside", "")), nil
				}
				t.Fatal("unexpected request")
				return nil, nil
			})
			link, err := m.RefreshLink(context.Background(), identity.exportID(), 1)
			if change == "none" {
				if err != nil || link != driveMediaURL("a") {
					t.Fatalf("refresh %q %v", link, err)
				}
			} else if err == nil || link != "" {
				t.Fatal("stale metadata received link", link, err)
			}
		})
	}
}

func TestDriveContextAndCloseCancelAPI(t *testing.T) {
	started := make(chan struct{})
	m := testManager(t, Options{DriveFolderID: "root-folder"}, func(r *http.Request) (*http.Response, error) {
		close(started)
		<-r.Context().Done()
		return nil, r.Context().Err()
	})
	done := make(chan error, 1)
	go func() { _, err := m.List(context.Background()); done <- err }()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("request did not start")
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("Close did not cancel active operation")
	}
}

func TestDriveTraversalLimitsFailWithoutPartialListing(t *testing.T) {
	for _, limit := range []string{"pages", "entries", "folders", "depth"} {
		t.Run(limit, func(t *testing.T) {
			calls := 0
			m := testManager(t, Options{DriveFolderID: "root-folder"}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/drive/v3/files/root-folder" {
					return responseJSON(r, 200, folderFile("root-folder", "")), nil
				}
				calls++
				page := &drive.FileList{}
				switch limit {
				case "pages":
					page.NextPageToken = "page" + strconv.Itoa(calls)
				case "entries":
					for i := 0; i <= maxDriveEntries; i++ {
						f := archiveFile("id"+strconv.Itoa(i), "readme.txt", "root-folder")
						page.Files = append(page.Files, f)
					}
				case "folders":
					for i := 0; i < maxDriveFolders; i++ {
						page.Files = append(page.Files, folderFile("sub"+strconv.Itoa(i), "root-folder"))
					}
				case "depth":
					parent := strings.Split(r.URL.Query().Get("q"), "'")[1]
					page.Files = []*drive.File{folderFile("sub"+strconv.Itoa(calls), parent)}
				}
				return responseJSON(r, 200, page), nil
			})
			if exports, err := m.List(context.Background()); err == nil || exports != nil {
				t.Fatal("limit returned a partial listing", exports, err)
			}
			if limit == "pages" && calls != maxDrivePages {
				t.Fatal("pagination calls exceeded limit", calls)
			}
			if limit == "depth" && calls > maxDriveDepth+1 {
				t.Fatal("depth traversal exceeded limit", calls)
			}
		})
	}
}
