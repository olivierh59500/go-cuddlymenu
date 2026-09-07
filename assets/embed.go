// Package assets exposes the game resources from inside every desktop and
// mobile build.
package assets

import "embed"

// Files contains the original Cuddly Demos menu graphics and music.
//
//go:embed menu/*.png menu/*.ym
var Files embed.FS
