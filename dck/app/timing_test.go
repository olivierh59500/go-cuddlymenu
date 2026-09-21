package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"testing"
)

func TestChangingRatePreservesSceneAndMusic(t *testing.T) {
	previous := ebiten.TPS()
	defer ebiten.SetTPS(previous)
	g, err := New(Config{Screen: "colorshock", Muted: true})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if g.TickRate() != 60 {
		t.Fatal("default animation rate is not 60 Hz")
	}
	scene, track := g.scene, g.track
	for _, rate := range []int{50, 60} {
		if err = g.SetTickRate(rate); err != nil {
			t.Fatal(err)
		}
		if g.TickRate() != rate || ebiten.TPS() != rate {
			t.Fatal("application and engine rates disagree")
		}
		if g.scene != scene || g.track != track {
			t.Fatal("comparison reset the scene or music")
		}
	}
	if err = g.SetTickRate(144); err == nil || ebiten.TPS() != 60 {
		t.Fatal("invalid rate changed the clock")
	}
}

func TestComparisonRateReachesLoader(t *testing.T) {
	g, err := New(Config{Muted: true, TickRate: 50})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if err = g.Begin("menu"); err != nil {
		t.Fatal(err)
	}
	for !g.transition.Blank() {
		g.transition.Update()
	}
	ticks := 0
	for !g.transition.Done() {
		g.transition.Update()
		ticks++
	}
	if ticks != 13 {
		t.Fatalf("50 Hz loader pause took %d ticks, want 13", ticks)
	}
}
