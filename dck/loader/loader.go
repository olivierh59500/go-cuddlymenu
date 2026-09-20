// Package loader recreates the original image, sector countdown and decrunch phase.
package loader

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/assets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	media "go-cuddlymenu/assets/cuddly"
)

const Music = "menu/resources/loader.wav"
const Width, Height = 768, 536
const BlackTicks = 13 // The original 250 ms pause, rounded up at PAL 50 Hz.

type Spec struct {
	Sector, Blipps       int
	Title1, Title2, Text string
}
type Screen struct {
	Canvas           *ebiten.Image
	spec             Spec
	tick             int
	store            *assets.Store
	background, text *ebiten.Image
	font             scrolling.BitmapGrid
	scroll           *scrolling.Ring
}

func New(name string) (*Screen, error) {
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
	l := &Screen{spec: spec, store: assets.New(media.Files), Canvas: ebiten.NewImage(Width, Height), text: ebiten.NewImage(640, 16)}
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
	l.font = scrolling.BitmapGrid{Image: brown, Width: 16, Height: 16, Columns: brown.Bounds().Dx() / 16, First: 32}
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
	return max(0, l.spec.Sector-l.tick), l.spec.Blipps - max(0, l.tick-l.spec.Sector+1)
}
func (l *Screen) Blank() bool { return l.tick >= l.spec.Sector+l.spec.Blipps+24 }
func (l *Screen) Done() bool  { return l.tick >= l.spec.Sector+l.spec.Blipps+24+BlackTicks }
func (l *Screen) Volume() float64 {
	return max(0, 1-float64(max(0, l.tick-(l.spec.Sector+l.spec.Blipps-42)))/75)
}
func (l *Screen) Update() error {
	if !l.Done() {
		l.tick++
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
