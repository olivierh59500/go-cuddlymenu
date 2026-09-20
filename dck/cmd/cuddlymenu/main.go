package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"go-cuddlymenu/dck/app"
	"go-cuddlymenu/dck/screens"
)

func main() {
	id := flag.String("screen", "menu", "screen ID, or menu")
	muted := flag.Bool("mute", false, "disable device audio")
	list := flag.Bool("list", false, "list screens and port status")
	flag.Parse()
	if *list {
		for _, d := range screens.Catalog() {
			state := "planned"
			if d.Ready {
				state = "available"
			}
			fmt.Printf("%-14s %-10s %s\n", d.ID, state, d.Title)
		}
		return
	}
	ebiten.SetWindowSize(768, 536)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Cuddly Demos / DCK - F1: screens, Esc: menu")
	ebiten.SetScreenClearedEveryFrame(false)
	game, err := app.New(app.Config{Screen: *id, Muted: *muted})
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()
	if err := ebiten.RunGame(newDrawOnUpdateGame(game)); err != nil {
		log.Fatal(err)
	}
}
