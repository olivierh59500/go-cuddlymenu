// Package cuddly embeds the native demo artwork, fonts, music and scene data.
package cuddly

import "embed"

// Files supplies self-contained assets for desktop and mobile playback.
//
//go:embed 3d_doc big_sprite colorshock2 digi dna_demo ehh fullscreen intro led_scroller megaball megascroller menu reset spreadpoint starwars tex data.json loader.json ym/*.ym
var Files embed.FS
