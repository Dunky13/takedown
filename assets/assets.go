// Package assets contains the application artwork embedded into the executable.
package assets

import _ "embed"

// IconPNG is the application icon generated from icon.svg.
//
//go:embed icon.png
var IconPNG []byte
