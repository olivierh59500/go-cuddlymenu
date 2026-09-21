package loader

import (
	"math"
	"testing"
)

func TestCountdownBoundaries(t *testing.T) {
	l := &Screen{spec: Spec{Sector: 49, Blipps: 162}, rate: 60}
	for _, tt := range []struct {
		tick, sector, blipps int
		blank                bool
	}{{0, 49, 162, false}, {48, 1, 162, false}, {49, 0, 161, false}, {234, 0, -24, false}, {235, 0, -25, true}} {
		l.tick = tt.tick
		s, b := l.Counters()
		if s != tt.sector || b != tt.blipps || l.Blank() != tt.blank {
			t.Fatalf("tick %d: counters %d/%d, blank %v", tt.tick, s, b, l.Blank())
		}
	}
}

func TestLoaderDurationsAtBothRates(t *testing.T) {
	for _, rate := range []int{50, 60} {
		l := &Screen{spec: Spec{Sector: 49, Blipps: 162}, rate: rate}
		for l.tick < 169 {
			l.advance()
		}
		if l.Volume() != 1 {
			t.Fatal("fade started early")
		}
		for i := 0; i < rate; i++ {
			l.advance()
		}
		if math.Abs(l.Volume()-1.0/3) > 1e-9 {
			t.Fatalf("%d Hz: fade duration changed: %g", rate, l.Volume())
		}
		for i := 0; i < rate/2; i++ {
			l.advance()
		}
		if l.Volume() > 1e-9 {
			t.Fatal("1.5 second fade did not finish")
		}
		l = &Screen{spec: Spec{Sector: 49, Blipps: 162}, rate: rate, tick: 235}
		ticks := 0
		for !l.Done() {
			l.advance()
			ticks++
		}
		want := int(math.Ceil(.25 * float64(rate)))
		if ticks != want {
			t.Fatalf("%d Hz: black pause took %d ticks, want %d", rate, ticks, want)
		}
	}
}

func TestRateSwitchPreservesElapsedPauseAndFade(t *testing.T) {
	l := &Screen{spec: Spec{Sector: 49, Blipps: 162}, rate: 50, tick: 235, fadeElapsed: .5}
	for i := 0; i < 5; i++ {
		l.advance()
	}
	volume := l.Volume()
	if err := l.SetTickRate(60); err != nil {
		t.Fatal(err)
	}
	if l.Volume() != volume {
		t.Fatal("rate switch changed elapsed fade")
	}
	for i := 0; i < 8; i++ {
		l.advance()
	}
	if l.Done() {
		t.Fatal("black pause finished early after rate switch")
	}
	l.advance()
	if !l.Done() {
		t.Fatal("black pause did not finish at 250 ms")
	}
}
