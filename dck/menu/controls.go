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
	Load  bool
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
	Load  controlButton
}

type controlSprites struct {
	left  [2]*ebiten.Image
	right [2]*ebiten.Image
	fly   [2]*ebiten.Image
	load  [2]*ebiten.Image
}

func newControlSprites() controlSprites {
	return controlSprites{
		left: [2]*ebiten.Image{
			newControlSprite(padButtonRadius, false, -1, false),
			newControlSprite(padButtonRadius, true, -1, false),
		},
		right: [2]*ebiten.Image{
			newControlSprite(padButtonRadius, false, 1, false),
			newControlSprite(padButtonRadius, true, 1, false),
		},
		fly: [2]*ebiten.Image{
			newControlSprite(flyButtonRadius, false, 0, true),
			newControlSprite(flyButtonRadius, true, 0, true),
		},
		load: [2]*ebiten.Image{newLoadButton(false), newLoadButton(true)},
	}
}

func newLoadButton(pressed bool) *ebiten.Image {
	img := ebiten.NewImage(108, 108)
	drawRoundButton(img, controlButton{X: 54, Y: 54, Radius: 50}, pressed)
	ebitenutil.DebugPrintAt(img, "ENTER", 39, 48)
	return img
}

func newControlSprite(radius float64, pressed bool, direction float32, fly bool) *ebiten.Image {
	size := int(radius*2) + 8
	center := float64(size) / 2
	img := ebiten.NewImage(size, size)
	drawRoundButton(img, controlButton{X: center, Y: center, Radius: radius}, pressed)

	icon := color.RGBA{R: 255, G: 255, B: 255, A: 235}
	if fly {
		drawUpArrow(img, center, center-6, icon)
		ebitenutil.DebugPrintAt(img, "VOLER", int(center)-15, int(center)+23)
	} else {
		drawChevron(img, center, center, direction, icon)
	}
	return img
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
		Load:  controlButton{X: flyX, Y: y - 118, Radius: flyButtonRadius},
	}
}

func (g *Game) readVirtualControls() (controlState, bool) {
	if !g.virtualControlsVisible() {
		g.controls = controlState{}
		return g.controls, false
	}

	layout := makeControlLayout(g.layoutWidth, screenHeight)
	state := controlState{}

	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		x, y := ebiten.TouchPosition(id)
		state.press(layout, x, y)
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
	if layout.Load.contains(x, y) {
		s.Load = true
	}
}

func (g *Game) virtualControlsVisible() bool {
	return g.touchControls || virtualControlsEnabled()
}

func virtualControlsEnabled() bool {
	return runtime.GOOS == "android" || runtime.GOOS == "ios"
}

func (g *Game) drawVirtualControls(dst *ebiten.Image) {
	layout := makeControlLayout(dst.Bounds().Dx(), dst.Bounds().Dy())
	drawControlSprite(dst, g.controlUI.left[boolIndex(g.controls.Left)], layout.Left)
	drawControlSprite(dst, g.controlUI.right[boolIndex(g.controls.Right)], layout.Right)
	drawControlSprite(dst, g.controlUI.fly[boolIndex(g.controls.Fly)], layout.Fly)
	drawControlSprite(dst, g.controlUI.load[boolIndex(g.controls.Load)], layout.Load)
}

func boolIndex(value bool) int {
	if value {
		return 1
	}
	return 0
}

func drawControlSprite(dst, sprite *ebiten.Image, button controlButton) {
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(
		button.X-float64(sprite.Bounds().Dx())/2,
		button.Y-float64(sprite.Bounds().Dy())/2,
	)
	dst.DrawImage(sprite, &op)
}

func drawRoundButton(dst *ebiten.Image, button controlButton, pressed bool) {
	fill := color.RGBA{R: 13, G: 22, B: 38, A: 170}
	border := color.RGBA{R: 126, G: 211, B: 255, A: 220}
	if pressed {
		fill = color.RGBA{R: 14, G: 116, B: 144, A: 225}
		border = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	vector.FillCircle(dst, float32(button.X), float32(button.Y), float32(button.Radius), fill, true)
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
