package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"go-cuddlymenu/menu"
)

func main() {
	ebiten.SetWindowSize(768, 536)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Cuddly Demos - Menu")
	ebiten.SetScreenClearedEveryFrame(false)
	if err := ebiten.RunGame(newDrawOnUpdateGame(menu.NewGame())); err != nil {
		log.Fatal(err)
	}
}
