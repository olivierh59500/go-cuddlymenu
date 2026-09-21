package app

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"go-cuddlymenu/dck/screens"
)

// TourOptions controls the unattended, complete tour. Durations only count
// time inside a scene; loaders and walking between doors retain their timing.
type TourOptions struct {
	IntroDuration  time.Duration
	ScreenDuration time.Duration
	MenuDuration   time.Duration
}

func DefaultTourOptions() TourOptions {
	return TourOptions{IntroDuration: time.Minute, ScreenDuration: time.Minute, MenuDuration: 4 * time.Second}
}

// The route follows the menu's doors from left to right, then visits Reset.
var tourScreens = []string{"big-sprite", "colorshock", "ehh", "megascroller", "digi", "spreadpoint", "led", "3d-doc", "fullscreen", "starwars", "knucklebuster", "dna", "megaball"}

type Tour struct {
	*Game
	options      TourOptions
	phase        string
	index, ticks int
}

func NewTour(c Config, options TourOptions) (*Tour, error) {
	if options.IntroDuration <= 0 || options.ScreenDuration <= 0 || options.MenuDuration <= 0 {
		return nil, fmt.Errorf("tour: durations must be positive")
	}
	c.Screen = "intro"
	g, err := New(c)
	if err != nil {
		return nil, err
	}
	return &Tour{Game: g, options: options, phase: "intro"}, nil
}

func (t *Tour) RecordingChapter() string { return t.CurrentScreen() }

func (t *Tour) Update() error {
	if err := t.Game.Update(); err != nil {
		return err
	}
	id := t.CurrentScreen()
	if t.transition != nil {
		return nil
	}
	t.ticks++
	due := func(d time.Duration) bool { return int64(t.ticks)*int64(time.Second) >= int64(d)*int64(t.TickRate()) }
	returnToMenu := func() error {
		t.phase, t.ticks = "menu", 0
		t.menu.DelayAutoPilot(int(t.options.MenuDuration.Seconds()*float64(t.TickRate())) + 2)
		return t.Begin("menu")
	}
	switch t.phase {
	case "intro":
		if due(t.options.IntroDuration) {
			return returnToMenu()
		}
	case "menu":
		if id != "menu" {
			return fmt.Errorf("tour: expected menu, got %s", id)
		}
		if due(t.options.MenuDuration) {
			t.ticks = 0
			if t.index < len(tourScreens) {
				t.phase = "walking"
				return t.menu.NavigateTo(screens.DoorName(tourScreens[t.index]))
			}
			if t.index == len(tourScreens) {
				t.phase = "reset"
				return t.Begin("reset")
			}
			return ebiten.Termination
		}
	case "walking":
		if id == tourScreens[t.index] {
			t.phase, t.ticks = "screen", 0
		} else if id != "menu" {
			return fmt.Errorf("tour: unexpected screen %s", id)
		}
		if t.ticks > 120*t.TickRate() {
			return fmt.Errorf("tour: navigation to %s timed out", tourScreens[t.index])
		}
	case "screen", "reset":
		if due(t.options.ScreenDuration) {
			t.index++
			return returnToMenu()
		}
	}
	return nil
}
