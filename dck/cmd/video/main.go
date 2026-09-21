// Command video records the introduction, every screen and the complete route.
package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/video"
	"go-cuddlymenu/dck/app"
)

func main() {
	c := video.Config{Output: "cuddly-demo.mp4", Title: "The Cuddly Demos", Width: 832, Height: 552, FPS: 60, TPS: 60, SampleRate: 48000}
	c.Flags(flag.CommandLine)
	tour := app.DefaultTourOptions()
	flag.DurationVar(&tour.IntroDuration, "intro-duration", tour.IntroDuration, "time spent in the introduction")
	flag.DurationVar(&tour.ScreenDuration, "screen-duration", tour.ScreenDuration, "time spent in each screen, excluding transitions")
	flag.DurationVar(&tour.MenuDuration, "menu-duration", tour.MenuDuration, "pause in the menu before walking to the next door")
	flag.Parse()
	if err := video.Run(c, func() (ebiten.Game, error) { return app.NewTour(app.Config{}, tour) }); err != nil {
		log.Fatal(err)
	}
}
