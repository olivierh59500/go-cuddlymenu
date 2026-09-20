package loader

import "testing"

func TestCountdownBoundaries(t *testing.T) {
	l := &Screen{spec: Spec{Sector: 49, Blipps: 162}}
	for _, tt := range []struct {
		tick, sector, blipps int
		blank, done          bool
	}{{0, 49, 162, false, false}, {48, 1, 162, false, false}, {49, 0, 161, false, false}, {234, 0, -24, false, false}, {235, 0, -25, true, false}, {247, 0, -37, true, false}, {248, 0, -38, true, true}} {
		l.tick = tt.tick
		s, b := l.Counters()
		if s != tt.sector || b != tt.blipps || l.Blank() != tt.blank || l.Done() != tt.done {
			t.Fatalf("tick %d: counters %d/%d, blank %v, done %v", tt.tick, s, b, l.Blank(), l.Done())
		}
	}
	l.tick = 169
	if l.Volume() != 1 {
		t.Fatal("fade started early")
	}
	l.tick = 244
	if l.Volume() != 0 {
		t.Fatal("fade failed to finish")
	}
}
