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
	rate := flag.Int("hz", screens.TicksPerSecond, "fixed animation rate: 60 (default) or 50")
	tourMode := flag.Bool("tour", false, "play the complete introduction and screen tour once")
	tourOptions := app.DefaultTourOptions()
	flag.DurationVar(&tourOptions.ScreenDuration, "screen-duration", tourOptions.ScreenDuration, "time spent in each tour screen")
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
	ebiten.SetScreenClearedEveryFrame(false)
	config := app.Config{Screen: *id, Muted: *muted, TouchControls: *touch, TickRate: *rate}
	var g *app.Game
	var game ebiten.Game
	var err error
	if *tourMode {
		var tour *app.Tour
		tour, err = app.NewTour(config, tourOptions)
		if err == nil {
			g, game = tour.Game, tour
		}
	} else {
		g, err = app.New(config)
		game = g
	}
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	if err = g.SetTickRate(*rate); err != nil {
		log.Fatal(err)
	}
	if err = ebiten.RunGame(app.CacheDraws(game)); err != nil {
		log.Fatal(err)
	}
}

func ConfigureTiming() { ebiten.SetTPS(screens.TicksPerSecond) }
