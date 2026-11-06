//go:build !headless

package gui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Colours are read from the theme in each Refresh, so a light or dark switch
// repaints these custom pieces too.
func tint(name fyne.ThemeColorName, percent int) color.Color {
	return withAlpha(color.NRGBAModel.Convert(theme.Color(name)).(color.NRGBA), percent)
}

func colorHex(c color.Color) string {
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	return fmt.Sprintf("#%02x%02x%02x", n.R, n.G, n.B)
}

func wrappedLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

// helperText is secondary prose: smaller and muted.
func helperText(text string) *widget.Label {
	label := wrappedLabel(text)
	label.Importance = widget.LowImportance
	label.SizeName = theme.SizeNameCaptionText
	return label
}

// tweakTheme overrides some sizes of whatever theme the app currently uses, so
// light and dark and any test variant carry through.
type tweakTheme struct {
	sizes map[fyne.ThemeSizeName]float32
}

func (t tweakTheme) current() fyne.Theme { return fyne.CurrentApp().Settings().Theme() }
func (t tweakTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return t.current().Color(name, variant)
}
func (t tweakTheme) Font(style fyne.TextStyle) fyne.Resource { return t.current().Font(style) }
func (t tweakTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.current().Icon(name)
}
func (t tweakTheme) Size(name fyne.ThemeSizeName) float32 {
	if size, ok := t.sizes[name]; ok {
		return size
	}
	return t.current().Size(name)
}

// flush removes the label inset, so text lines up with the edge of the boxes
// (inputs, cards, buttons) beside it.
func flush(obj fyne.CanvasObject) fyne.CanvasObject {
	return container.NewThemeOverride(obj, tweakTheme{sizes: map[fyne.ThemeSizeName]float32{theme.SizeNameInnerPadding: 0}})
}

// section separates groups; inside a group the VBox spacing is enough.
func section() fyne.CanvasObject { return gap(theme.Padding() * 2) }

// spaced stacks groups with a section gap between them.
func spaced(groups ...fyne.CanvasObject) fyne.CanvasObject {
	var rows []fyne.CanvasObject
	for i, g := range groups {
		if i > 0 {
			rows = append(rows, section())
		}
		rows = append(rows, g)
	}
	return container.NewVBox(rows...)
}

// pageTitle is the heading each page opens with.
func pageTitle(text string) *widget.Label {
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	label.Wrapping = fyne.TextWrapWord
	label.SizeName = theme.SizeNameSubHeadingText
	return label
}

func pageHeader(title, help string) fyne.CanvasObject {
	if help == "" {
		return flush(pageTitle(title))
	}
	return container.NewVBox(flush(pageTitle(title)), flush(helperText(help)))
}

// field puts the label above its input, with optional help beneath.
func field(label string, input fyne.CanvasObject, help string) fyne.CanvasObject {
	rows := []fyne.CanvasObject{flush(widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})), input}
	if help != "" {
		rows = append(rows, flush(helperText(help)))
	}
	return container.NewVBox(rows...)
}

// scrolled keeps controls clear of the scroll bar that overlays the right edge,
// and inset enough that a focus ring is not clipped at the top or left.
func scrolled(content fyne.CanvasObject) *container.Scroll {
	ring := float32(focusRingSpace + focusRingWidth)
	return container.NewVScroll(container.New(layout.NewCustomPaddedLayout(ring, 0, ring, theme.ScrollBarSize()+theme.Padding()), content))
}

func gap(height float32) fyne.CanvasObject {
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(1, height))
	return spacer
}

// maxWidth centres one object in a column of at most max width, keeping a
// margin at either side and a little air above it.
type maxWidth struct{ max, margin, top float32 }

func (m maxWidth) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	width := fyne.Min(m.max, size.Width-2*m.margin)
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		o.Resize(fyne.NewSize(width, size.Height-m.top))
		o.Move(fyne.NewPos((size.Width-width)/2, m.top))
	}
}

func (m maxWidth) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var minimum fyne.Size
	top := float32(0)
	for _, o := range objects {
		if o.Visible() {
			minimum, top = minimum.Max(o.MinSize()), m.top
		}
	}
	return fyne.NewSize(minimum.Width+2*m.margin, minimum.Height+top)
}

// callout is a label on a rounded warning-tinted box. Embedding Label keeps
// Text, SetText, Show and Hide working for callers.
type callout struct{ widget.Label }

func newCallout(text string) *callout {
	c := &callout{}
	c.Text, c.Wrapping = text, fyne.TextWrapWord
	c.ExtendBaseWidget(c)
	return c
}

func (c *callout) CreateRenderer() fyne.WidgetRenderer {
	r := &calloutRenderer{box: canvas.NewRectangle(color.Transparent), label: c.Label.CreateRenderer()}
	r.paint()
	return r
}

type calloutRenderer struct {
	box   *canvas.Rectangle
	label fyne.WidgetRenderer
}

func (r *calloutRenderer) paint() {
	r.box.FillColor = tint(theme.ColorNameWarning, 14)
	r.box.StrokeColor = tint(theme.ColorNameWarning, 45)
	r.box.StrokeWidth = 1
	r.box.CornerRadius = theme.Size(theme.SizeNameInputRadius)
	r.box.Refresh()
}

func (r *calloutRenderer) Layout(size fyne.Size) {
	r.box.Resize(size)
	r.label.Layout(size)
}
func (r *calloutRenderer) MinSize() fyne.Size { return r.label.MinSize() }
func (r *calloutRenderer) Objects() []fyne.CanvasObject {
	return append([]fyne.CanvasObject{r.box}, r.label.Objects()...)
}
func (r *calloutRenderer) Refresh() { r.paint(); r.label.Refresh() }
func (r *calloutRenderer) Destroy() { r.label.Destroy() }

// surface paints a themed fill behind its content, for the chrome layer.
type surface struct {
	widget.BaseWidget
	content fyne.CanvasObject
	fill    fyne.ThemeColorName
}

func newSurface(content fyne.CanvasObject, fill fyne.ThemeColorName) *surface {
	s := &surface{content: content, fill: fill}
	s.ExtendBaseWidget(s)
	return s
}

func (s *surface) CreateRenderer() fyne.WidgetRenderer {
	r := &surfaceRenderer{s: s, box: canvas.NewRectangle(theme.Color(s.fill))}
	return r
}

type surfaceRenderer struct {
	s   *surface
	box *canvas.Rectangle
}

func (r *surfaceRenderer) Layout(size fyne.Size) { r.box.Resize(size); r.s.content.Resize(size) }
func (r *surfaceRenderer) MinSize() fyne.Size    { return r.s.content.MinSize() }
func (r *surfaceRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.box, r.s.content}
}
func (r *surfaceRenderer) Refresh() {
	r.box.FillColor = theme.Color(r.s.fill)
	r.box.Refresh()
	r.s.content.Refresh()
}
func (r *surfaceRenderer) Destroy() {}

// panel is a rounded, bordered surface around its content.
type panel struct {
	widget.BaseWidget
	content fyne.CanvasObject
}

func (s *surface) inner() fyne.CanvasObject { return s.content }
func (p *panel) inner() fyne.CanvasObject   { return p.content }

func newPanel(content fyne.CanvasObject) *panel {
	p := &panel{content: content}
	p.ExtendBaseWidget(p)
	return p
}

func (p *panel) CreateRenderer() fyne.WidgetRenderer {
	r := &panelRenderer{p: p, box: canvas.NewRectangle(color.Transparent)}
	r.paint()
	return r
}

type panelRenderer struct {
	p   *panel
	box *canvas.Rectangle
}

func (r *panelRenderer) paint() {
	r.box.FillColor = theme.Color(theme.ColorNameInputBackground)
	r.box.StrokeColor = theme.Color(theme.ColorNameSeparator)
	r.box.StrokeWidth = theme.Size(theme.SizeNameSeparatorThickness)
	r.box.CornerRadius = theme.Size(theme.SizeNameCardRadius)
	r.box.Refresh()
}

func (r *panelRenderer) Layout(size fyne.Size) {
	r.box.Resize(size)
	inset := theme.Padding() * 2
	r.p.content.Move(fyne.NewPos(inset, inset))
	r.p.content.Resize(size.Subtract(fyne.NewSize(2*inset, 2*inset)))
}

func (r *panelRenderer) MinSize() fyne.Size {
	return r.p.content.MinSize().Add(fyne.NewSize(4*theme.Padding(), 4*theme.Padding()))
}
func (r *panelRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.box, r.p.content} }
func (r *panelRenderer) Refresh()                     { r.paint(); r.p.content.Refresh() }
func (r *panelRenderer) Destroy()                     {}

// choiceCard is a large selectable option. It is tappable, keyboard
// operable (Space or Enter) and inert while disabled.
type choiceCard struct {
	widget.BaseWidget
	title, description                   string
	selected, disabled, hovered, focused bool
	onTap                                func()
}

func newChoiceCard(title, description string, onTap func()) *choiceCard {
	c := &choiceCard{title: title, description: description, onTap: onTap}
	c.ExtendBaseWidget(c)
	return c
}

func (c *choiceCard) setSelected(on bool) {
	if c.selected != on {
		c.selected = on
		c.Refresh()
	}
}

func (c *choiceCard) activate() {
	if !c.disabled && c.onTap != nil {
		c.onTap()
	}
}

func (c *choiceCard) Tapped(*fyne.PointEvent) {
	if c.disabled {
		return
	}
	// The driver focuses a clicked control. Drop it, as Button.Tapped does, so the ring shows keyboard focus only.
	if owner := fyne.CurrentApp().Driver().CanvasForObject(c); owner != nil {
		owner.Unfocus()
	}
	c.activate()
}

func (c *choiceCard) TypedRune(rune) {}

func (c *choiceCard) TypedKey(e *fyne.KeyEvent) {
	switch e.Name {
	case fyne.KeySpace, fyne.KeyReturn, fyne.KeyEnter:
		c.activate()
	}
}

func (c *choiceCard) FocusGained() { c.focused = true; c.Refresh() }
func (c *choiceCard) FocusLost()   { c.focused = false; c.Refresh() }

func (c *choiceCard) MouseIn(*desktop.MouseEvent)    { c.hovered = true; c.Refresh() }
func (c *choiceCard) MouseMoved(*desktop.MouseEvent) {}
func (c *choiceCard) MouseOut()                      { c.hovered = false; c.Refresh() }

func (c *choiceCard) Cursor() desktop.Cursor {
	if c.disabled {
		return desktop.DefaultCursor
	}
	return desktop.PointerCursor
}

func (c *choiceCard) Disabled() bool { return c.disabled }
func (c *choiceCard) Enable()        { c.disabled = false; c.Refresh() }
func (c *choiceCard) Disable()       { c.disabled, c.hovered = true, false; c.Refresh() }

func (c *choiceCard) CreateRenderer() fyne.WidgetRenderer {
	r := &choiceCardRenderer{
		c: c, fill: canvas.NewRectangle(color.Transparent), ring: canvas.NewRectangle(color.Transparent), border: canvas.NewRectangle(color.Transparent),
		title: widget.NewLabelWithStyle(c.title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), description: helperText(c.description),
	}
	r.paint()
	return r
}

type choiceCardRenderer struct {
	c                  *choiceCard
	fill, ring, border *canvas.Rectangle
	title, description *widget.Label
}

func (r *choiceCardRenderer) paint() {
	c, radius := r.c, theme.Size(theme.SizeNameCardRadius)
	fill, edge, width := color.Color(color.Transparent), theme.Color(theme.ColorNameInputBorder), float32(1)
	switch {
	case c.selected:
		fill, edge, width = tint(theme.ColorNamePrimary, 8), theme.Color(theme.ColorNamePrimary), 2
		if c.disabled {
			edge = tint(theme.ColorNamePrimary, 45)
		}
	case c.disabled:
		edge = theme.Color(theme.ColorNameSeparator)
	case c.hovered:
		fill = theme.Color(theme.ColorNameHover)
	}
	r.fill.FillColor, r.fill.CornerRadius = fill, radius
	r.border.StrokeColor, r.border.StrokeWidth, r.border.CornerRadius = edge, width, radius
	r.ring.StrokeColor, r.ring.StrokeWidth, r.ring.CornerRadius = theme.Color(theme.ColorNameForeground), 2, radius+focusRingGap
	setVisible(r.ring, c.focused && !c.disabled)
	r.title.Importance = widget.MediumImportance
	if c.disabled {
		r.title.Importance = widget.LowImportance
	}
	for _, o := range []fyne.CanvasObject{r.fill, r.ring, r.border, r.title} {
		o.Refresh()
	}
	r.description.Refresh()
}

const focusRingGap = 3

func (r *choiceCardRenderer) inset() float32 { return theme.InnerPadding() + theme.Padding() }

func (r *choiceCardRenderer) Layout(size fyne.Size) {
	for _, box := range []*canvas.Rectangle{r.fill, r.border} {
		box.Resize(size)
	}
	// The ring sits outside the border so it never reads as the selected edge.
	r.ring.Move(fyne.NewPos(-focusRingGap, -focusRingGap))
	r.ring.Resize(size.AddWidthHeight(2*focusRingGap, 2*focusRingGap))
	inset := r.inset()
	width := fyne.Max(size.Width-2*inset, 0)
	titleHeight := r.title.MinSize().Height
	r.title.Move(fyne.NewPos(inset, inset))
	r.title.Resize(fyne.NewSize(width, titleHeight))
	// Wrapped text reports its height for the width it was last given.
	r.description.Resize(fyne.NewSize(width, r.description.Size().Height))
	r.description.Move(fyne.NewPos(inset, inset+titleHeight))
	r.description.Resize(fyne.NewSize(width, r.description.MinSize().Height))
}

func (r *choiceCardRenderer) MinSize() fyne.Size {
	inset := r.inset()
	return fyne.NewSize(r.title.MinSize().Width+2*inset, 2*inset+r.title.MinSize().Height+r.description.MinSize().Height)
}

func (r *choiceCardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.fill, r.ring, r.border, r.title, r.description}
}
func (r *choiceCardRenderer) Refresh() { r.paint() }
func (r *choiceCardRenderer) Destroy() {}

// digitStrokes draws 1 to 4 inside the 24 by 24 number icon. SVG text is not
// rendered by Fyne, so the digits are paths.
var digitStrokes = map[int]string{
	1: "M10 9.4 L12.6 7 V17",
	2: "M9.3 9.6 C9.6 7.9 10.8 7 12.1 7 C13.7 7 14.8 8 14.8 9.4 C14.8 11 13.5 12 12 13.4 L9.2 17 H15.2",
	3: "M9.4 8.2 C10 7.5 11 7 12 7 C13.5 7 14.6 7.9 14.6 9.2 C14.6 10.6 13.4 11.4 11.8 11.5 C13.6 11.5 14.9 12.4 14.9 14 C14.9 15.8 13.6 17 12 17 C10.7 17 9.7 16.5 9.1 15.6",
	4: "M13.6 17 V7 L8.6 13.4 H15.6",
}

// numberIcon is a filled circle holding n, with explicit colours so it reads
// in both light and dark variants.
func numberIcon(n int, fill color.Color) fyne.Resource {
	hex := colorHex(fill)
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><circle cx="12" cy="12" r="12" fill="%s"/><path d="%s" fill="none" stroke="#fafbff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`, hex, digitStrokes[n])
	return fyne.NewStaticResource(fmt.Sprintf("number-%d-%s.svg", n, hex[1:]), []byte(svg))
}

var numberGrey = rgb(0x6b7280)

// instructionRow is a muted digit beside wrapped text, with an optional
// control at the trailing edge. No circle: the stepper owns those numbers.
func instructionRow(n int, text string, trailing fyne.CanvasObject) fyne.CanvasObject {
	digit := widget.NewLabel(strconv.Itoa(n))
	digit.Importance = widget.LowImportance
	// The text keeps its inner padding, so give the flush digit the same top offset.
	aligned := container.New(layout.NewCustomPaddedLayout(theme.InnerPadding(), 0, 0, 0), flush(digit))
	return border(nil, nil, aligned, trailing, wrappedLabel(text))
}

// border is container.NewBorder with the objects in visual order. Fyne walks
// focus in container order and NewBorder lists the centre first, which would
// focus a page before its header.
func border(top, bottom, left, right, centre fyne.CanvasObject) *fyne.Container {
	var objects []fyne.CanvasObject
	for _, o := range []fyne.CanvasObject{top, left, centre, right, bottom} {
		if o != nil {
			objects = append(objects, o)
		}
	}
	return container.New(layout.NewBorderLayout(top, bottom, left, right), objects...)
}

// hint is a centred, wrapping line of guidance for an empty area. Unlike
// NewCenter it gives the label the whole width, so it wraps instead of
// running past its pane.
func hint(label *widget.Label) fyne.CanvasObject {
	label.Alignment = fyne.TextAlignCenter
	return container.New(layout.NewCustomPaddedLayout(theme.Padding()*6, 0, theme.Padding()*2, theme.Padding()*2), label)
}

// dialogBody pulls a dialog's content out by the dialog's own inset, so its
// text lines up with the title above it.
func dialogBody(content fyne.CanvasObject) fyne.CanvasObject {
	return container.New(layout.NewCustomPaddedLayout(0, 0, -theme.Padding()*2, -theme.Padding()*2), content)
}

// button is a widget.Button with a focus ring. Fyne marks focus by blending
// the theme Focus colour into the fill, which cannot show on a primary button.
// Clicks drop focus (Button.Tapped), so the ring shows keyboard focus only.
type button struct {
	widget.Button
	focused bool
}

func newButton(text string, icon fyne.Resource, tapped func()) *button {
	b := &button{}
	b.Text, b.Icon, b.OnTapped = text, icon, tapped
	b.ExtendBaseWidget(b)
	return b
}

func (b *button) FocusGained() { b.focused = true; b.Button.FocusGained() }
func (b *button) FocusLost()   { b.focused = false; b.Button.FocusLost() }

// focusRingColour is the one ring colour: the foreground contrasts with every
// surface a control sits on, which no fill-relative colour does.
func focusRingColour() color.Color { return theme.Color(theme.ColorNameForeground) }

const (
	focusRingWidth = 2
	focusRingSpace = 1 // between the control and its ring
)

func (b *button) CreateRenderer() fyne.WidgetRenderer {
	return newRingRenderer(b.Button.CreateRenderer(), func() bool { return b.focused && !b.Disabled() }, theme.Size(theme.SizeNameButtonRadius))
}

// entry is a widget.Entry with the same focus ring as button. Fyne's own focus
// cue is a 1 px border colour change, which is under 2:1 in light and 1.06:1 in dark.
type entry struct {
	widget.Entry
	focused bool
}

func newEntry() *entry {
	e := &entry{}
	e.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
	e.ExtendBaseWidget(e)
	return e
}

// newPasswordEntry has its own reveal button, because Fyne's is unexported.
func newPasswordEntry() *entry {
	e := newEntry()
	e.Password = true
	reveal := newButton("", theme.VisibilityIcon(), nil)
	reveal.Importance = widget.LowImportance
	reveal.OnTapped = func() {
		// As Fyne's own revealer: inert while disabled, and focus returns to the entry.
		if e.Disabled() {
			return
		}
		e.Password = !e.Password
		reveal.SetIcon(theme.VisibilityIcon())
		if e.Password {
			reveal.SetIcon(theme.VisibilityOffIcon())
		}
		e.Refresh()
		if c := fyne.CurrentApp().Driver().CanvasForObject(e); c != nil {
			c.Focus(e)
		}
	}
	reveal.SetIcon(theme.VisibilityOffIcon())
	e.ActionItem = reveal
	return e
}

func (e *entry) FocusGained() { e.focused = true; e.Entry.FocusGained() }
func (e *entry) FocusLost()   { e.focused = false; e.Entry.FocusLost() }

func (e *entry) CreateRenderer() fyne.WidgetRenderer {
	return newRingRenderer(e.Entry.CreateRenderer(), func() bool { return e.focused && !e.Disabled() }, theme.Size(theme.SizeNameInputRadius))
}

// ringRenderer wraps a control's renderer and draws a ring just outside it.
type ringRenderer struct {
	inner  fyne.WidgetRenderer
	ring   *canvas.Rectangle
	shown  func() bool
	radius float32
}

func newRingRenderer(inner fyne.WidgetRenderer, shown func() bool, radius float32) *ringRenderer {
	r := &ringRenderer{inner: inner, ring: canvas.NewRectangle(color.Transparent), shown: shown, radius: radius}
	r.paint()
	return r
}

func (r *ringRenderer) paint() {
	r.ring.StrokeColor, r.ring.StrokeWidth = focusRingColour(), focusRingWidth
	r.ring.CornerRadius = r.radius + focusRingSpace + focusRingWidth/2
	setVisible(r.ring, r.shown())
	r.ring.Refresh()
}

func (r *ringRenderer) Layout(size fyne.Size) {
	r.inner.Layout(size)
	out := float32(focusRingSpace + focusRingWidth)
	r.ring.Move(fyne.NewPos(-out, -out))
	r.ring.Resize(size.AddWidthHeight(2*out, 2*out))
}
func (r *ringRenderer) MinSize() fyne.Size { return r.inner.MinSize() }
func (r *ringRenderer) Objects() []fyne.CanvasObject {
	return append(r.inner.Objects(), r.ring)
}
func (r *ringRenderer) Refresh() { r.inner.Refresh(); r.paint() }
func (r *ringRenderer) Destroy() { r.inner.Destroy() }

// slimBar draws a progress bar at a fixed thin height, centred in its row.
type slimBar struct{ height float32 }

func (l slimBar) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(fyne.NewSize(size.Width, l.height))
		o.Move(fyne.NewPos(0, (size.Height-l.height)/2))
	}
}

func (l slimBar) MinSize([]fyne.CanvasObject) fyne.Size { return fyne.NewSize(60, l.height) }

// dialogSize fits a text dialog to its window and text, since wrapped text
// reports no height until it has a width.
func (s *screen) dialogSize(text string) fyne.Size {
	canvas := s.w.Canvas().Size()
	width := fyne.Min(520, canvas.Width-32)
	lines := float32(strings.Count(text, "\n") + 1 + len(text)*8/int(width))
	return fyne.NewSize(width, fyne.Min(canvas.Height-64, 120+lines*fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameText)*1.6))
}

// inform shows a message with one OK button.
func (s *screen) inform(title, text string) {
	var d *dialog.CustomDialog
	ok := newButton("OK", nil, func() { d.Hide() })
	ok.Importance = widget.HighImportance
	d = s.dialogWithButtons(title, container.NewVScroll(flush(archiveLabel(text))), ok)
	d.Resize(s.dialogSize(text))
	d.Show()
}

// confirmDialog shows body with Cancel and one primary confirm button. done
// gets true for the confirm button and false for Cancel.
func (s *screen) confirmDialog(title string, body fyne.CanvasObject, confirm, cancel string, size fyne.Size, done func(bool)) {
	var d *dialog.CustomDialog
	no := newButton(cancel, nil, func() { d.Hide(); done(false) })
	yes := newButton(confirm, nil, func() { d.Hide(); done(true) })
	yes.Importance = widget.HighImportance
	d = s.dialogWithButtons(title, body, no, yes)
	d.Resize(size)
	d.Show()
}

// confirmText is confirmDialog for a plain message.
func (s *screen) confirmText(title, text, confirm string, done func(bool)) {
	s.confirmDialog(title, container.NewVScroll(flush(archiveLabel(text))), confirm, "Cancel", s.dialogSize(text), done)
}

// dialogWithButtons opens a titled dialog with its buttons at natural width on
// the trailing edge. Callers hide it from the buttons they pass.
func (s *screen) dialogWithButtons(title string, body fyne.CanvasObject, buttons ...fyne.CanvasObject) *dialog.CustomDialog {
	footer := container.NewHBox(append([]fyne.CanvasObject{layout.NewSpacer()}, buttons...)...)
	return dialog.NewCustomWithoutButtons(title, dialogBody(border(nil, container.NewVBox(gap(theme.Padding()), footer), nil, nil, body)), s.w)
}

// resultPanel is a titled panel with an optional Copy button beside the title,
// for the text a diagnostic or startup action reports.
func resultPanel(title string, copy *button, text *widget.Label) fyne.CanvasObject {
	var trailing fyne.CanvasObject
	if copy != nil {
		trailing = copy
	}
	header := border(nil, nil, widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), trailing, nil)
	return newPanel(border(header, nil, nil, nil, container.NewVScroll(text)))
}

// relayout re-runs the layouts of o and its nested containers, then repaints.
// Container.Refresh would also refresh every widget below, re-rendering every
// icon on every page, which made each step change and busy toggle slow.
func relayout(o fyne.CanvasObject) {
	var walk func(fyne.CanvasObject)
	walk = func(o fyne.CanvasObject) {
		if w, ok := o.(interface{ inner() fyne.CanvasObject }); ok {
			walk(w.inner())
			return
		}
		c, ok := o.(*fyne.Container)
		if !ok {
			return
		}
		if c.Layout != nil {
			c.Layout.Layout(c.Objects, c.Size())
		}
		for _, child := range c.Objects {
			walk(child)
		}
	}
	walk(o)
	if c := fyne.CurrentApp().Driver().CanvasForObject(o); c != nil {
		c.Refresh(o)
	}
}
