package mobilehost

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"go-cuddlymenu/dck/app"
)

type counterGame struct{ updates int }

func (g *counterGame) Update() error              { g.updates++; return nil }
func (g *counterGame) Draw(*ebiten.Image)         {}
func (g *counterGame) Layout(int, int) (int, int) { return 768, 540 }

func TestWarmupDoesNotAccumulateOffscreenFrames(t *testing.T) {
	game := &counterGame{}
	h := New("intro")
	h.game = &app.Game{}
	h.presentation = game
	h.warmup = 30
	for i := 0; i < 30; i++ {
		if err := h.Update(); err != nil {
			t.Fatal(err)
		}
		if game.updates != i+1 {
			t.Fatal("warmup bypassed normal presentation")
		}
	}
	if h.ticks != 0 || h.warmup != 0 {
		t.Fatal("warmup contaminated measured frames")
	}
	h.Update()
	if h.ticks != 1 || game.updates != 31 {
		t.Fatal("normal playback did not resume")
	}
}

func TestConfigurationDefersConstructionAndKeepsLatestRequest(t *testing.T) {
	h := New("intro")
	h.Configure("dna", false, 100)
	h.Configure("reset", true, 3000)
	r := <-h.requests
	if r.screen != "reset" || !r.metrics || r.warmup != 3000 || h.game != nil {
		t.Fatal("configuration created the game on the caller thread or lost the latest request")
	}
}

func TestMobileComparisonRateReachesLaunchRequest(t *testing.T) {
	h := New("intro")
	h.ConfigureAtRate("", true, 120, 50)
	r := <-h.requests
	if r.screen != "intro" || r.rate != 50 || r.warmup != 120 {
		t.Fatal("mobile rate override was lost")
	}
	h.Configure("intro", false, 0)
	if r = <-h.requests; r.rate != 0 {
		t.Fatal("normal launch did not select the default rate")
	}
}

func TestTourRequestKeepsAndroidRateAndMetrics(t *testing.T) {
	host := New("intro")
	host.ConfigureAtRate("dna", true, 12, 50)
	host.ConfigureMenuCRT(true)
	if host.ConfigureTour(0, 6, 1) {
		t.Fatal("accepted a zero-length introduction")
	}
	if !host.ConfigureTour(10, 6, 1) {
		t.Fatal("rejected a valid complete tour")
	}
	select {
	case request := <-host.requests:
		if request.rate != 50 || !request.metrics || request.warmup != 12 || !request.crt || request.tour == nil ||
			request.tour.IntroDuration != 10*time.Second ||
			request.tour.ScreenDuration != 6*time.Second ||
			request.tour.MenuDuration != time.Second {
			t.Fatalf("tour lost the Android launch configuration: %+v", request)
		}
	default:
		t.Fatal("tour request was not queued")
	}
}
