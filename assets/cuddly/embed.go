// Package cuddly exposes the archived Cuddly artwork, sound and source data.
package cuddly

import "embed"

// Files contains original assets. The the original implementation reference is archived separately
// and is never needed or executed by the native Go application.
//
//go:embed 3d_doc big_sprite colorshock2 digi dna_demo ehh fullscreen intro led_scroller megaball megascroller menu reset spreadpoint starwars tex data.json ym/*.ym
var Files embed.FS
