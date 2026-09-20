// Package app connects the original Cuddly menu to native DCK screens and audio.
package app

import (
	"bytes"
	"fmt"
	"image/color"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/olivierh59500/democonstructionkit/sound"
	device "github.com/olivierh59500/democonstructionkit/sound/ebiten"
	media "go-cuddlymenu/assets/cuddly"
	"go-cuddlymenu/dck/menu"
	"go-cuddlymenu/dck/screens"
)

type Config struct {
	Screen string
	Muted  bool
}
type Game struct {
	menu           *menu.Game
	scene          *screens.Scene
	context        *audio.Context
	player         *device.Player
	config         Config
	track, pending string
	chooser        bool
	selection      int
	notice         string
	noticeTicks    int
	loopMusic      bool
	loaderPaused   bool
}

func New(c Config) (*Game, error) {
	g := &Game{config: c, track: "menu/resources/menu.ym", loopMusic: true}
	g.menu = menu.NewGame()
	g.menu.UseExternalAudio()
	g.menu.SetScreenHandler(func(name string) { g.pending = screens.DoorID(name) })
	if c.Screen != "" && c.Screen != "menu" {
		if err := g.open(c.Screen); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (g *Game) open(id string) error {
	if id == "menu" || id == "" {
		if g.scene != nil {
			g.scene.Close()
			g.scene = nil
		}
		g.setTrack("menu/resources/menu.ym")
		g.loopMusic = true
		return nil
	}
	s, err := screens.New(id)
	if err != nil {
		return err
	}
	if g.scene != nil {
		g.scene.Close()
	}
	g.scene = s
	g.consumeAudioCue()
	return nil
}
func (g *Game) consumeAudioCue() {
	if g.scene == nil {
		return
	}
	cue := g.scene.TakeAudioCue()
	if cue == nil {
		return
	}
	track := ""
	if cue.File != "" {
		track = "" + g.scene.Descriptor.Directory + "/" + cue.File
		if strings.HasPrefix(cue.File, "@") {
			track = cue.File[1:]
		}
	}
	g.setTrack(track)
	g.loopMusic = cue.Loop
}
func (g *Game) setTrack(name string) {
	if g.player != nil {
		g.player.Close()
		g.player = nil
	}
	g.track = name
	g.loaderPaused = false
}

func (g *Game) startAudio() error {
	if g.config.Muted || g.player != nil || g.track == "" {
		return nil
	}
	if g.context == nil {
		g.context = audio.CurrentContext()
		if g.context == nil {
			g.context = audio.NewContext(48000)
		}
	}
	data, err := media.Files.ReadFile(g.track)
	if err != nil {
		return err
	}
	var stream *sound.Stream
	switch strings.ToLower(path.Ext(g.track)) {
	case ".ym":
		stream, err = sound.NewYM(data, sound.YMOptions{SampleRate: 48000, Loop: g.loopMusic})
	case ".mp3":
		var decoded *mp3.Stream
		decoded, err = mp3.DecodeWithSampleRate(48000, bytes.NewReader(data))
		if err == nil {
			stream, err = sound.NewPCM16(decoded, sound.PCM16Options{SampleRate: 48000, Loop: g.loopMusic})
		}
	default:
		return fmt.Errorf("unsupported audio asset %q", g.track)
	}
	if err != nil {
		return err
	}
	g.player, err = device.NewPlayer(g.context, stream)
	if err != nil {
		stream.Close()
		return err
	}
	g.player.SetVolume(.7)
	g.player.Play()
	return nil
}

func (g *Game) Update() error {
	if g.noticeTicks > 0 {
		g.noticeTicks--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.chooser = !g.chooser
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if g.chooser {
			g.chooser = false
		} else {
			g.open("menu")
		}
	}
	if g.scene != nil && !g.chooser && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.open("menu")
	}
	if g.scene != nil && !g.chooser && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if err := g.open("reset"); err != nil {
			return err
		}
	}
	if g.chooser {
		list := screens.Catalog()
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.selection = (g.selection + 1) % len(list)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.selection = (g.selection + len(list) - 1) % len(list)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if err := g.open(list[g.selection].ID); err != nil {
				g.notice = err.Error()
				g.noticeTicks = 180
			} else {
				g.chooser = false
			}
		}
		return g.startAudio()
	}
	var err error
	if g.scene != nil {
		g.scene.Input(screens.Input{Left: inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), Right: inpututil.IsKeyJustPressed(ebiten.KeyArrowRight), Up: inpututil.IsKeyJustPressed(ebiten.KeyArrowUp), Down: inpututil.IsKeyJustPressed(ebiten.KeyArrowDown)})
		err = g.scene.Update()
		g.consumeAudioCue()
	} else {
		err = g.menu.Update()
	}
	if err != nil {
		return err
	}
	if g.pending != "" {
		id := g.pending
		g.pending = ""
		if err := g.open(id); err != nil {
			g.notice = err.Error()
			g.noticeTicks = 180
		}
	}
	if err = g.startAudio(); err != nil {
		return err
	}
	if g.player != nil {
		if g.scene == nil && g.menu.IsLoading() {
			g.player.Pause()
			g.loaderPaused = true
		} else if g.loaderPaused {
			g.player.Play()
			g.loaderPaused = false
		}
	}
	return nil
}
func (g *Game) Draw(dst *ebiten.Image) {
	if g.scene != nil {
		g.scene.Draw(dst)
	} else {
		g.menu.Draw(dst)
	}
	if g.chooser {
		dst.Fill(color.RGBA{10, 14, 28, 255})
		ebitenutil.DebugPrintAt(dst, "CUDDLY DEMOS  /  Up, Down, Enter   /   Esc: back", 24, 20)
		for i, d := range screens.Catalog() {
			prefix := "  "
			if i == g.selection {
				prefix = "> "
			}
			suffix := ""
			if !d.Ready {
				suffix = " (coming later)"
			}
			ebitenutil.DebugPrintAt(dst, prefix+d.Title+suffix, 24, 54+i*24)
		}
	}
	if g.noticeTicks > 0 {
		ebitenutil.DebugPrintAt(dst, g.notice, 20, 10)
	}
}
func (g *Game) Layout(w, h int) (int, int) {
	if g.scene != nil {
		return g.scene.Layout(w, h)
	}
	return g.menu.Layout(w, h)
}
func (g *Game) Close() error {
	if g.player != nil {
		g.player.Close()
	}
	if g.scene != nil {
		g.scene.Close()
	}
	return g.menu.Close()
}
