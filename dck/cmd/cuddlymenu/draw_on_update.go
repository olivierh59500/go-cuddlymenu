package main

import "github.com/hajimehoshi/ebiten/v2"

// drawOnUpdateGame preserves the last frame when the monitor refreshes faster
// than the game's update rate. All visual state changes happen in Update, so
// rebuilding the same frame multiple times only wastes CPU and allocations.
type drawOnUpdateGame struct {
	game  ebiten.Game
	dirty bool
}

func newDrawOnUpdateGame(game ebiten.Game) *drawOnUpdateGame {
	return &drawOnUpdateGame{game: game, dirty: true}
}

func (g *drawOnUpdateGame) Update() error {
	if err := g.game.Update(); err != nil {
		return err
	}
	g.dirty = true
	return nil
}

func (g *drawOnUpdateGame) Draw(screen *ebiten.Image) {
	if !g.dirty {
		return
	}
	g.game.Draw(screen)
	g.dirty = false
}

func (g *drawOnUpdateGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.dirty = true
	return g.game.Layout(outsideWidth, outsideHeight)
}
