package mobilehost

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go-cuddlymenu/dck/app"
	"testing"
)

type counterGame struct{ updates int }

func (g *counterGame) Update() error              { g.updates++; return nil }
func (g *counterGame) Draw(*ebiten.Image)         {}
func (g *counterGame) Layout(int, int) (int, int) { return 768, 540 }

func TestWarmupDoesNotAccumulateOffscreenFrames(t *testing.T) {
	game := &counterGame{}
	h := New("intro")
	h.game = &app.Game{}
	h.presentation = game
	h.warmup = 30
	for i := 0; i < 30; i++ {
		if err := h.Update(); err != nil {
			t.Fatal(err)
		}
		if game.updates != i+1 {
			t.Fatal("warmup bypassed normal presentation")
		}
	}
	if h.ticks != 0 || h.warmup != 0 {
		t.Fatal("warmup contaminated measured frames")
	}
	h.Update()
	if h.ticks != 1 || game.updates != 31 {
		t.Fatal("normal playback did not resume")
	}
}

func TestConfigurationDefersConstructionAndKeepsLatestRequest(t *testing.T) {
	h := New("intro")
	h.Configure("dna", false, 100)
	h.Configure("reset", true, 3000)
	r := <-h.requests
	if r.screen != "reset" || !r.metrics || r.warmup != 3000 || h.game != nil {
		t.Fatal("configuration created the game on the caller thread or lost the latest request")
	}
}
