//go:build !headless

package gui

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/source"
)

type resourceChooser struct {
	s        *screen
	filter   *widget.Entry
	checks   map[string]*widget.Check
	names    []string
	summary  *widget.Label
	scroll   *container.Scroll
	manual   *widget.Button
	updating bool
}

func (s *screen) resourceContent() fyne.CanvasObject {
	v := &resourceChooser{s: s, names: source.Resources(), checks: make(map[string]*widget.Check)}
	s.resourceUI = v
	sort.Strings(v.names)
	v.filter = widget.NewEntry()
	v.filter.SetPlaceHolder("Filter resource groups, e.g. Chrome or myactivity")
	v.summary = wrappedLabel("")
	rows := container.NewVBox()
	for _, name := range v.names {
		v.checks[name] = widget.NewCheck(name, func(bool) { v.choose() })
		rows.Add(v.checks[name])
	}
	v.filter.OnChanged = func(text string) {
		for name, check := range v.checks {
			if strings.Contains(strings.ToLower(name), strings.ToLower(strings.TrimSpace(text))) {
				check.Show()
			} else {
				check.Hide()
			}
		}
		rows.Refresh()
	}
	manual := container.NewVBox(wrappedLabel("Advanced: comma-separated published resource group names. These entries and the checkboxes stay synchronized."), s.resources)
	manual.Hide()
	v.manual = widget.NewButton("Advanced: enter resource names…", func() {
		if manual.Visible() {
			manual.Hide()
		} else {
			manual.Show()
		}
		s.portabilityOptions.Refresh()
	})
	v.sync()
	v.scroll = container.NewVScroll(rows)
	return container.NewBorder(container.NewVBox(wrappedLabel("Choose resource groups to authorize and export. Filtering keeps other selections."), v.filter), container.NewVBox(v.summary, v.manual, manual), nil, nil, v.scroll)
}

func (v *resourceChooser) sync() {
	v.updating = true
	defer func() { v.updating = false }()
	selected := make(map[string]bool)
	for _, name := range strings.Split(v.s.resources.Text, ",") {
		if name = strings.TrimSpace(name); name != "" {
			selected[name] = true
		}
	}
	for name, check := range v.checks {
		check.SetChecked(selected[name])
	}
	unknown := 0
	for name := range selected {
		if _, ok := v.checks[name]; !ok {
			unknown++
		}
	}
	text := "No resource groups selected."
	if len(selected) > 0 {
		names := make([]string, 0, len(selected))
		for name := range selected {
			names = append(names, name)
		}
		sort.Strings(names)
		count := len(names)
		if count > 2 {
			names = names[:2]
		}
		text = fmt.Sprintf("%d groups: %s", count, strings.Join(names, ", "))
		if count > 2 {
			text += fmt.Sprintf(" (+%d more)", count-2)
		}
	}
	if unknown > 0 {
		text += " (includes unpublished names; review Advanced before authorizing)"
	}
	v.summary.SetText("Selected: " + text)
}

func (v *resourceChooser) choose() {
	if v.updating || v.s.busy {
		return
	}
	var selected []string
	for _, name := range v.names {
		if v.checks[name].Checked {
			selected = append(selected, name)
		}
	}
	v.s.resources.SetText(strings.Join(selected, ", "))
}

func (v *resourceChooser) setBusy(busy bool) {
	for _, check := range v.checks {
		if busy {
			check.Disable()
		} else {
			check.Enable()
		}
	}
	if busy {
		v.filter.Disable()
		v.manual.Disable()
	} else {
		v.filter.Enable()
		v.manual.Enable()
	}
}
