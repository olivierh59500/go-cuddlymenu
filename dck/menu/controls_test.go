package menu

import (
	"runtime"
	"testing"
)

func TestLogicalWidth(t *testing.T) {
	tests := []struct {
		name          string
		outsideWidth  int
		outsideHeight int
		want          int
	}{
		{name: "invalid dimensions", want: screenWidth},
		{name: "portrait keeps minimum", outsideWidth: 1080, outsideHeight: 2400, want: screenWidth},
		{name: "pixel landscape aspect", outsideWidth: 2400, outsideHeight: 1080, want: 1192},
		{name: "ultrawide is capped", outsideWidth: 4000, outsideHeight: 1000, want: maxLogicalWidth},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logicalWidth(tt.outsideWidth, tt.outsideHeight); got != tt.want {
				t.Fatalf("logicalWidth(%d, %d) = %d, want %d", tt.outsideWidth, tt.outsideHeight, got, tt.want)
			}
		})
	}
}

func TestWideControlLayoutUsesSideAreas(t *testing.T) {
	const width = 1192
	layout := makeControlLayout(width, screenHeight)
	sceneLeft := float64(width-screenWidth) / 2
	sceneRight := sceneLeft + screenWidth

	if layout.Left.X-layout.Left.Radius < 0 || layout.Right.X+layout.Right.Radius > sceneLeft {
		t.Fatalf("direction pad is not contained in left side area: %#v", layout)
	}
	if layout.Fly.X-layout.Fly.Radius < sceneRight || layout.Fly.X+layout.Fly.Radius > width {
		t.Fatalf("fly button is not contained in right side area: %#v", layout.Fly)
	}
}

func TestControlStateSupportsMultitouch(t *testing.T) {
	layout := makeControlLayout(1192, screenHeight)
	state := controlState{}
	state.press(layout, int(layout.Right.X), int(layout.Right.Y))
	state.press(layout, int(layout.Fly.X), int(layout.Fly.Y))

	if state.Left || !state.Right || !state.Fly {
		t.Fatalf("unexpected control state: %#v", state)
	}
}

func TestVirtualControlsAreDisabledOnDesktop(t *testing.T) {
	if runtime.GOOS == "android" || runtime.GOOS == "ios" {
		t.Skip("mobile platform")
	}
	if virtualControlsEnabled() {
		t.Fatal("virtual controls must stay disabled on desktop")
	}
}

func TestEnterButtonIsIndependentOfFlight(t *testing.T) {
	layout := makeControlLayout(1212, screenHeight)
	state := controlState{}
	state.press(layout, int(layout.Load.X), int(layout.Load.Y))
	if !state.Load || state.Fly || state.Left || state.Right {
		t.Fatalf("incorrect enter input: %+v", state)
	}
}
