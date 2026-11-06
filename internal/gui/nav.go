//go:build !headless

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// stackLayout sizes every page to the whole area. Its minimum covers hidden
// pages too, so the window does not resize when the page changes.
type stackLayout struct{}

func (stackLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (stackLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var minimum fyne.Size
	for _, o := range objects {
		minimum = minimum.Max(o.MinSize())
	}
	return minimum
}

// navigator shows one page at a time. It keeps TabItem as the page descriptor
// and the AppTabs method names, but draws no tab bar; the stepper and toolbar
// are the controls.
type navigator struct {
	items      []*container.TabItem
	stack      *fyne.Container
	selected   int
	disabled   map[*container.TabItem]bool
	OnSelected func(*container.TabItem)
	// onEnabledChanged lets the toolbar follow EnableItem and DisableItem.
	onEnabledChanged func()
}

func newNavigator(items ...*container.TabItem) *navigator {
	n := &navigator{items: items, disabled: make(map[*container.TabItem]bool)}
	pages := make([]fyne.CanvasObject, len(items))
	for i, item := range items {
		pages[i] = item.Content
	}
	n.stack = container.New(stackLayout{}, pages...)
	n.show(0)
	return n
}

func (n *navigator) show(index int) {
	n.selected = index
	for i, page := range n.stack.Objects {
		setVisible(page, i == index)
	}
	relayout(n.stack)
}

// setPage shows page for item instead of item.Content, which stays the
// object other code holds on to.
func (n *navigator) setPage(item *container.TabItem, page fyne.CanvasObject) {
	for i, it := range n.items {
		if it == item {
			n.stack.Objects[i] = page
			item.Content.Show() // the page, not its content, is what gets hidden
		}
	}
	n.show(n.selected)
}

// Select shows the page. A disabled page cannot be selected.
func (n *navigator) Select(item *container.TabItem) {
	for i, it := range n.items {
		if it == item && i != n.selected && !n.disabled[item] {
			// Keyboard focus would leave the activated button ringed on a page it no longer controls.
			if c := fyne.CurrentApp().Driver().CanvasForObject(n.stack); c != nil {
				c.Unfocus()
			}
			n.show(i)
			if n.OnSelected != nil {
				n.OnSelected(item)
			}
			return
		}
	}
}

func (n *navigator) Selected() *container.TabItem { return n.items[n.selected] }

func (n *navigator) Disabled(item *container.TabItem) bool { return n.disabled[item] }

func (n *navigator) setDisabled(item *container.TabItem, disabled bool) {
	n.disabled[item] = disabled
	if n.onEnabledChanged != nil {
		n.onEnabledChanged()
	}
}

func (n *navigator) EnableItem(item *container.TabItem)  { n.setDisabled(item, false) }
func (n *navigator) DisableItem(item *container.TabItem) { n.setDisabled(item, true) }

// subTabs is a row of text buttons over a navigator, for the pages inside
// Diagnostics. Unlike AppTabs its buttons keep the normal theme, take keyboard
// focus, and look like the toolbar's.
type subTabs struct {
	nav     *navigator
	buttons []*button
	view    fyne.CanvasObject
}

func newSubTabs(items ...*container.TabItem) *subTabs {
	t := &subTabs{nav: newNavigator(items...)}
	var row []fyne.CanvasObject
	for _, item := range items {
		b := newButton(item.Text, nil, func() { t.nav.Select(item); t.sync() })
		t.buttons = append(t.buttons, b)
		row = append(row, b)
	}
	t.sync()
	t.view = border(container.NewHBox(row...), nil, nil, nil, t.nav.stack)
	return t
}

func (t *subTabs) sync() {
	for i, b := range t.buttons {
		b.Importance = widget.LowImportance
		if i == t.nav.selected {
			b.Importance = widget.MediumImportance
		}
		b.Refresh()
	}
}

func (t *subTabs) SelectIndex(i int) { t.nav.Select(t.nav.items[i]); t.sync() }

// stepper is the horizontal progress strip over the four guide steps. It
// shows position only: it never claims a step is complete.
type stepper struct {
	widget.BaseWidget
	buttons []*button
	current int // -1 on a tool page
}

func newStepper(steps []*container.TabItem, onTap func(*container.TabItem)) *stepper {
	st := &stepper{current: -1}
	for i, step := range steps {
		b := newButton(step.Text, numberIcon(i+1, numberGrey), func() { onTap(step) })
		b.Importance = widget.LowImportance
		st.buttons = append(st.buttons, b)
	}
	st.ExtendBaseWidget(st)
	return st
}

func (st *stepper) setCurrent(index int) {
	st.current = index
	st.Refresh()
}

func (st *stepper) CreateRenderer() fyne.WidgetRenderer {
	r := &stepperRenderer{st: st, underline: canvas.NewRectangle(color.Transparent)}
	for range st.buttons[1:] {
		r.connectors = append(r.connectors, canvas.NewRectangle(color.Transparent))
	}
	r.paint()
	return r
}

type stepperRenderer struct {
	st         *stepper
	connectors []*canvas.Rectangle
	underline  *canvas.Rectangle
}

const (
	stepperRule         = 2
	stepperMinConnector = 12
)

func (r *stepperRenderer) paint() {
	st := r.st
	for i, b := range st.buttons {
		fill := color.Color(numberGrey)
		if i == st.current { // position, never completion
			fill = theme.Color(theme.ColorNamePrimary)
		}
		b.SetIcon(numberIcon(i+1, fill))
	}
	for _, c := range r.connectors {
		c.FillColor = theme.Color(theme.ColorNameSeparator)
		c.Refresh()
	}
	r.underline.FillColor = theme.Color(theme.ColorNamePrimary)
	setVisible(r.underline, st.current >= 0)
	r.underline.Refresh()
}

func (r *stepperRenderer) rowHeight() float32 {
	var height float32
	for _, b := range r.st.buttons {
		height = fyne.Max(height, b.MinSize().Height)
	}
	return height
}

func (r *stepperRenderer) Layout(size fyne.Size) {
	var buttonsWidth float32
	for _, b := range r.st.buttons {
		buttonsWidth += b.MinSize().Width
	}
	margin, height := theme.Padding()*2, r.rowHeight()
	connector := fyne.Max((size.Width-buttonsWidth-float32(len(r.connectors))*2*margin)/float32(len(r.connectors)), stepperMinConnector)
	x := float32(0)
	for i, b := range r.st.buttons {
		b.Move(fyne.NewPos(x, 0))
		b.Resize(fyne.NewSize(b.MinSize().Width, height))
		if i == r.st.current {
			r.underline.Move(fyne.NewPos(x, height))
			r.underline.Resize(fyne.NewSize(b.MinSize().Width, stepperRule))
		}
		x += b.MinSize().Width
		if i < len(r.connectors) {
			r.connectors[i].Move(fyne.NewPos(x+margin, (height-stepperRule)/2))
			r.connectors[i].Resize(fyne.NewSize(connector, stepperRule))
			x += connector + 2*margin
		}
	}
}

func (r *stepperRenderer) MinSize() fyne.Size {
	var width float32
	for _, b := range r.st.buttons {
		width += b.MinSize().Width
	}
	width += float32(len(r.connectors)) * (stepperMinConnector + 4*theme.Padding())
	return fyne.NewSize(width, r.rowHeight()+stepperRule)
}

func (r *stepperRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.underline}
	for _, c := range r.connectors {
		objects = append(objects, c)
	}
	for _, b := range r.st.buttons {
		objects = append(objects, b)
	}
	return objects
}

func (r *stepperRenderer) Refresh() {
	r.paint()
	r.Layout(r.st.Size())
}

func (r *stepperRenderer) Destroy() {}
