package app

import (
	"io"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/sound/output"
)

func TestKnucklebusterFollowsAudiblePlaybackAcrossPauseSeekAndRateChanges(t *testing.T) {
	previousRate := ebiten.TPS()
	defer ebiten.SetTPS(previousRate)
	session, err := output.Begin(48000)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	g, err := New(Config{Screen: "knucklebuster"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if err := g.startAudio(); err != nil {
		t.Fatal(err)
	}
	if err := session.Advance(8, 1, io.Discard); err != nil {
		t.Fatal(err)
	}
	if got := g.scene.MusicPosition(); got != 8*time.Second {
		t.Fatalf("scene used a frame counter instead of audible playback: %v", got)
	}
	g.player.Pause()
	for _, rate := range []int{50, 60} {
		if err := g.SetTickRate(rate); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < rate; i++ {
			if err := g.scene.Update(); err != nil {
				t.Fatal(err)
			}
		}
		if got := g.scene.MusicPosition(); got != 8*time.Second {
			t.Fatalf("paused soundtrack advanced with %d Hz graphics: %v", rate, got)
		}
	}
	if err := g.player.Seek(7 * time.Second); err != nil {
		t.Fatal(err)
	}
	if got := g.scene.MusicPosition(); got != 7*time.Second {
		t.Fatalf("backwards audio seek left stale animation time: %v", got)
	}
}
