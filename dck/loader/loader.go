// Package loader recreates the original image, sector countdown and decrunch phase.
package loader

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/assets"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/timeline"
	media "go-cuddlymenu/assets/cuddly"
	"go-cuddlymenu/dck/timing"
)

const Music = "menu/resources/loader.wav"
const Width, Height = 768, 536
const blackSeconds, fadeSeconds = .25, 1.5
const blackCue, fadeCue = 0, 1

type Spec struct {
	Sector, Blipps       int
	Title1, Title2, Text string
}
type Screen struct {
	Canvas           *ebiten.Image
	spec             Spec
	countdown        *timeline.Countdown
	clock            *timeline.CueClock
	gain             *timeline.CueRamp
	store            *assets.Store
	background, text *ebiten.Image
	font             scrolling.BitmapGrid
	scroll           *scrolling.Ring
}

func newScreenClock(countdown *timeline.Countdown, rate int) (*timeline.CueClock, error) {
	return timeline.NewCueClock(timeline.CueClockConfig{Rate: rate, Windows: []timeline.CueWindow{
		{StartTick: countdown.BlankTick(), Duration: blackSeconds, Tolerance: 1e-9},
		{StartTick: countdown.FadeStartTick(), Duration: fadeSeconds},
	}})
}

func New(name string) (*Screen, error) {
	return NewAtRate(name, timing.DefaultRate)
}

func NewAtRate(name string, rate int) (*Screen, error) {
	rate, err := timing.Normalize(rate)
	if err != nil {
		return nil, err
	}
	b, err := media.Files.ReadFile("loader.json")
	if err != nil {
		return nil, err
	}
	var specs map[string]Spec
	if err = json.Unmarshal(b, &specs); err != nil {
		return nil, err
	}
	spec, ok := specs[name]
	if !ok {
		return nil, fmt.Errorf("unknown loading screen %q", name)
	}
	countdown, err := timeline.NewCountdown(timeline.CountdownConfig{First: spec.Sector, Second: spec.Blipps, Hold: 24, FadeLead: 42})
	if err != nil {
		return nil, err
	}
	clock, err := newScreenClock(countdown, rate)
	if err != nil {
		return nil, err
	}
	gain, err := timeline.NewCueRamp(clock, timeline.CueRampConfig{Window: fadeCue, From: 1, To: 0})
	if err != nil {
		return nil, err
	}
	l := &Screen{spec: spec, countdown: countdown, clock: clock, gain: gain, store: assets.New(media.Files), Canvas: ebiten.NewImage(Width, Height), text: ebiten.NewImage(640, 16)}
	load := func(name string) *ebiten.Image {
		img, e := l.store.Texture("menu/resources/" + name)
		if e != nil {
			err = e
		}
		return img
	}
	l.background = load("loader.png")
	brown, grey := load("fontN.png"), load("fontL.png")
	if err != nil {
		l.Close()
		return nil, err
	}
	l.font, err = presets.BitmapFont("cuddly-loader", brown, ebiten.FilterNearest)
	if err != nil {
		l.Close()
		return nil, err
	}
	f := l.font
	f.Image = grey
	f.Columns = grey.Bounds().Dx() / 16
	l.scroll, err = scrolling.NewRing(scrolling.RingConfig{Text: spec.Text, Font: f, Viewport: 640, Speed: 6})
	if err != nil {
		l.Close()
		return nil, err
	}
	l.render()
	return l, nil
}

// Counters describe the displayed values before the source's end-of-frame decrement.
func (l *Screen) Counters() (sector, blipps int) {
	state := l.countdown.At(l.clock.Tick())
	return state.First, state.Second
}
func (l *Screen) Blank() bool     { return l.countdown.At(l.clock.Tick()).Blank }
func (l *Screen) Done() bool      { return l.Blank() && l.clock.Done(blackCue) }
func (l *Screen) Volume() float64 { return l.gain.Value() }
func (l *Screen) SetTickRate(rate int) error {
	rate, err := timing.Normalize(rate)
	if err != nil {
		return err
	}
	return l.clock.SetRate(rate)
}
func (l *Screen) advance() { l.clock.Step() }
func (l *Screen) Update() error {
	if !l.Done() {
		l.advance()
		l.render()
	}
	return nil
}
func (l *Screen) Draw(dst *ebiten.Image)     { dst.DrawImage(l.Canvas, nil) }
func (l *Screen) Layout(int, int) (int, int) { return Width, Height }
func (l *Screen) render() {
	l.Canvas.Fill(color.Black)
	if l.Blank() {
		return
	}
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(64, 64)
	l.Canvas.DrawImage(l.background, &op)
	l.font.Print(l.Canvas, l.spec.Title1, 496, 124, 1, 1)
	l.font.Print(l.Canvas, l.spec.Title2, 448, 140, 1, 1)
	sector, blipps := l.Counters()
	l.font.Print(l.Canvas, fmt.Sprintf("%3d", sector), 562, 216, 1, 1)
	if sector == 0 {
		l.font.Print(l.Canvas, "DECRUNCH MICRO", 482, 244, 1, 1)
		l.font.Print(l.Canvas, " BLIPPS TO GO", 482, 268, 1, 1)
		l.font.Print(l.Canvas, fmt.Sprintf("%3d", max(0, blipps)), 562, 284, 1, 1)
	}
	l.text.Clear()
	l.scroll.Step()
	l.scroll.DrawAt(l.text, 0, 0)
	op.GeoM.Reset()
	op.GeoM.Translate(64, 450)
	l.Canvas.DrawImage(l.text, &op)
}
func (l *Screen) Close() error { l.Canvas.Deallocate(); l.text.Deallocate(); return l.store.Close() }
