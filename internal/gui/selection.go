//go:build !headless

package gui

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/app"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/state"
)

// humanSize uses decimal units, matching the sizes Google reports.
func humanSize(n int64) string {
	switch {
	case n >= 1e9:
		return fmt.Sprintf("%.2f GB", float64(n)/1e9)
	case n >= 1e6:
		return fmt.Sprintf("%.1f MB", float64(n)/1e6)
	case n >= 1e3:
		return fmt.Sprintf("%.0f KB", float64(n)/1e3)
	}
	return fmt.Sprintf("%d B", n)
}

// downloadedParts lists parts already complete at the current destination.
// A later session can then start from what is left. Unreadable state selects
// everything, and the download engine still verifies what is on disk.
func (s *screen) downloadedParts() map[string]bool {
	identity, err := app.DestinationIdentity(strings.TrimSpace(s.destination.Text))
	if err != nil {
		return nil
	}
	records, err := state.ReadSnapshot(filepath.Join(s.a.DataDir, "takedown-state.json"))
	if err != nil {
		return nil
	}
	done := make(map[string]bool)
	for _, record := range records {
		if record.Complete && record.Destination == identity {
			done[record.Part.Key()] = true
		}
	}
	return done
}

// showExports lists the discovered archives. Drive delivers each file as its own
// one-part export, so only multi-part exports get a parent row. s.skippedParts
// holds the unticked parts, so an empty map selects all.
func (s *screen) showExports(binding sourceSnapshot) {
	s.selected = make(map[string]bool)
	s.selectionChecks = nil
	s.setAllParts = nil
	s.exports.RemoveAll()
	if len(s.ready) == 0 {
		s.exports.Add(wrappedLabel("No ready archives in this source. Create an export, then discover again or start watching."))
		s.updateSelectionSummary()
		return
	}
	// A destination change re-derives only these automatic skips; a person's unticks stay.
	for key := range s.doneParts {
		delete(s.skippedParts, key)
	}
	done := s.downloadedParts()
	s.doneParts = done
	// Programmatic check changes must not re-enter the handlers below.
	updating := false
	allowed := func() bool { return !updating && !s.busy && binding.matches(s) }
	for _, export := range s.ready {
		var total int64
		for _, part := range export.Parts {
			total += part.Size
			if done[part.Key()] {
				s.skippedParts[part.Key()] = true
			}
		}
		created := "date unavailable"
		if !export.CreatedAt.IsZero() {
			created = export.CreatedAt.Format("2006-01-02")
		}
		single := len(export.Parts) == 1
		var exportCheck *widget.Check
		partChecks := make([]*widget.Check, len(export.Parts))
		refresh := func() {
			picked := 0
			for _, part := range export.Parts {
				if !s.skippedParts[part.Key()] {
					picked++
				}
			}
			s.selected[export.ID] = picked > 0
			if exportCheck != nil {
				exportCheck.Checked, exportCheck.Partial = picked > 0 && picked == len(export.Parts), picked > 0 && picked < len(export.Parts)
				exportCheck.Refresh()
			}
			s.updateSelectionSummary()
		}
		setAll := func(on bool) {
			if !allowed() {
				return
			}
			updating = true
			for i, part := range export.Parts {
				s.skipPart(part, !on)
				partChecks[i].SetChecked(on)
			}
			updating = false
			refresh()
		}
		s.setAllParts = append(s.setAllParts, setAll)
		if !single {
			exportCheck = widget.NewCheck(fmt.Sprintf("%s  ·  %d parts  ·  %s", created, len(export.Parts), humanSize(total)), setAll)
			s.exports.Add(exportCheck)
			s.selectionChecks = append(s.selectionChecks, exportCheck)
		}
		for i, part := range export.Parts {
			check := widget.NewCheck("", func(on bool) {
				if !allowed() {
					return
				}
				s.skipPart(part, !on)
				refresh()
			})
			check.Checked = !s.skippedParts[part.Key()]
			partChecks[i] = check
			s.selectionChecks = append(s.selectionChecks, check)
			text := part.Filename + "  ·  " + humanSize(part.Size)
			row := container.New(layout.NewCustomPaddedLayout(0, 0, theme.Padding()*6, 0), check)
			if single {
				text += "  ·  " + created
				row = container.NewStack(check)
			}
			if done[part.Key()] {
				text += "  ·  downloaded earlier"
			}
			s.exports.Add(border(nil, nil, row, nil, archiveLabel(text)))
		}
		refresh()
	}
}

func (s *screen) skipPart(part model.Part, skip bool) {
	if skip {
		s.skippedParts[part.Key()] = true
	} else {
		delete(s.skippedParts, part.Key())
	}
}

// selectedParts is the one definition of what Download fetches and the summary counts.
func (s *screen) selectedParts() []model.Part {
	var parts []model.Part
	for _, export := range s.ready {
		if !s.selected[export.ID] {
			continue
		}
		for _, part := range export.Parts {
			if !s.skippedParts[part.Key()] {
				parts = append(parts, part)
			}
		}
	}
	return parts
}

func (s *screen) selectAllParts(on bool) {
	for _, set := range s.setAllParts {
		set(on)
	}
}

func (s *screen) updateSelectionSummary() {
	total := 0
	for _, export := range s.ready {
		total += len(export.Parts)
	}
	picked := s.selectedParts()
	var size int64
	for _, part := range picked {
		size += part.Size
	}
	s.selectionSummary.SetText(fmt.Sprintf("Selected %d of %d archive files  ·  %s", len(picked), total, humanSize(size)))
	if shown := total > 0; shown != s.selectionRow.Visible() {
		setVisible(s.selectionRow, shown)
		setVisible(s.downloadRow, shown)
		s.find.Importance = widget.HighImportance
		if shown {
			s.find.Importance = widget.MediumImportance
		}
		s.find.Refresh()
		relayout(s.content)
	}
}
