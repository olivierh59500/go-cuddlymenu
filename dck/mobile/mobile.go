// Package mobile binds the complete controller with a menu-first entry point.
package mobile

import (
	"github.com/hajimehoshi/ebiten/v2"
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"
	"go-cuddlymenu/dck/mobilehost"
	"go-cuddlymenu/dck/screens"
)

var host = mobilehost.New("menu")

func init() {
	ebiten.SetTPS(screens.TicksPerSecond)
	ebiten.SetScreenClearedEveryFrame(false)
	enginemobile.SetGame(host)
}

// ConfigureCRT enables the optional menu material for device profiling.
func ConfigureCRT(enabled bool) { host.ConfigureMenuCRT(enabled) }

// Dummy forces gomobile to include this package in the Android binding.
func Dummy() {}
