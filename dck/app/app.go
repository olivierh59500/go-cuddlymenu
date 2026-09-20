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
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/olivierh59500/democonstructionkit/sound"
	device "github.com/olivierh59500/democonstructionkit/sound/ebiten"
	media "go-cuddlymenu/assets/cuddly"
	"go-cuddlymenu/dck/loader"
	"go-cuddlymenu/dck/menu"
	"go-cuddlymenu/dck/screens"
)

type Config struct {
	Screen        string
	Muted         bool
	TouchControls bool
}
type Game struct {
	menu                      *menu.Game
	scene                     *screens.Scene
	context                   *audio.Context
	player                    *device.Player
	config                    Config
	track, pending            string
	chooser                   bool
	selection                 int
	notice                    string
	noticeTicks               int
	loopMusic                 bool
	loaderPaused              bool
	transition                *loader.Screen
	target                    string
	layoutWidth, layoutHeight int
	touchIDs                  []ebiten.TouchID
}

func New(c Config) (*Game, error) {
	g := &Game{config: c, track: "menu/resources/menu.ym", loopMusic: true}
	g.menu = menu.NewGame()
	g.menu.UseExternalAudio()
	if c.TouchControls {
		g.menu.UseTouchControls()
	}
	g.menu.SetScreenHandler(func(name string) { g.pending = screens.DoorID(name) })
	if c.Screen != "" && c.Screen != "menu" {
		if err := g.open(c.Screen); err != nil {
			return nil, err
		}
	}
	return g, nil
}

// Begin loads a production through the original countdown. Intro and Reset have
// their own entry phases and do not use the floppy loader.
func (g *Game) Begin(id string) error {
	name := screens.DoorName(id)
	if name == "" {
		return g.open(id)
	}
	next, err := loader.New(name)
	if err != nil {
		return err
	}
	if g.transition != nil {
		g.transition.Close()
	}
	g.transition = next
	g.target = id
	g.setTrack(loader.Music)
	g.loopMusic = true
	return nil
}

func (g *Game) open(id string) error {
	if g.transition != nil {
		g.transition.Close()
		g.transition = nil
		g.target = ""
	}

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
	case ".wav":
		var decoded *wav.Stream
		decoded, err = wav.DecodeWithSampleRate(48000, bytes.NewReader(data))
		if err == nil {
			stream, err = sound.NewPCM16(decoded, sound.PCM16Options{SampleRate: 48000, Loop: g.loopMusic})
		}
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
	in := g.readControls()
	if in.chooser && g.transition == nil {
		g.chooser = !g.chooser
	}
	if in.back {
		if g.chooser {
			g.chooser = false
		} else if g.scene != nil || g.transition != nil {
			if err := g.Begin("menu"); err != nil {
				return err
			}
		}
	}
	if in.reset && !g.chooser {
		if err := g.open("reset"); err != nil {
			return err
		}
	}
	if g.transition != nil {
		if err := g.transition.Update(); err != nil {
			return err
		}
		if g.transition.Done() {
			if err := g.open(g.target); err != nil {
				return err
			}
		}
		if err := g.startAudio(); err != nil {
			return err
		}
		if g.transition != nil && g.player != nil {
			g.player.SetVolume(.7 * g.transition.Volume())
		}
		return nil
	}
	if g.chooser {
		list := screens.Catalog()
		if in.down {
			g.selection = (g.selection + 1) % len(list)
		}
		if in.up {
			g.selection = (g.selection + len(list) - 1) % len(list)
		}
		if in.chosen >= 0 {
			g.selection = in.chosen
			in.confirm = true
		}
		if in.confirm {
			if err := g.Begin(list[g.selection].ID); err != nil {
				g.notice = err.Error()
				g.noticeTicks = 150
			} else {
				g.chooser = false
			}
		}
		return g.startAudio()
	}
	var err error
	if g.scene != nil {
		g.scene.Input(screens.Input{Left: in.left, Right: in.right, Up: in.up, Down: in.down})
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
		if err = g.Begin(id); err != nil {
			return err
		}
	}
	return g.startAudio()
}
func (g *Game) Draw(dst *ebiten.Image) {
	if g.transition != nil {
		dst.Fill(color.Black)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(dst.Bounds().Dx()-768)/2, 0)
		dst.DrawImage(g.transition.Canvas, &op)
	} else if g.scene != nil {
		dst.Fill(color.Black)
		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(dst.Bounds().Dx()-g.scene.Descriptor.Width)/2, 0)
		dst.DrawImage(g.scene.Canvas, &op)
	} else {
		g.menu.Draw(dst)
	}
	if g.chooser {
		dst.Fill(color.RGBA{10, 14, 28, 255})
		ebitenutil.DebugPrintAt(dst, "CUDDLY DEMOS  /  Up, Down, Enter   /   Esc: back", g.panelX(), 20)
		for i, d := range screens.Catalog() {
			prefix := "  "
			if i == g.selection {
				prefix = "> "
			}
			suffix := ""
			if !d.Ready {
				suffix = " (coming later)"
			}
			ebitenutil.DebugPrintAt(dst, prefix+d.Title+suffix, g.panelX(), 54+i*28)
		}
	}
	g.drawControls(dst)
	if g.noticeTicks > 0 {
		ebitenutil.DebugPrintAt(dst, g.notice, 20, 10)
	}
}
func (g *Game) Layout(w, h int) (int, int) {
	width, height := 768, 540
	if g.transition != nil {
		width, height = loader.Width, loader.Height
	}
	if g.transition == nil {
		if g.scene != nil {
			width, height = g.scene.Layout(w, h)
		} else {
			width, height = g.menu.Layout(w, h)
		}
	}
	if g.touchEnabled() && h > 0 {
		width = max(width, min(1280, (w*height+h-1)/h))
	}
	g.layoutWidth, g.layoutHeight = width, height
	return width, height
}
func (g *Game) Close() error {
	if g.transition != nil {
		g.transition.Close()
	}
	if g.player != nil {
		g.player.Close()
	}
	if g.scene != nil {
		g.scene.Close()
	}
	return g.menu.Close()
}

func (g *Game) CurrentScreen() string {
	if g.transition != nil {
		return "loader:" + g.target
	}
	if g.scene != nil {
		return g.scene.Descriptor.ID
	}
	return "menu"
}

// AdvanceScene supports deterministic profiling at late animation phases.
// It does not poll user input or open the audio device.
func (g *Game) AdvanceScene(ticks int) error {
	if g.scene == nil {
		return nil
	}
	for i := 0; i < ticks; i++ {
		if err := g.scene.Update(); err != nil {
			return err
		}
	}
	g.consumeAudioCue()
	return nil
}
