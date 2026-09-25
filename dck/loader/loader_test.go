package loader

import (
	"github.com/olivierh59500/democonstructionkit/timeline"
	"math"
	"testing"
)

func screenAtTick(t *testing.T, spec Spec, rate, tick int) *Screen {
	t.Helper()
	countdown, err := timeline.NewCountdown(timeline.CountdownConfig{First: spec.Sector, Second: spec.Blipps, Hold: 24, FadeLead: 42})
	if err != nil {
		t.Fatal(err)
	}
	clock, err := newScreenClock(countdown, rate)
	if err != nil {
		t.Fatal(err)
	}
	l := &Screen{spec: spec, countdown: countdown, clock: clock}
	for i := 0; i < tick; i++ {
		l.advance()
	}
	return l
}

func TestCountdownBoundaries(t *testing.T) {
	spec := Spec{Sector: 49, Blipps: 162}
	for _, tt := range []struct {
		tick, sector, blipps int
		blank                bool
	}{{0, 49, 162, false}, {48, 1, 162, false}, {49, 0, 161, false}, {234, 0, -24, false}, {235, 0, -25, true}} {
		l := screenAtTick(t, spec, 60, tt.tick)
		s, b := l.Counters()
		if s != tt.sector || b != tt.blipps || l.Blank() != tt.blank {
			t.Fatalf("tick %d: counters %d/%d, blank %v", tt.tick, s, b, l.Blank())
		}
	}
}

func TestLoaderDurationsAtBothRates(t *testing.T) {
	for _, rate := range []int{50, 60} {
		spec := Spec{Sector: 49, Blipps: 162}
		l := screenAtTick(t, spec, rate, 0)
		for l.clock.Tick() < 169 {
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
		l = screenAtTick(t, spec, rate, 235)
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
	l := screenAtTick(t, Spec{Sector: 49, Blipps: 162}, 50, 235)
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
