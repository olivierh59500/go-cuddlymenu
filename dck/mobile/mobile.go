// Package mobile binds the complete controller with a menu-first entry point.
package mobile

import (
	"github.com/hajimehoshi/ebiten/v2"
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"
	"go-cuddlymenu/dck/mobilehost"
	"go-cuddlymenu/dck/screens"
)

func init() {
	ebiten.SetTPS(screens.TicksPerSecond)
	ebiten.SetScreenClearedEveryFrame(false)
	enginemobile.SetGame(mobilehost.New("menu"))
}

// Dummy forces gomobile to include this package in the Android binding.
func Dummy() {}
