package source

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"google.golang.org/api/drive/v3"
)

func TestBrowseFoldersExplicitRootAndSharedDrivePages(t *testing.T) {
	for _, shared := range []bool{false, true} {
		t.Run(map[bool]string{false: "root", true: "shared"}[shared], func(t *testing.T) {
			requests := 0
			parentID := "root"
			if shared {
				parentID = "shared-parent"
			}
			m := testManager(t, Options{DriveFolderID: "configured-download-subtree"}, func(r *http.Request) (*http.Response, error) {
				requests++
				if r.Method != http.MethodGet {
					t.Fatal("folder browse mutation")
				}
				if r.URL.Path == "/drive/v3/files/"+parentID {
					f := folderFile("actual-root", "")
					if shared {
						f.Id = parentID
						f.DriveId = "shared-drive"
					}
					return responseJSON(r, 200, f), nil
				}
				actual := "actual-root"
				if shared {
					actual = parentID
					if r.URL.Query().Get("corpora") != "drive" || r.URL.Query().Get("driveId") != "shared-drive" {
						t.Fatal("shared drive omitted")
					}
				}
				if r.URL.Query().Get("q") != "'"+actual+"' in parents and trashed = false and mimeType = '"+driveFolder+"'" {
					t.Fatal("browse searched beyond direct folder children")
				}
				if r.URL.Query().Get("pageToken") == "" {
					return responseJSON(r, 200, &drive.FileList{NextPageToken: "next"}), nil
				}
				a, b := folderFile("a", actual), folderFile("b", actual)
				a.Name = "Takeout"
				b.Name = "Takeout"
				return responseJSON(r, 200, &drive.FileList{Files: []*drive.File{b, a}}), nil
			})
			selected := parentID
			if !shared {
				selected = ""
			}
			got, err := m.BrowseFolders(context.Background(), selected)
			if err != nil || !reflect.DeepEqual(got, []Folder{{ID: "a", Name: "Takeout"}, {ID: "b", Name: "Takeout"}}) || requests != 3 {
				t.Fatal(got, err, requests)
			}
			if m.opts.DriveFolderID != "configured-download-subtree" {
				t.Fatal("browse changed archive discovery scope")
			}
		})
	}
}

func TestBrowseFoldersRejectsIncompleteOrAmbiguousChoices(t *testing.T) {
	for _, kind := range []string{"shortcut", "partial", "duplicate", "outside", "missing-capability", "loop", "unsafe-id", "wrong-provider"} {
		t.Run(kind, func(t *testing.T) {
			calls := 0
			opts := Options{}
			if kind == "wrong-provider" {
				opts.Provider = "portability"
				opts.Resources = []string{"myactivity.search"}
			}
			m := testManager(t, opts, func(r *http.Request) (*http.Response, error) {
				calls++
				if r.URL.Path == "/drive/v3/files/parent" {
					return responseJSON(r, 200, folderFile("parent", "")), nil
				}
				p := &drive.FileList{Files: []*drive.File{folderFile("child", "parent")}}
				switch kind {
				case "shortcut":
					p.Files[0].MimeType = driveShortcut
				case "partial":
					p.IncompleteSearch = true
				case "duplicate":
					p.Files = append(p.Files, folderFile("child", "parent"))
				case "outside":
					p.Files[0].Parents = []string{"elsewhere"}
				case "missing-capability":
					p.Files[0].Capabilities = nil
				case "loop":
					p.Files = nil
					p.NextPageToken = "repeated"
				}
				return responseJSON(r, 200, p), nil
			})
			id := "parent"
			if kind == "unsafe-id" {
				id = "../outside"
			}
			got, err := m.BrowseFolders(context.Background(), id)
			if err == nil || got != nil {
				t.Fatal("ambiguous browse accepted", got, err)
			}
			if (kind == "unsafe-id" || kind == "wrong-provider") && calls != 0 {
				t.Fatal("invalid browse made network call")
			}
		})
	}
}

func TestResourceCatalogIndependentAndContextAwareBrowsing(t *testing.T) {
	catalog := Resources()
	if len(catalog) == 0 {
		t.Fatal("empty resource catalog")
	}
	catalog[0] = "invented"
	if _, err := Scopes("portability", Resources()); err != nil {
		t.Fatal("unsupported catalog", err)
	}
	if Resources()[0] == "invented" {
		t.Fatal("catalog is mutable global")
	}
	m := testManager(t, Options{}, func(r *http.Request) (*http.Response, error) {
		t.Fatal("cancelled browse made request")
		return nil, nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := m.BrowseFolders(ctx, ""); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := m.BrowseFolders(context.Background(), "https://secret.example/?token=private"); err == nil || strings.Contains(err.Error(), "private") {
		t.Fatal("unsafe diagnostic", err)
	}
}
