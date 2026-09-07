package menu

import (
	"image/color"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	maxLogicalWidth = 1280
	padButtonRadius = 42
	flyButtonRadius = 50
)

type controlState struct {
	Left  bool
	Right bool
	Fly   bool
}

type controlButton struct {
	X      float64
	Y      float64
	Radius float64
}

func (b controlButton) contains(x, y int) bool {
	dx := float64(x) - b.X
	dy := float64(y) - b.Y
	return dx*dx+dy*dy <= b.Radius*b.Radius
}

type controlLayout struct {
	Left  controlButton
	Right controlButton
	Fly   controlButton
}

func logicalWidth(outsideWidth, outsideHeight int) int {
	if outsideWidth <= 0 || outsideHeight <= 0 {
		return screenWidth
	}
	w := (outsideWidth*screenHeight + outsideHeight - 1) / outsideHeight
	if w < screenWidth {
		return screenWidth
	}
	if w > maxLogicalWidth {
		return maxLogicalWidth
	}
	return w
}

func makeControlLayout(width, height int) controlLayout {
	sideWidth := float64(width-screenWidth) / 2
	padX := 106.0
	flyX := float64(width) - 90
	if sideWidth >= 184 {
		padX = sideWidth / 2
		flyX = float64(width) - sideWidth/2
	}
	y := float64(height) - 88
	return controlLayout{
		Left:  controlButton{X: padX - 44, Y: y, Radius: padButtonRadius},
		Right: controlButton{X: padX + 44, Y: y, Radius: padButtonRadius},
		Fly:   controlButton{X: flyX, Y: y, Radius: flyButtonRadius},
	}
}

func (g *Game) readVirtualControls() (controlState, bool) {
	layout := makeControlLayout(g.layoutWidth, screenHeight)
	state := controlState{}

	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		x, y := ebiten.TouchPosition(id)
		state.press(layout, x, y)
	}
	if len(g.touchIDs) > 0 {
		g.touchSeen = true
	}

	mouseActive := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if mouseActive {
		x, y := ebiten.CursorPosition()
		state.press(layout, x, y)
	}

	g.controls = state
	return state, len(g.touchIDs) > 0 || mouseActive
}

func (s *controlState) press(layout controlLayout, x, y int) {
	if layout.Left.contains(x, y) {
		s.Left = true
	}
	if layout.Right.contains(x, y) {
		s.Right = true
	}
	if layout.Fly.contains(x, y) {
		s.Fly = true
	}
}

func (g *Game) virtualControlsVisible() bool {
	return runtime.GOOS == "android" || runtime.GOOS == "ios" || g.touchSeen || g.layoutWidth > screenWidth
}

func (g *Game) drawVirtualControls(dst *ebiten.Image) {
	layout := makeControlLayout(dst.Bounds().Dx(), dst.Bounds().Dy())
	drawRoundButton(dst, layout.Left, g.controls.Left)
	drawRoundButton(dst, layout.Right, g.controls.Right)
	drawRoundButton(dst, layout.Fly, g.controls.Fly)

	icon := color.RGBA{R: 255, G: 255, B: 255, A: 235}
	drawChevron(dst, layout.Left.X, layout.Left.Y, -1, icon)
	drawChevron(dst, layout.Right.X, layout.Right.Y, 1, icon)
	drawUpArrow(dst, layout.Fly.X, layout.Fly.Y-6, icon)
	ebitenutil.DebugPrintAt(dst, "VOLER", int(layout.Fly.X)-15, int(layout.Fly.Y)+23)
}

func drawRoundButton(dst *ebiten.Image, button controlButton, pressed bool) {
	fill := color.RGBA{R: 13, G: 22, B: 38, A: 170}
	border := color.RGBA{R: 126, G: 211, B: 255, A: 220}
	if pressed {
		fill = color.RGBA{R: 14, G: 116, B: 144, A: 225}
		border = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	vector.DrawFilledCircle(dst, float32(button.X), float32(button.Y), float32(button.Radius), fill, true)
	vector.StrokeCircle(dst, float32(button.X), float32(button.Y), float32(button.Radius), 3, border, true)
}

func drawChevron(dst *ebiten.Image, x, y float64, direction float32, clr color.Color) {
	cx := float32(x)
	cy := float32(y)
	vector.StrokeLine(dst, cx-direction*11, cy-15, cx+direction*7, cy, 7, clr, true)
	vector.StrokeLine(dst, cx+direction*7, cy, cx-direction*11, cy+15, 7, clr, true)
}

func drawUpArrow(dst *ebiten.Image, x, y float64, clr color.Color) {
	cx := float32(x)
	cy := float32(y)
	vector.StrokeLine(dst, cx, cy+18, cx, cy-16, 6, clr, true)
	vector.StrokeLine(dst, cx, cy-16, cx-12, cy-4, 6, clr, true)
	vector.StrokeLine(dst, cx, cy-16, cx+12, cy-4, 6, clr, true)
}
