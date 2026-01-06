package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type SineSprites struct {
	Tiles *TileSet
}

func (s *SineSprites) Draw(dst *ebiten.Image, t float64) {
	if s == nil || s.Tiles == nil {
		return
	}

	const (
		letterCount       = 12
		animationCount    = 7
		animationDuration = 8.0
		scrollDuration    = 1.0
		scrollIndex       = 0.5
	)

	cycle := animationDuration + (scrollDuration * 2)
	maxTime := float64(animationCount) * cycle
	if t >= maxTime {
		t = math.Mod(t, maxTime)
	}

	// Make sure we have enough space to subtract from time
	t += scrollDuration

	timerCycle := math.Floor(deriveFromTime(t, cycle, 0, cycle))
	xDisplacement := 0.0
	if timerCycle <= (scrollDuration - scrollIndex) {
		xDisplacement = -float64(dst.Bounds().Dx()) * deriveFromTime(t, scrollDuration, 0, 1)
	} else if timerCycle <= (scrollDuration + scrollDuration - scrollIndex) {
		xDisplacement = float64(dst.Bounds().Dx()) * deriveFromTime(t, scrollDuration, 1, 0)
	}

	animation := int(math.Floor((t-scrollDuration)/cycle)) % animationCount

	centerX := float64(dst.Bounds().Dx()) * 0.5
	centerY := float64(dst.Bounds().Dy()) * 0.5
	width := centerX * 0.9
	height := centerY * 0.88

	centerSpriteX := float64(carebearTileW) / 2
	centerSpriteY := float64(carebearTileH) / 2

	for i := 0; i < letterCount; i++ {
		p := sineSpritePoint(animation, t, i, width, height)
		p.x += centerX + xDisplacement
		p.y += centerY
		p.x = math.Floor(p.x) - centerSpriteX
		p.y = math.Floor(p.y) - centerSpriteY

		tile := s.Tiles.Tile(i)
		var op ebiten.DrawImageOptions
		op.GeoM.Translate(p.x, p.y)
		dst.DrawImage(tile, &op)
	}
}

type sinePoint struct {
	x float64
	y float64
}

func sineSpritePoint(animation int, t float64, i int, width, height float64) sinePoint {
	switch animation {
	case 0:
		return sineFunc0(t, i, width, height)
	case 1:
		return sineFunc1(t, i, width, height)
	case 2:
		return sineFunc2(t, i, width, height)
	case 3:
		return sineFunc3(t, i, width, height)
	case 4:
		return sineFunc4(t, i, width, height)
	case 5:
		return sineFunc5(t, i, width, height)
	default:
		return sineFunc6(t, i, width, height)
	}
}

func sineFunc0(t float64, i int, width, height float64) sinePoint {
	speedX := 25.0
	speedY := 4.5
	spacingY := 0.15
	spacingX := 0.3

	w := width * 0.5 * 0.5
	h := height * 0.25 * 0.25

	x := ((float64(i) - (12-1)/2.0) * spacingX) * w
	y := math.Cos((t-float64(i)*spacingY)*speedY) * h

	spinX := 6.1
	o := math.Sin(t*spinX + float64(i)*spinX*3)
	x += o * speedX

	x *= 2.1
	y *= 3

	return sinePoint{x: x, y: y}
}

func sineFunc1(t float64, i int, width, height float64) sinePoint {
	speedX := 2.5 * 0.7
	speedY := 5.0 * 0.7
	spacing := 0.1
	twistSpeed := 1.6 * 0.7
	twist := (math.Sin(t*twistSpeed) + math.Cos(t*twistSpeed)) * 0.75

	x := math.Sin((t-float64(i)*spacing)*speedX) * width * 1.05
	y := math.Sin((t-float64(i)*spacing)*speedY) * height * twist

	return sinePoint{x: x, y: y}
}

func sineFunc2(t float64, i int, width, height float64) sinePoint {
	speedX := 4.0
	speedY := 3.0
	spacing := 0.07
	twistSpeed := 2.0
	twist := (math.Sin(t*twistSpeed) + math.Cos(t*twistSpeed)) * 0.75

	x := math.Sin((t-float64(i)*spacing)*speedX) * width * 1.05
	y := math.Sin((t-float64(i)*spacing)*speedY) * height * twist
	x *= math.Cos((t - float64(i)) * 0.25)

	return sinePoint{x: x, y: y}
}

func sineFunc3(t float64, i int, width, height float64) sinePoint {
	speedX := 4.0 * 0.7
	speedY := 3.0 * 0.7
	spacing := 0.07
	twistSpeed := 2.0
	twist := (math.Sin(t*twistSpeed) + math.Cos(t*twistSpeed)) * 0.73

	x := math.Sin((t-float64(i)*spacing)*speedX) * width * twist
	y := math.Cos((t-float64(i)*spacing)*speedY) * height
	y *= math.Sin((t - float64(i)) * 0.25)

	return sinePoint{x: x, y: y}
}

func sineFunc4(t float64, i int, width, height float64) sinePoint {
	speedX := 1.0
	speedY := 1.2
	spacing := 0.1

	twistSpeedY := 4.0
	twistSpeedX := 4.5
	twistY := (1 - math.Sin((t-float64(i)*spacing)*twistSpeedY)) * 0.52
	twistX := (1 - math.Sin((t-float64(i)*spacing)*twistSpeedX)) * 0.52

	x := math.Sin((t-float64(i)*spacing)*speedX) * width * twistX
	y := math.Cos((t-float64(i)*spacing)*speedY) * height * twistY

	return sinePoint{x: x, y: y}
}

func sineFunc5(t float64, i int, width, height float64) sinePoint {
	speedX := 3.5 * 0.7
	speedY := 3.5 * 0.7
	spacing := 0.1
	h := height * 0.9
	twistSpeed := 3.0 * 0.5
	twist := (math.Sin(t*twistSpeed) + math.Cos(t*twistSpeed)) * 0.82

	x := math.Cos((t-float64(i)*spacing)*speedX) * width * 1.05
	y := math.Sin((t-float64(i)*spacing)*speedY) * h * twist

	return sinePoint{x: x, y: y}
}

func sineFunc6(t float64, i int, width, height float64) sinePoint {
	speedX := 4.0 * 0.7
	speedY := 3.0 * 0.7
	spacing := 0.1
	twistSpeed := 3.0
	twist := (math.Sin(t*twistSpeed) + math.Cos(t*twistSpeed)) * 0.75

	x := math.Cos((t-float64(i)*spacing)*speedX) * width * 1.04
	y := math.Sin((t-float64(i)*spacing)*speedY) * height * twist
	x *= math.Sin((t - float64(i)*spacing) * 0.5)

	return sinePoint{x: x, y: y}
}

func deriveFromTime(time, duration, min, max float64) float64 {
	if duration == 0 {
		return min
	}
	return ((math.Mod(time, duration) * (max - min)) / duration) + min
}
