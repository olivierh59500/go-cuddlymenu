// Package mobile exposes the game to ebitenmobile.
package mobile

import (
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"

	"go-cuddlymenu/dck/menu"
)

func init() {
	enginemobile.SetGame(menu.NewGame())
}

// Dummy forces gomobile to include this package in the Android binding.
func Dummy() {}
