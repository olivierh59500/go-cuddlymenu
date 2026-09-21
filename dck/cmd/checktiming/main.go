// Command checktiming verifies fixed updates independently of unrestricted draws.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"go-cuddlymenu/dck/app"
	"go-cuddlymenu/dck/screens"
)

type check struct {
	game           *app.Game
	start          time.Time
	updates, draws int
	elapsed        time.Duration
}

func (c *check) Update() error {
	if c.start.IsZero() {
		c.start = time.Now()
	}
	c.elapsed = time.Since(c.start)
	if c.elapsed >= 3*time.Second {
		return ebiten.Termination
	}
	c.updates++
	return c.game.Update()
}
func (c *check) Draw(dst *ebiten.Image)     { c.draws++; c.game.Draw(dst) }
func (c *check) Layout(w, h int) (int, int) { return c.game.Layout(w, h) }
func main() {
	rate := flag.Int("hz", screens.TicksPerSecond, "fixed animation rate: 50 or 60")
	out := flag.String("out", "docs/cuddly-timing.json", "timing report path")
	flag.Parse()
	ebiten.SetVsyncEnabled(false)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(768, 540)
	g, err := app.New(app.Config{Screen: "starwars", Muted: true, TickRate: *rate})
	if err != nil {
		panic(err)
	}
	defer g.Close()
	if err = g.SetTickRate(*rate); err != nil {
		panic(err)
	}
	c := &check{game: g}
	if err := ebiten.RunGame(c); err != nil {
		panic(err)
	}
	actual := float64(c.updates) / c.elapsed.Seconds()
	report := map[string]any{"Screen": "starwars", "ConfiguredTPS": ebiten.TPS(), "ElapsedSeconds": c.elapsed.Seconds(), "Updates": c.updates, "DrawCalls": c.draws, "MeasuredUpdatesPerSecond": actual, "Vsync": false, "Audio": "muted; native audio sample rate is independent of visual TPS"}
	b, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(b))
	if actual < float64(*rate)-2 || actual > float64(*rate)+2 {
		panic("animation cadence outside tolerance")
	}
	if err := os.WriteFile(*out, append(b, '\n'), 0644); err != nil {
		panic(err)
	}
}
