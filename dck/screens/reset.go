package screens

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

func (s *Scene) reset() {
	backdrop, raster, fontW, fontP, fontC := s.asset("backdrop.png"), s.asset("raster6.png"), s.asset("fontw.png"), s.asset("fontp.png"), s.asset("cuddlyfont1.png")
	upA, downA, upB, downB := s.asset("raster2_up.png"), s.asset("raster2_down.png"), s.asset("raster1_up.png"), s.asset("raster1_down.png")
	back, plain, scroll, merge, rotated, off, bars := s.surface(640, 400), s.surface(640, 400), s.surface(640, 400), s.surface(640, 400), s.surface(640, 400), s.surface(640, 400), s.surface(768, 540)
	r1, r2 := s.ring(plain, fontW, "cuddly-reset", s.data.Strings["text1"], 6), s.ring(plain, fontW, "cuddly-reset", s.data.Strings["text2"], 6)
	r3, r4 := s.ring(scroll, fontP, "cuddly-reset", s.data.Strings["text3"], 6), s.ring(scroll, fontP, "cuddly-reset", s.data.Strings["text4"], 6)
	grid := s.bitmap(fontC, "cuddly-reset-letters", ebiten.FilterLinear)
	slotConfig := presets.CuddlyResetSlots(grid, s.data.Strings["text5"])
	letters, err := scrolling.New(scrolling.Config{Slots: &slotConfig})
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, letters.Close)
	phases := make([]float64, 16)
	for i := range phases {
		phases[i] = float64(i) * .25
	}
	part, time, index := 1, 0, 12.0
	offset, yy, fade := 0.0, 0.0, 0.0
	s.music("", false)
	drawBack := func() {
		angle := 2 * math.Pi / 384 * index
		x := 335 - 220*math.Sin(angle)*math.Sin(offset)
		y := 250 + 170*math.Sin(angle+math.Pi/3)
		s.transform(back, backdrop, x, y, 1, 1, 0, float64(backdrop.Bounds().Dx()/2), float64(backdrop.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		index++
		offset = math.Mod(offset+math.Pi/32, 2*math.Pi)
	}
	drawBars := func() {
		pair := func(i int, a, b *ebiten.Image, alpha float64) {
			y := 270 + 200*math.Cos(phases[i])
			s.transform(bars, a, 100, y, 2, 1, 0, float64(a.Bounds().Dx()/2), float64(a.Bounds().Dy()/2), alpha, ebiten.BlendSourceOver)
			s.transform(bars, b, 668, y, 2, 1, 0, float64(b.Bounds().Dx()/2), float64(b.Bounds().Dy()/2), alpha, ebiten.BlendSourceOver)
		}
		for i, a := range phases {
			if a > math.Pi && a < 2*math.Pi {
				pair(i, downA, downB, .5)
			}
		}
		for i, a := range phases {
			if a >= 0 && a <= math.Pi/2 {
				pair(i, upA, upB, 1)
			}
		}
		for i := len(phases) - 1; i >= 0; i-- {
			a := phases[i]
			if a > math.Pi/2 && a <= math.Pi {
				pair(i, upA, upB, 1)
			}
			phases[i] += .05
			if phases[i] >= 2*math.Pi {
				phases[i] -= 2 * math.Pi
			}
		}
	}
	maskScroll := func() {
		s.draw(merge, scroll, 0, 70)
		s.transform(merge, off, 0, 0, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
		s.draw(s.Canvas, merge, 64, 0)
	}
	s.render = func() {
		clearBlack(s.Canvas)
		for _, img := range []*ebiten.Image{back, plain, scroll, merge, rotated, bars} {
			img.Clear()
		}
		s.draw(off, raster, 0, yy)
		yy -= 2
		if yy <= -480 {
			yy = 3
		}
		// These are independent checks: a control can start the next part in
		// this same tick, as in the original screen.
		if part == 1 {
			r1.Step()
			r1.DrawAt(plain, 0, 184)
			s.draw(s.Canvas, plain, 64, 70)
			if r1.NextRune() == '\\' {
				s.music("cuddlyreset.ym", true)
				part++
			}
		}
		if part == 2 {
			r2.Step()
			r2.DrawAt(plain, 0, 184)
			s.draw(s.Canvas, plain, 64, 70)
			if r2.NextRune() == ']' {
				part++
			}
		}
		if part == 3 {
			r3.Step()
			r3.DrawAt(scroll, 0, 184)
			maskScroll()
			if r3.NextRune() == '\\' {
				part++
			}
		}
		if part == 4 {
			drawBack()
			drawBars()
			if time <= 40 {
				s.transform(s.Canvas, back, 64, 70, 1, 1, 0, 0, 0, math.Min(1, fade), ebiten.BlendSourceOver)
				fade += .025
			}
			if time == 41 {
				fade = 0
			}
			if time >= 41 && time <= 80 {
				s.transform(s.Canvas, bars, 0, 40, 1, .85, 0, 0, 0, math.Min(1, fade), ebiten.BlendSourceOver)
				s.draw(s.Canvas, back, 64, 70)
				fade += .025
			}
			if time >= 81 {
				s.transform(s.Canvas, bars, 0, 40, 1, .85, 0, 0, 0, 1, ebiten.BlendSourceOver)
				s.draw(s.Canvas, back, 64, 70)
			}
			time++
			if r4.NextRune() == ']' {
				part++
			}
			r4.Step()
			r4.DrawAt(scroll, 0, 184)
			maskScroll()
		}
		if part == 5 {
			drawBack()
			drawBars()
			s.transform(s.Canvas, bars, 0, 40, 1, .85, 0, 0, 0, 1, ebiten.BlendSourceOver)
			s.draw(s.Canvas, back, 64, 70)
			if err := letters.Update(kit.Frame{}); err != nil {
				s.err = err
				return
			}
			letters.Draw(rotated)
			s.draw(merge, rotated, 0, 0)
			s.transform(merge, off, 0, 0, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
			s.draw(s.Canvas, merge, 64, 0)
		}
	}
}
