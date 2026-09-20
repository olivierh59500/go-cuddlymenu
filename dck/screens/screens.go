// Package screens reconstructs the Cuddly NATIVE screens with native DCK effects.
package screens

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/assets"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	media "go-cuddlymenu/assets/cuddly"
)

type Descriptor struct {
	ID, Title, Directory, Page, Music string
	Width, Height                     int
	Ready                             bool
}

// TicksPerSecond is the PAL playback cadence, independent of monitor refresh.
// Source frame counters retain their original ordering at this fixed rate.
const TicksPerSecond = 50

var catalog = []Descriptor{
	{"big-sprite", "The Big Sprite Demo", "big_sprite", "big-sprite", "cuddly_bigsprite.ym", 768, 540, true},
	{"colorshock", "Colorshock II", "colorshock2", "colorshock", "Colorshock.ym", 768, 540, true},
	{"ehh", "Ehhh!!!! / No Name 1", "ehh", "ehh", "@ym/ehh.ym", 768, 540, true},
	{"megascroller", "The Mega Scroller", "megascroller", "megascroller", "@ym/megascroller.ym", 768, 540, true},
	{"spreadpoint", "Spreadpoint", "spreadpoint", "spreadpoint", "master.mp3", 832, 552, true},
	{"digi", "Digi Sound", "digi", "digi", "tcbdigi.mp3", 768, 540, true},
	{"led", "The LED Scroller", "led_scroller", "led", "@ym/led.ym", 768, 540, true},
	{"3d-doc", "The 3D DOC", "3d_doc", "3d-doc", "Cuddly - 3D doc.ym", 768, 540, true},
	{"fullscreen", "The Fullscreen Demo", "fullscreen", "fullscreen", "@ym/fullscreen.ym", 768, 536, true},
	{"starwars", "The Starwars Demo", "starwars", "starwars", "Cuddly - Star-Wars.ym", 768, 540, true},
	{"knucklebuster", "Knucklebuster", "tex", "knucklebuster", "@ym/knucklebusters.ym", 768, 540, true},
	{"dna", "The DNA Demo", "dna_demo", "dna", "bankok-knights-1.ym", 832, 552, true},
	{"megaball", "The Megaball Demo / No Name 2", "megaball", "megaball", "Cuddly - Megaballs.ym", 768, 540, true},
	{"intro", "Introduction", "intro", "intro", "cuddlyintro.mp3", 768, 540, true},
	{"reset", "Reset Screen", "reset", "reset", "cuddlyreset.ym", 768, 540, true},
}

func Catalog() []Descriptor { return append([]Descriptor(nil), catalog...) }
func Find(id string) (Descriptor, bool) {
	for _, d := range catalog {
		if d.ID == id {
			return d, true
		}
	}
	return Descriptor{}, false
}
func DoorID(name string) string {
	return map[string]string{"BIG_SPRITE": "big-sprite", "COLORSHOCK_II": "colorshock", "NO_NAME_1": "ehh", "MEGA_SCROLLER": "megascroller", "SPREADPOINT": "spreadpoint", "DIGI_DEMO": "digi", "LED_SCROLLER": "led", "DOC": "3d-doc", "FULLSCREEN": "fullscreen", "STARWARS_DEMO": "starwars", "KNUCKLE_BUSTER": "knucklebuster", "DNA_DEMO": "dna", "NO_NAME_2": "megaball"}[name]
}

func DoorName(id string) string {
	if id == "menu" {
		return "MENU"
	}
	for _, name := range []string{"BIG_SPRITE", "COLORSHOCK_II", "NO_NAME_1", "MEGA_SCROLLER", "SPREADPOINT", "DIGI_DEMO", "LED_SCROLLER", "DOC", "FULLSCREEN", "STARWARS_DEMO", "KNUCKLE_BUSTER", "DNA_DEMO", "NO_NAME_2"} {
		if DoorID(name) == id {
			return name
		}
	}
	return ""
}

// AudioCue selects scene-relative music, or silence when File is empty.
type AudioCue struct {
	File string
	Loop bool
}
type Input struct{ Left, Right, Up, Down bool }

type sourceData struct {
	Strings map[string]string
	Numbers map[string][]float64
	Lists   map[string][]string
}
type Scene struct {
	Descriptor   Descriptor
	Canvas       *ebiten.Image
	store        *assets.Store
	surfaces     []*ebiten.Image
	filters      map[*ebiten.Image]ebiten.Filter
	data         sourceData
	render       func()
	frame        uint64
	random       uint32
	err          error
	input        func(Input)
	audioCue     *AudioCue
	closeEffects []func() error
}

func New(id string) (*Scene, error) {
	d, ok := Find(id)
	if !ok {
		return nil, fmt.Errorf("unknown screen %q", id)
	}
	if !d.Ready {
		return nil, fmt.Errorf("%s is not available yet", d.Title)
	}
	s := &Scene{Descriptor: d, store: assets.New(media.Files), filters: map[*ebiten.Image]ebiten.Filter{}, random: 42}
	// The reference consumes one random value for the initial random sequence offset
	// before initializing either stars or the drummer's trigger sequence.
	s.rnd()
	bytes, err := media.Files.ReadFile("data.json")
	if err != nil {
		return nil, err
	}
	var all map[string]sourceData
	if err = json.Unmarshal(bytes, &all); err != nil {
		return nil, err
	}
	s.data = all[d.Page]
	s.Canvas = s.surface(d.Width, d.Height)
	s.music(d.Music, true)
	switch id {
	case "big-sprite":
		s.bigSprite()
	case "colorshock":
		s.colorshock()
	case "megascroller":
		s.megaScroller()
	case "fullscreen":
		s.fullscreen()
	case "digi":
		s.digi()
	case "led":
		s.led()
	case "knucklebuster":
		s.knucklebuster()
	case "ehh":
		s.ehh()
	case "intro":
		s.intro()
	case "reset":
		s.reset()
	case "megaball":
		s.megaball()
	case "spreadpoint":
		s.spreadpoint()
	case "dna":
		s.dna()
	case "3d-doc":
		s.doc()
	case "starwars":
		s.starwars()
	}
	if s.err != nil {
		s.Close()
		return nil, s.err
	}
	if s.render == nil {
		s.Close()
		return nil, fmt.Errorf("missing screen renderer %q", id)
	}
	s.render() // NATIVE init() draws the first frame before scheduling the next one.
	return s, nil
}

func (s *Scene) Update() error {
	if s.err != nil {
		return s.err
	}
	s.frame++
	s.render()
	return s.err
}
func (s *Scene) Draw(dst *ebiten.Image)     { dst.DrawImage(s.Canvas, nil) }
func (s *Scene) Layout(int, int) (int, int) { return s.Descriptor.Width, s.Descriptor.Height }
func (s *Scene) Close() error {
	for _, close := range s.closeEffects {
		close()
	}
	s.closeEffects = nil
	for _, img := range s.surfaces {
		img.Deallocate()
	}
	s.surfaces = nil
	return s.store.Close()
}

func (s *Scene) Input(in Input) {
	if s.input != nil {
		s.input(in)
	}
}
func (s *Scene) music(file string, loop bool) { s.audioCue = &AudioCue{File: file, Loop: loop} }
func (s *Scene) TakeAudioCue() *AudioCue      { cue := s.audioCue; s.audioCue = nil; return cue }
func (s *Scene) asset(name string) *ebiten.Image {
	img, err := s.store.Texture("" + s.Descriptor.Directory + "/" + name)
	if err != nil {
		s.err = err
		return s.surface(1, 1)
	}
	return img
}
func (s *Scene) surface(w, h int) *ebiten.Image {
	img := ebiten.NewImageWithOptions(image.Rect(0, 0, w, h), &ebiten.NewImageOptions{Unmanaged: true})
	s.surfaces = append(s.surfaces, img)
	s.filters[img] = ebiten.FilterLinear
	return img
}
func (s *Scene) ring(dst *ebiten.Image, font *ebiten.Image, w, h float64, first rune, text string, speed float64) *scrolling.Ring {
	r, err := scrolling.NewRing(scrolling.RingConfig{Text: text, Font: scrolling.BitmapGrid{Image: font, Width: w, Height: h, Columns: int(float64(font.Bounds().Dx()) / w), ColumnSpan: float64(font.Bounds().Dx()) / w, First: first, Filter: s.filters[dst]}, Viewport: float64(dst.Bounds().Dx()), Speed: speed, Controls: true})
	if err != nil {
		s.err = err
	}
	return r
}
func (s *Scene) draw(dst, src *ebiten.Image, x, y float64) {
	s.transform(dst, src, x, y, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
}
func (s *Scene) transform(dst, src *ebiten.Image, x, y, sx, sy, angle, hx, hy, alpha float64, blend ebiten.Blend) {
	op := ebiten.DrawImageOptions{Filter: s.filters[dst], Blend: blend}
	op.GeoM.Translate(-hx, -hy)
	op.GeoM.Scale(sx, sy)
	op.GeoM.Rotate(angle * math.Pi / 180)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleAlpha(float32(alpha))
	composite.Instance{Image: src, Options: op}.Draw(dst)
}
func (s *Scene) part(dst, src *ebiten.Image, r composite.Region, x, y, sx, sy float64) {
	op := ebiten.DrawImageOptions{Filter: s.filters[dst]}
	op.GeoM.Scale(sx, sy)
	op.GeoM.Translate(x, y)
	composite.DrawRegion(dst, src, r, &op)
}
func (s *Scene) rnd() float64 {
	s.random = s.random*1664525 + 1013904223
	return float64(s.random) / 4294967296
}
func clearBlack(img *ebiten.Image) { img.Fill(color.Black) }
