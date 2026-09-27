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
	rate    int
	tour    *app.TourOptions
}
type Host struct {
	first                string
	requests             chan request
	game                 *app.Game
	presentation         ebiten.Game
	metrics              bool
	touring              bool
	lastScreen           string
	warmup, ticks        int
	updateTime, drawTime time.Duration
	draws                int
}

func New(first string) *Host { return &Host{first: first, requests: make(chan request, 1)} }

// Configure is safe on Android's UI thread. Construction happens in Update.
// Warmup is a debug-only number of normally presented frames before profiling.
func (h *Host) Configure(screen string, metrics bool, warmup int) {
	h.ConfigureAtRate(screen, metrics, warmup, 0)
}

// ConfigureAtRate optionally selects 50 Hz for comparison; zero defaults to 60.
func (h *Host) ConfigureAtRate(screen string, metrics bool, warmup, rate int) {
	if screen == "" {
		screen = h.first
	}
	r := request{screen: screen, metrics: metrics, warmup: max(0, min(10000, warmup)), rate: rate}
	select {
	case <-h.requests:
	default:
	}
	h.requests <- r
}

// ConfigureTour queues an unattended route on the UI thread. It retains any
// preceding rate/metrics request and builds the tour in Update, on the graphics
// thread. Durations are seconds and remain independent of the 50/60 Hz rate.
func (h *Host) ConfigureTour(introSeconds, screenSeconds, menuSeconds int) bool {
	if introSeconds < 1 || introSeconds > 600 || screenSeconds < 1 || screenSeconds > 600 ||
		menuSeconds < 1 || menuSeconds > 600 {
		return false
	}
	r := request{screen: h.first}
	select {
	case r = <-h.requests:
	default:
	}
	r.tour = &app.TourOptions{
		IntroDuration:  time.Duration(introSeconds) * time.Second,
		ScreenDuration: time.Duration(screenSeconds) * time.Second,
		MenuDuration:   time.Duration(menuSeconds) * time.Second,
	}
	h.requests <- r
	return true
}

func (h *Host) Update() error {
	select {
	case r := <-h.requests:
		if h.game != nil {
			h.game.Close()
		}
		var g *app.Game
		var presentation ebiten.Game
		var err error
		if r.tour != nil {
			var tour *app.Tour
			tour, err = app.NewTour(app.Config{TickRate: r.rate}, *r.tour)
			if err == nil {
				g, presentation = tour.Game, app.CacheDraws(tour)
			}
		} else {
			g, err = app.New(app.Config{Screen: r.screen, TickRate: r.rate})
			if err == nil {
				presentation = app.CacheDraws(g)
			}
		}
		if err != nil {
			return err
		}
		if err = g.SetTickRate(g.TickRate()); err != nil {
			g.Close()
			return err
		}
		h.game = g
		h.presentation = presentation
		h.metrics = r.metrics
		h.touring = r.tour != nil
		h.lastScreen = ""
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
	if h.touring {
		if screen := h.game.CurrentScreen(); screen != h.lastScreen {
			log.Printf("cuddly_tour screen=%s", screen)
			h.lastScreen = screen
		}
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
		log.Print(fmt.Sprintf("cuddly_perf screen=%s hz=%d tps=%.2f fps=%.2f update_ms=%.3f draw_ms=%.3f frames=%d", h.game.CurrentScreen(), h.game.TickRate(), ebiten.ActualTPS(), ebiten.ActualFPS(), float64(h.updateTime.Microseconds())/250000, float64(h.drawTime.Microseconds())/float64(max(1, h.draws)*1000), h.ticks))
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
