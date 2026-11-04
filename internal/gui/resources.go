//go:build !headless

package gui

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/source"
)

var resourceGroupTitles = map[string]string{
	"alerts": "Alerts", "chrome": "Chrome", "discover": "Discover", "maps": "Maps", "myactivity": "My Activity",
	"mymaps": "My Maps", "nest": "Nest", "order_reserve": "Orders & reservations", "pixel": "Pixel", "play": "Google Play",
	"saved": "Saved", "search_ugc": "Search contributions", "searchnotifications": "Search notifications",
	"shopping": "Shopping", "streetview": "Street View", "youtube": "YouTube",
}

// resourceGroup is one API prefix. Its children keep the exact published names.
type resourceGroup struct {
	title    string
	names    []string
	check    *widget.Check
	toggle   *widget.Button
	row      *fyne.Container
	children *fyne.Container
	expanded bool
}

type resourceChooser struct {
	s        *screen
	filter   *widget.Entry
	checks   map[string]*widget.Check
	labels   map[string]string // lower-cased "<group> <child>" text each check shows
	groups   []*resourceGroup
	names    []string
	summary  *widget.Label
	scroll   *container.Scroll
	manual   *widget.Button
	updating bool
}

// humanize turns an API identifier such as reading_list into "Reading list".
func humanize(id string) string {
	text := strings.NewReplacer("_", " ", ".", " ").Replace(id)
	if text == "" {
		return text
	}
	return strings.ToUpper(text[:1]) + text[1:]
}

// groupTitle falls back to the humanized prefix so a newly published resource still appears.
func groupTitle(prefix string) string {
	if title, ok := resourceGroupTitles[prefix]; ok {
		return title
	}
	return humanize(prefix)
}

func (s *screen) resourceContent() fyne.CanvasObject {
	v := &resourceChooser{s: s, names: source.Resources(), checks: make(map[string]*widget.Check), labels: make(map[string]string)}
	s.resourceUI = v
	sort.Strings(v.names)
	v.filter = widget.NewEntry()
	v.filter.SetPlaceHolder("Filter, e.g. Chrome or YouTube")
	v.summary = wrappedLabel("")
	byPrefix := make(map[string]*resourceGroup)
	for _, name := range v.names {
		prefix, child, found := strings.Cut(name, ".")
		if !found {
			child = name
		}
		g := byPrefix[prefix]
		if g == nil {
			g = &resourceGroup{title: groupTitle(prefix)}
			byPrefix[prefix] = g
			v.groups = append(v.groups, g)
		}
		g.names = append(g.names, name)
		v.checks[name] = widget.NewCheck(humanize(child), func(bool) { v.choose() })
		v.labels[name] = strings.ToLower(g.title + " " + humanize(child))
	}
	sort.Slice(v.groups, func(i, j int) bool { return v.groups[i].title < v.groups[j].title })
	rows := container.NewVBox()
	for _, g := range v.groups {
		kids := container.NewVBox()
		for _, name := range g.names {
			kids.Add(v.checks[name])
		}
		g.children = container.New(layout.NewCustomPaddedLayout(0, 0, theme.Padding()*8, 0), kids)
		g.check = widget.NewCheck(g.title, func(on bool) { v.toggleGroup(g, on) })
		g.toggle = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
			g.expanded = !g.expanded
			v.applyFilter()
		})
		g.toggle.Importance = widget.LowImportance
		g.row = container.NewVBox(container.NewHBox(g.toggle, g.check), g.children)
		rows.Add(g.row)
	}
	v.filter.OnChanged = func(string) { v.applyFilter() }
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
	v.manual.Importance = widget.LowImportance
	v.sync()
	for _, g := range v.groups {
		g.expanded = v.selectedCount(g) > 0
	}
	v.applyFilter()
	v.scroll = container.NewVScroll(rows)
	return container.NewBorder(container.NewVBox(wrappedLabel("Tick only the data you want. takedown asks Google for access to these items and nothing else."), v.filter), container.NewVBox(v.summary, v.manual, manual), nil, nil, v.scroll)
}

func (v *resourceChooser) selectedCount(g *resourceGroup) int {
	n := 0
	for _, name := range g.names {
		if v.checks[name].Checked {
			n++
		}
	}
	return n
}

func (v *resourceChooser) query() string { return strings.ToLower(strings.TrimSpace(v.filter.Text)) }

// matches accepts the exact API name or the text a person sees.
func (v *resourceChooser) matches(name, query string) bool {
	return query == "" || strings.Contains(strings.ToLower(name), query) || strings.Contains(v.labels[name], query)
}

// applyFilter shows matching children, hides empty groups and, while filtering,
// treats every remaining group as expanded without losing the user's own state.
func (v *resourceChooser) applyFilter() {
	query := v.query()
	for _, g := range v.groups {
		hasMatch := false
		for _, name := range g.names {
			matched := v.matches(name, query)
			setVisible(v.checks[name], matched)
			hasMatch = hasMatch || matched
		}
		open := g.expanded || query != ""
		setVisible(g.children, open)
		if open {
			g.toggle.SetIcon(theme.MenuDropDownIcon())
		} else {
			g.toggle.SetIcon(theme.NavigateNextIcon())
		}
		setVisible(g.row, hasMatch)
	}
	if v.scroll != nil {
		v.scroll.Content.Refresh()
	}
}

// toggleGroup changes only the children the current filter shows, so filtering
// keeps its promise that other selections are untouched.
func (v *resourceChooser) toggleGroup(g *resourceGroup, on bool) {
	if v.updating || v.s.busy {
		return
	}
	query := v.query()
	v.updating = true
	for _, name := range g.names {
		if v.matches(name, query) {
			v.checks[name].SetChecked(on)
		}
	}
	v.updating = false
	v.choose()
	// The write may not change the text; derive the group state from it either way.
	v.sync()
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
	var parts []string
	for _, g := range v.groups {
		n := v.selectedCount(g)
		g.check.Checked, g.check.Partial = n == len(g.names), n > 0 && n < len(g.names)
		g.check.Text = fmt.Sprintf("%s (%d of %d)", g.title, n, len(g.names))
		g.check.Refresh()
		if n > 0 {
			parts = append(parts, fmt.Sprintf("%s (%d)", g.title, n))
		}
	}
	unknown := 0
	for name := range selected {
		if _, ok := v.checks[name]; !ok {
			unknown++
		}
	}
	text := "No data selected."
	if len(selected) > 0 {
		shown := parts
		if len(shown) > 3 {
			shown = shown[:3]
		}
		text = fmt.Sprintf("Selected %d data types", len(selected))
		if len(shown) > 0 {
			text += ": " + strings.Join(shown, ", ")
		}
		if len(parts) > 3 {
			text += fmt.Sprintf(", +%d more", len(parts)-3)
		}
	}
	if unknown > 0 {
		text += " (includes unpublished names; review Advanced before authorizing)"
	}
	v.summary.SetText(text)
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
	for _, g := range v.groups {
		if busy {
			g.check.Disable()
			g.toggle.Disable()
		} else {
			g.check.Enable()
			g.toggle.Enable()
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
