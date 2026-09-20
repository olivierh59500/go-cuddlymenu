package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go-cuddlymenu/dck/app"
)

// Keep the existing entry-point tests attached to the shared presenter.
func newDrawOnUpdateGame(game ebiten.Game) ebiten.Game { return app.CacheDraws(game) }
