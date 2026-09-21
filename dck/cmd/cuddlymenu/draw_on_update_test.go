package main

import (
	"errors"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type fakeGame struct {
	updateErr error
	updates   int
	draws     int
	layouts   int
}

func (g *fakeGame) Update() error {
	g.updates++
	return g.updateErr
}

func (g *fakeGame) Draw(*ebiten.Image) {
	g.draws++
}

func (g *fakeGame) Layout(_, _ int) (int, int) {
	g.layouts++
	return 320, 240
}

func TestDrawOnUpdateGame(t *testing.T) {
	inner := &fakeGame{}
	game := newDrawOnUpdateGame(inner)

	game.Draw(nil)
	game.Draw(nil)
	if inner.draws != 1 {
		t.Fatalf("initial Draw count = %d, want 1", inner.draws)
	}

	if err := game.Update(); err != nil {
		t.Fatal(err)
	}
	game.Draw(nil)
	game.Draw(nil)
	if inner.draws != 2 {
		t.Fatalf("Draw count after Update = %d, want 2", inner.draws)
	}
}

func TestDrawOnUpdateGameRedrawsAfterLayout(t *testing.T) {
	inner := &fakeGame{}
	game := newDrawOnUpdateGame(inner)
	game.Draw(nil)

	width, height := game.Layout(640, 480)
	if width != 320 || height != 240 {
		t.Fatalf("Layout = (%d, %d), want (320, 240)", width, height)
	}
	game.Draw(nil)
	if inner.layouts != 1 || inner.draws != 2 {
		t.Fatalf("layouts=%d draws=%d, want layouts=1 draws=2", inner.layouts, inner.draws)
	}
}

func TestDrawOnUpdateGameDoesNotRedrawAfterFailedUpdate(t *testing.T) {
	wantErr := errors.New("update failed")
	inner := &fakeGame{updateErr: wantErr}
	game := newDrawOnUpdateGame(inner)
	game.Draw(nil)

	if err := game.Update(); !errors.Is(err, wantErr) {
		t.Fatalf("Update error = %v, want %v", err, wantErr)
	}
	game.Draw(nil)
	if inner.draws != 1 {
		t.Fatalf("Draw count after failed Update = %d, want 1", inner.draws)
	}
}

func TestNATIVECadenceIsIndependentOfDisplayRefresh(t *testing.T) {
	previous := ebiten.TPS()
	defer ebiten.SetTPS(previous)
	configureTiming()
	if ebiten.TPS() != 60 {
		t.Fatalf("playback rate = %d, want fixed 60 Hz", ebiten.TPS())
	}
	for _, refresh := range []int{60, 120, 144, 240} {
		inner := &fakeGame{}
		game := newDrawOnUpdateGame(inner)
		for tick := 0; tick < 60; tick++ {
			if err := game.Update(); err != nil {
				t.Fatal(err)
			}
			for draw := 0; draw < refresh; draw++ {
				game.Draw(nil)
			}
		}
		if inner.updates != 60 || inner.draws != 60 {
			t.Fatalf("refresh %d advanced animation: updates=%d draws=%d", refresh, inner.updates, inner.draws)
		}
	}
}

func TestRepeatedLayoutDoesNotForceRedraw(t *testing.T) {
	inner := &fakeGame{}
	game := newDrawOnUpdateGame(inner)
	game.Layout(640, 480)
	game.Draw(nil)
	for i := 0; i < 144; i++ {
		game.Layout(640, 480)
		game.Draw(nil)
	}
	if inner.draws != 1 {
		t.Fatalf("unchanged layout rebuilt %d frames", inner.draws)
	}
}
