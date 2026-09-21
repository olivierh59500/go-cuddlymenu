// Package mobilehost delays graphics/audio construction until Android's view is ready.
package mobilehost

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"go-cuddlymenu/dck/app"
)

type request struct {
	screen  string
	metrics bool
	warmup  int
}
type Host struct {
	first                string
	requests             chan request
	game                 *app.Game
	presentation         ebiten.Game
	metrics              bool
	warmup, ticks        int
	updateTime, drawTime time.Duration
	draws                int
}

func New(first string) *Host { return &Host{first: first, requests: make(chan request, 1)} }

// Configure is safe on Android's UI thread. Construction happens in Update.
// Warmup is a debug-only number of normally presented frames before profiling.
func (h *Host) Configure(screen string, metrics bool, warmup int) {
	if screen == "" {
		screen = h.first
	}
	r := request{screen, metrics, max(0, min(10000, warmup))}
	select {
	case <-h.requests:
	default:
	}
	h.requests <- r
}
func (h *Host) Update() error {
	select {
	case r := <-h.requests:
		if h.game != nil {
			h.game.Close()
		}
		g, err := app.New(app.Config{Screen: r.screen})
		if err != nil {
			return err
		}
		h.game = g
		h.presentation = app.CacheDraws(g)
		h.metrics = r.metrics
		h.warmup = r.warmup
		h.ticks = 0
		h.updateTime = 0
		h.drawTime = 0
		h.draws = 0
	default:
	}
	if h.game == nil {
		h.Configure(h.first, false, 0)
		return nil
	}
	start := time.Now()
	if err := h.presentation.Update(); err != nil {
		return err
	}
	// Keep one presented animation step per update. Accelerated offscreen
	// warmup can retain large antialiased/feedback GPU command chains on Android.
	if h.warmup > 0 {
		h.warmup--
		h.drawTime = 0
		h.draws = 0
		return nil
	}
	h.updateTime += time.Since(start)
	h.ticks++
	if h.metrics && h.ticks%250 == 0 {
		log.Print(fmt.Sprintf("cuddly_perf screen=%s tps=%.2f fps=%.2f update_ms=%.3f draw_ms=%.3f frames=%d", h.game.CurrentScreen(), ebiten.ActualTPS(), ebiten.ActualFPS(), float64(h.updateTime.Microseconds())/250000, float64(h.drawTime.Microseconds())/float64(max(1, h.draws)*1000), h.ticks))
		h.updateTime = 0
		h.drawTime = 0
		h.draws = 0
	}
	return nil
}
func (h *Host) Draw(dst *ebiten.Image) {
	if h.presentation != nil {
		start := time.Now()
		h.presentation.Draw(dst)
		h.drawTime += time.Since(start)
		h.draws++
	}
}
func (h *Host) Layout(w, hgt int) (int, int) {
	if h.presentation != nil {
		return h.presentation.Layout(w, hgt)
	}
	return 768, 540
}
