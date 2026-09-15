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
