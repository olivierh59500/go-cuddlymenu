package app

import (
	"fmt"
	"image"
	"image/color"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"go-cuddlymenu/dck/screens"
)

type controls struct {
	back, reset, chooser, confirm, left, right, up, down, rate bool
	chosen                                                     int
}
type button struct {
	label, action string
	rect          image.Rectangle
}

func (g *Game) touchEnabled() bool {
	return g.config.TouchControls || runtime.GOOS == "android" || runtime.GOOS == "ios"
}
func (g *Game) panelX() int { return max(24, (g.layoutWidth-768)/2+24) }
func (g *Game) buttons() []button {
	if !g.touchEnabled() {
		return nil
	}
	w, h := g.layoutWidth, g.layoutHeight
	if w == 0 {
		w, h = 768, 540
	}
	content := 768
	if g.scene != nil && g.transition == nil {
		content = g.scene.Descriptor.Width
	}
	x := max(62, (w-content)/4)
	makeButton := func(label, action string, cx, cy int) button {
		return button{label, action, image.Rect(cx-52, cy-27, cx+52, cy+27)}
	}
	b := []button{makeButton("SCREENS", "chooser", x, 50), makeButton("RESET", "reset", x, 124)}
	b = append(b, makeButton(fmt.Sprintf("%d HZ", g.TickRate()), "rate", x, 198))
	if g.scene != nil || g.transition != nil || g.chooser {
		b = append(b, makeButton("MENU", "back", x, h-70))
	}
	if g.scene != nil && g.scene.Descriptor.ID == "megaball" && !g.chooser && g.transition == nil {
		for i, s := range []struct{ label, action string }{{"PARAM +", "up"}, {"PARAM -", "down"}, {"VALUE +", "right"}, {"VALUE -", "left"}} {
			b = append(b, makeButton(s.label, s.action, w-x, 150+i*72))
		}
	}
	return b
}
func (g *Game) readControls() controls {
	in := controls{chosen: -1, chooser: inpututil.IsKeyJustPressed(ebiten.KeyF1), back: inpututil.IsKeyJustPressed(ebiten.KeyEscape) || (g.scene != nil && inpututil.IsKeyJustPressed(ebiten.KeySpace)), reset: inpututil.IsKeyJustPressed(ebiten.KeyR), confirm: inpututil.IsKeyJustPressed(ebiten.KeyEnter), left: inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), right: inpututil.IsKeyJustPressed(ebiten.KeyArrowRight), up: inpututil.IsKeyJustPressed(ebiten.KeyArrowUp), down: inpututil.IsKeyJustPressed(ebiten.KeyArrowDown)}
	in.rate = inpututil.IsKeyJustPressed(ebiten.KeyF3)
	if !g.touchEnabled() {
		return in
	}
	press := func(x, y int) {
		for _, b := range g.buttons() {
			if image.Pt(x, y).In(b.rect) {
				switch b.action {
				case "rate":
					in.rate = true
				case "back":
					in.back = true
				case "reset":
					in.reset = true
				case "chooser":
					in.chooser = true
				case "left":
					in.left = true
				case "right":
					in.right = true
				case "up":
					in.up = true
				case "down":
					in.down = true
				}
				return
			}
		}
		if g.chooser && x >= g.panelX() && x < g.panelX()+700 && y >= 54 {
			i := (y - 54) / 28
			if i < len(screens.Catalog()) {
				in.chosen = i
			}
		}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		x, y := ebiten.TouchPosition(id)
		press(x, y)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		press(x, y)
	}
	return in
}
func (g *Game) drawControls(dst *ebiten.Image) {
	for _, b := range g.buttons() {
		vector.FillRect(dst, float32(b.rect.Min.X), float32(b.rect.Min.Y), float32(b.rect.Dx()), float32(b.rect.Dy()), color.RGBA{20, 38, 52, 230}, false)
		vector.StrokeRect(dst, float32(b.rect.Min.X), float32(b.rect.Min.Y), float32(b.rect.Dx()), float32(b.rect.Dy()), 2, color.RGBA{126, 211, 255, 230}, false)
		ebitenutil.DebugPrintAt(dst, b.label, b.rect.Min.X+(b.rect.Dx()-len(b.label)*6)/2, b.rect.Min.Y+20)
	}
}
