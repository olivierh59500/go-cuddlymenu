// Package mobile binds the complete Cuddly production, starting with Calvin.
package mobile

import (
	"github.com/hajimehoshi/ebiten/v2"
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"
	"go-cuddlymenu/dck/mobilehost"
	"go-cuddlymenu/dck/screens"
)

var host = mobilehost.New("intro")

func init() {
	ebiten.SetTPS(screens.TicksPerSecond)
	ebiten.SetScreenClearedEveryFrame(false)
	enginemobile.SetGame(host)
}

// Configure supplies optional debug launch arguments before the first frame.
func Configure(screen string, metrics bool, startFrame int) {
	host.Configure(screen, metrics, startFrame)
}

// ConfigureAtRate selects the optional comparison rate before the first frame.
func ConfigureAtRate(screen string, metrics bool, startFrame, rate int) {
	host.ConfigureAtRate(screen, metrics, startFrame, rate)
}

// ConfigureCRT enables the optional menu material for device profiling.
func ConfigureCRT(enabled bool) { host.ConfigureMenuCRT(enabled) }

// ConfigureTour enables an unattended complete route with editable dwell
// times. Call it after ConfigureAtRate and before Android starts its view.
func ConfigureTour(introSeconds, screenSeconds, menuSeconds int) bool {
	return host.ConfigureTour(introSeconds, screenSeconds, menuSeconds)
}
