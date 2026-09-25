package menu

import "fmt"

// NavigateTo uses the original walking and thrusting autopilot to enter a door.
func (g *Game) NavigateTo(name string) error {
	for i, door := range demoScreens {
		if door.Name == name {
			g.autoPilot = AutoPilot{NextScreen: i, WaitToLoad: 80}
			return nil
		}
	}
	return fmt.Errorf("menu: unknown door %q", name)
}

// DelayAutoPilot leaves a readable pause after returning from a screen.
func (g *Game) DelayAutoPilot(ticks int) { g.autoPilot.ActivateIn = ticks }

// UseExternalAudio lets a screen controller own the single application player.
// Call before the first Update. Standalone menu behavior remains the default.
func (g *Game) UseExternalAudio()                     { g.audioReady = true }
func (g *Game) SetScreenHandler(handler func(string)) { g.screenHandler = handler }
func (g *Game) IsLoading() bool                       { return g.loading.Active }

func (g *Game) UseTouchControls() {
	if !g.virtualControlsVisible() {
		g.controlUI = newControlSprites()
	}
	g.touchControls = true
}

func (g *Game) Close() error {
	if g.background != nil {
		g.background.Close()
		g.background = nil
	}
	if g.audioPlayer != nil {
		g.audioPlayer.Close()
		g.audioPlayer = nil
	}
	if g.musicStream != nil {
		g.musicStream.Close()
		g.musicStream = nil
	}
	return nil
}
