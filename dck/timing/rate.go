// Package timing defines the fixed animation rates supported by the production.
package timing

import "fmt"

// DefaultRate matches NATIVE's 60-step pauses and nominal animation timebase.
// It is independent of the display refresh rate and the audio sample clock.
const DefaultRate = 60

func Normalize(rate int) (int, error) {
	if rate == 0 {
		rate = DefaultRate
	}
	if rate != 50 && rate != 60 {
		return 0, fmt.Errorf("animation rate must be 50 or 60 Hz, got %d", rate)
	}
	return rate, nil
}
