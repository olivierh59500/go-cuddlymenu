// Package desktop shares the menu and complete-production desktop launchers.
package desktop

import (
	"flag"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go-cuddlymenu/dck/app"
	"go-cuddlymenu/dck/screens"
)

func Run(first string) {
	id := flag.String("screen", first, "initial screen ID, or menu")
	muted := flag.Bool("mute", false, "disable device audio")
	touch := flag.Bool("touch", false, "show touch controls on desktop")
	list := flag.Bool("list", false, "list native screens")
	flag.Parse()
	if *list {
		for _, d := range screens.Catalog() {
			fmt.Printf("%-14s %s\n", d.ID, d.Title)
		}
		return
	}
	ebiten.SetWindowSize(768, 540)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Cuddly Demo / DCK - Space: menu, F1: screens")
	ConfigureTiming()
	ebiten.SetScreenClearedEveryFrame(false)
	g, err := app.New(app.Config{Screen: *id, Muted: *muted, TouchControls: *touch})
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	if err = ebiten.RunGame(app.CacheDraws(g)); err != nil {
		log.Fatal(err)
	}
}

func ConfigureTiming() { ebiten.SetTPS(screens.TicksPerSecond) }
