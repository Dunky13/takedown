//go:build !headless

package gui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
)

func groupByPrefix(t *testing.T, v *resourceChooser, prefix string) *resourceGroup {
	t.Helper()
	for _, g := range v.groups {
		if g.title == groupTitle(prefix) {
			return g
		}
	}
	t.Fatalf("no %q group", prefix)
	return nil
}

func TestGroupedCatalogWritesExactNamesAndShowsMixedState(t *testing.T) {
	s, _ := newManagedScreen(t)
	v := s.resourceUI
	s.resources.SetText("")
	chrome := groupByPrefix(t, v, "chrome")
	if chrome.title != "Chrome" || v.checks["chrome.bookmarks"].Text != "Bookmarks" {
		t.Fatalf("labels: group %q child %q", chrome.title, v.checks["chrome.bookmarks"].Text)
	}
	test.Tap(chrome.check)
	for _, name := range chrome.names {
		if !strings.Contains(s.resources.Text, name) {
			t.Fatalf("group toggle left out %q: %s", name, s.resources.Text)
		}
	}
	if !chrome.check.Checked || chrome.check.Partial || !strings.HasPrefix(v.summary.Text, "Selected ") {
		t.Fatalf("full group state: checked=%v partial=%v", chrome.check.Checked, chrome.check.Partial)
	}
	test.Tap(v.checks["chrome.bookmarks"])
	if chrome.check.Checked || !chrome.check.Partial || strings.Contains(s.resources.Text, "chrome.bookmarks") {
		t.Fatal("mixed group is not partial")
	}
	test.Tap(chrome.check) // partial -> everything
	if !chrome.check.Checked || !strings.Contains(s.resources.Text, "chrome.bookmarks") {
		t.Fatal("ticking a partial group did not select all its children")
	}
}

func TestGroupToggleWhileFilteredOnlyTouchesMatches(t *testing.T) {
	s, _ := newManagedScreen(t)
	v := s.resourceUI
	s.resources.SetText("")
	chrome := groupByPrefix(t, v, "chrome")
	v.filter.SetText("bookmarks")
	test.Tap(chrome.check)
	if s.resources.Text != "chrome.bookmarks" || !chrome.check.Partial {
		t.Fatalf("filtered toggle wrote %q (partial=%v)", s.resources.Text, chrome.check.Partial)
	}
	if !chrome.children.Visible() {
		t.Fatal("a filtered group is not shown expanded")
	}
	v.filter.SetText("")
	if chrome.children.Visible() {
		t.Fatal("clearing the filter lost the collapsed state")
	}
}

func TestUnknownPrefixFallsBackToHumanizedName(t *testing.T) {
	if got := groupTitle("brand_new"); got != "Brand new" {
		t.Fatalf("fallback title %q", got)
	}
	if got := humanize("media.reviews_and_stars"); got != "Media reviews and stars" {
		t.Fatalf("humanized %q", got)
	}
}
