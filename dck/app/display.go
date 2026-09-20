package app

import "github.com/hajimehoshi/ebiten/v2"

// CacheDraws retains the latest screen between updates. SetScreenClearedEveryFrame
// must be false. Repeated Layout calls only invalidate a changed logical size.
func CacheDraws(game ebiten.Game) ebiten.Game { return &display{game: game, dirty: true} }

type display struct {
	game          ebiten.Game
	dirty         bool
	width, height int
}

func (d *display) Update() error {
	if err := d.game.Update(); err != nil {
		return err
	}
	d.dirty = true
	return nil
}
func (d *display) Draw(dst *ebiten.Image) {
	if d.dirty {
		d.game.Draw(dst)
		d.dirty = false
	}
}
func (d *display) Layout(w, h int) (int, int) {
	width, height := d.game.Layout(w, h)
	if width != d.width || height != d.height {
		d.dirty = true
		d.width, d.height = width, height
	}
	return width, height
}
