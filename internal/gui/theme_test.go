//go:build !headless

package gui

import (
	"image/color"
	"math"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func relativeLuminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	linear := func(v uint32) float64 {
		s := float64(v) / 0xffff
		if s <= 0.03928 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	return 0.2126*linear(r) + 0.7152*linear(g) + 0.0722*linear(b)
}

// contrast is the WCAG 2.x contrast ratio between two opaque colours.
func contrast(a, b color.Color) float64 {
	la, lb := relativeLuminance(a), relativeLuminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func TestThemeColoursMeetContrastInBothVariants(t *testing.T) {
	for _, variant := range []struct {
		name string
		v    fyne.ThemeVariant
	}{{"light", theme.VariantLight}, {"dark", theme.VariantDark}} {
		t.Run(variant.name, func(t *testing.T) {
			colorOf := func(name fyne.ThemeColorName) color.Color { return appTheme{}.Color(name, variant.v) }
			for _, pair := range []struct {
				fg, bg  fyne.ThemeColorName
				minimum float64
			}{
				{theme.ColorNameForeground, theme.ColorNameBackground, 4.5},
				{theme.ColorNameForeground, theme.ColorNameInputBackground, 4.5},
				{theme.ColorNameDisabled, theme.ColorNameBackground, 4.5},
				{theme.ColorNameDisabled, theme.ColorNameInputBackground, 4.5},
				{theme.ColorNamePlaceHolder, theme.ColorNameBackground, 4.5},
				{theme.ColorNamePlaceHolder, theme.ColorNameInputBackground, 4.5},
				{theme.ColorNameHyperlink, theme.ColorNameBackground, 4.5},
				{theme.ColorNameForegroundOnPrimary, theme.ColorNamePrimary, 4.5},
				{theme.ColorNameForegroundOnSuccess, theme.ColorNameSuccess, 4.5},
				{theme.ColorNameForegroundOnWarning, theme.ColorNameWarning, 4.5},
				{theme.ColorNameForegroundOnError, theme.ColorNameError, 4.5},
				{theme.ColorNameInputBorder, theme.ColorNameInputBackground, 3},
				{theme.ColorNamePrimary, theme.ColorNameBackground, 3},
			} {
				if got := contrast(colorOf(pair.fg), colorOf(pair.bg)); got < pair.minimum {
					t.Errorf("%s on %s has contrast %.2f, want at least %.1f", pair.fg, pair.bg, got, pair.minimum)
				}
			}
		})
	}
}

// over composites the translucent colour top on the opaque bottom.
func over(top, bottom color.Color) color.Color {
	tr, tg, tb, ta := top.RGBA()
	br, bg, bb, _ := bottom.RGBA()
	mix := func(t, b uint32) uint8 { return uint8((t + b*(0xffff-ta)/0xffff) >> 8) }
	return color.NRGBA{R: mix(tr, br), G: mix(tg, bg), B: mix(tb, bb), A: 0xff}
}

func TestFocusIsVisibleAndKeepsTextReadable(t *testing.T) {
	for _, variant := range []struct {
		name string
		v    fyne.ThemeVariant
	}{{"light", theme.VariantLight}, {"dark", theme.VariantDark}} {
		t.Run(variant.name, func(t *testing.T) {
			test.NewTempApp(t).Settings().SetTheme(variantTheme{variant: variant.v})
			colorOf := func(name fyne.ThemeColorName) color.Color { return appTheme{}.Color(name, variant.v) }
			// Other widgets blend the Focus colour into the page.
			focused := over(colorOf(theme.ColorNameFocus), colorOf(theme.ColorNameBackground))
			if got := contrast(focused, colorOf(theme.ColorNameBackground)); got < 3 {
				t.Errorf("focus against background %.2f, want at least 3", got)
			}
			// The ring is drawn outside the control, so it must contrast with every surface it can sit on.
			ring := focusRingColour()
			for _, surface := range []fyne.ThemeColorName{theme.ColorNameBackground, theme.ColorNameHeaderBackground, theme.ColorNameInputBackground, theme.ColorNameOverlayBackground} {
				if got := contrast(ring, colorOf(surface)); got < 3 {
					t.Errorf("focus ring against %s %.2f, want at least 3", surface, got)
				}
			}
			// A button also blends Focus into its fill, so the label must stay readable.
			for _, b := range []struct {
				name       string
				fill, text fyne.ThemeColorName
			}{
				{"primary", theme.ColorNamePrimary, theme.ColorNameForegroundOnPrimary},
				{"medium", theme.ColorNameButton, theme.ColorNameForeground},
				{"low on page", theme.ColorNameBackground, theme.ColorNameForeground},
				{"low on header", theme.ColorNameHeaderBackground, theme.ColorNameForeground},
			} {
				blended := over(colorOf(theme.ColorNameFocus), colorOf(b.fill))
				if got := contrast(colorOf(b.text), blended); got < 4.5 {
					t.Errorf("%s: label on the focused fill %.2f, want at least 4.5", b.name, got)
				}
			}
		})
	}
}

func TestButtonAndEntryShowKeyboardFocusNotClickFocus(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.tabs.Select(s.diagnosticsTab)
	b := s.diagnosticsUI.copy // Copy does not navigate, so only Button.Tapped can drop focus
	ringOf := func(w fyne.Widget) *canvas.Rectangle {
		return test.WidgetRenderer(w).(*ringRenderer).ring
	}
	if ringOf(b).Visible() {
		t.Fatal("ring shows before focus")
	}
	s.w.Canvas().Focus(b)
	if !ringOf(b).Visible() {
		t.Fatal("keyboard focus shows no ring")
	}
	b.Disable()
	if ringOf(b).Visible() {
		t.Fatal("a disabled button shows a ring")
	}
	b.Enable()
	test.Tap(b)
	if ringOf(b).Visible() || s.w.Canvas().Focused() != nil {
		t.Fatal("a click left the ring on")
	}
	s.tabs.Select(s.authorizeTab)
	s.w.Canvas().Focus(s.email)
	if !ringOf(s.email).Visible() || ringOf(s.email).Position().X >= 0 {
		t.Fatal("a focused entry shows no ring outside its bounds")
	}
	s.w.Canvas().Unfocus()
	if ringOf(s.email).Visible() {
		t.Fatal("an entry kept its ring after losing focus")
	}
}

func TestDiagnosticsButtonsKeepTheNormalTheme(t *testing.T) {
	for _, variant := range []fyne.ThemeVariant{theme.VariantLight, theme.VariantDark} {
		s, _ := newTestScreen(t, nil)
		fyne.CurrentApp().Settings().SetTheme(variantTheme{variant: variant})
		for i := range s.diagnosticsTabs.nav.items {
			s.diagnosticsTabs.SelectIndex(i)
			s.tabs.Select(s.diagnosticsTab)
		}
		buttons := []*button{s.diagnosticsUI.run, s.startupUI.buttons["preview"]}
		for _, b := range buttons {
			th := b.Theme()
			fill, text := th.Color(theme.ColorNamePrimary, variant), th.Color(theme.ColorNameForegroundOnPrimary, variant)
			if fill != (appTheme{}).Color(theme.ColorNamePrimary, variant) {
				t.Errorf("variant %v: %q is drawn on a remapped primary %v", variant, b.Text, fill)
			}
			if got := contrast(fill, text); got < 4.5 {
				t.Errorf("variant %v: %q text on its fill %.2f, want at least 4.5", variant, b.Text, got)
			}
		}
	}
}

func TestPasswordRevealIsInertWhileDisabled(t *testing.T) {
	test.NewTempApp(t)
	e := newPasswordEntry()
	w := test.NewWindow(e)
	defer w.Close()
	reveal := e.ActionItem.(*button)
	e.Disable()
	test.Tap(reveal)
	if !e.Password {
		t.Fatal("reveal showed a disabled password field")
	}
	e.Enable()
	test.Tap(reveal)
	if e.Password || w.Canvas().Focused() != e {
		t.Fatal("reveal did not show the password and return focus to the field")
	}
}
