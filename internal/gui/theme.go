//go:build !headless

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// appTheme is a restrained blue-tinted neutral palette with one accent. It
// keeps Fyne's bundled font and icons, and follows the OS light or dark variant.
type appTheme struct{}

func rgb(hex uint32) color.NRGBA {
	return color.NRGBA{R: uint8(hex >> 16), G: uint8(hex >> 8), B: uint8(hex), A: 0xff}
}

func withAlpha(c color.NRGBA, percent int) color.NRGBA {
	c.A = uint8(percent * 255 / 100)
	return c
}

func palette(light bool) map[fyne.ThemeColorName]color.Color {
	if light {
		foreground, primary := rgb(0x1b1e25), rgb(0x2c5cdd)
		onStatus := rgb(0xfafbff)
		return map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:          rgb(0xf6f7f9),
			theme.ColorNameOverlayBackground:   rgb(0xf6f7f9),
			theme.ColorNameHeaderBackground:    rgb(0xfdfdfe),
			theme.ColorNameMenuBackground:      rgb(0xfdfdfe),
			theme.ColorNameInputBackground:     rgb(0xfdfdfe),
			theme.ColorNameForeground:          foreground,
			theme.ColorNameDisabled:            rgb(0x5c6372),
			theme.ColorNamePlaceHolder:         rgb(0x666d7c),
			theme.ColorNameSeparator:           rgb(0xe2e4e9),
			theme.ColorNameInputBorder:         rgb(0x8b92a1),
			theme.ColorNameButton:              rgb(0xe8eaef),
			theme.ColorNameDisabledButton:      rgb(0xeceef2),
			theme.ColorNameHover:               withAlpha(foreground, 6),
			theme.ColorNamePressed:             withAlpha(foreground, 12),
			theme.ColorNamePrimary:             primary,
			theme.ColorNameForegroundOnPrimary: rgb(0xfafbff),
			theme.ColorNameFocus:               withAlpha(primary, 73),
			theme.ColorNameSelection:           withAlpha(primary, 25),
			theme.ColorNameHyperlink:           rgb(0x2452c4),
			theme.ColorNameSuccess:             rgb(0x17794a),
			theme.ColorNameWarning:             rgb(0x9a5b07),
			theme.ColorNameError:               rgb(0xc02d39),
			theme.ColorNameForegroundOnSuccess: onStatus,
			theme.ColorNameForegroundOnWarning: onStatus,
			theme.ColorNameForegroundOnError:   onStatus,
			theme.ColorNameShadow:              withAlpha(rgb(0x000000), 16),
			theme.ColorNameScrollBar:           withAlpha(foreground, 30),
		}
	}
	foreground, primary := rgb(0xe6e8ed), rgb(0x3a68e0)
	onStatus := rgb(0x14161b)
	return map[fyne.ThemeColorName]color.Color{
		theme.ColorNameBackground:          rgb(0x14161b),
		theme.ColorNameOverlayBackground:   rgb(0x14161b),
		theme.ColorNameHeaderBackground:    rgb(0x1a1c22),
		theme.ColorNameMenuBackground:      rgb(0x1a1c22),
		theme.ColorNameInputBackground:     rgb(0x1d2027),
		theme.ColorNameForeground:          foreground,
		theme.ColorNameDisabled:            rgb(0x9aa0ad),
		theme.ColorNamePlaceHolder:         rgb(0x8f96a4),
		theme.ColorNameSeparator:           rgb(0x2a2d35),
		theme.ColorNameInputBorder:         rgb(0x666c7a),
		theme.ColorNameButton:              rgb(0x282b33),
		theme.ColorNameDisabledButton:      rgb(0x22252b),
		theme.ColorNameHover:               withAlpha(foreground, 7),
		theme.ColorNamePressed:             withAlpha(foreground, 13),
		theme.ColorNamePrimary:             primary,
		theme.ColorNameForegroundOnPrimary: rgb(0xfafbff),
		theme.ColorNameFocus:               withAlpha(primary, 88),
		theme.ColorNameSelection:           withAlpha(primary, 35),
		theme.ColorNameHyperlink:           rgb(0x86a8f5),
		theme.ColorNameSuccess:             rgb(0x4cc38a),
		theme.ColorNameWarning:             rgb(0xe2a33d),
		theme.ColorNameError:               rgb(0xef6a74),
		theme.ColorNameForegroundOnSuccess: onStatus,
		theme.ColorNameForegroundOnWarning: onStatus,
		theme.ColorNameForegroundOnError:   onStatus,
		theme.ColorNameShadow:              withAlpha(rgb(0x000000), 40),
		theme.ColorNameScrollBar:           withAlpha(foreground, 30),
	}
}

var (
	lightPalette = palette(true)
	darkPalette  = palette(false)
)

func (appTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	colors := darkPalette
	if variant == theme.VariantLight {
		colors = lightPalette
	}
	if c, ok := colors[name]; ok {
		return c
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (appTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameSubHeadingText:
		return 21
	case theme.SizeNameHeadingText:
		return 26
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameInputBorder, theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameSplitThickness:
		return 4 // Fyne also uses this as the drag target
	case theme.SizeNameInputRadius:
		return 6 // a tighter corner leaves less of the scroll shadow outside it
	case theme.SizeNameButtonRadius, theme.SizeNamePopupRadius, theme.SizeNameMenuRadius:
		return 8
	case theme.SizeNameSelectionRadius:
		return 6
	case theme.SizeNameCardRadius, theme.SizeNameDialogRadius:
		return 12
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameScrollBarRadius:
		return 5
	}
	return theme.DefaultTheme().Size(name)
}

func (appTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(style) }

func (appTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(name) }
