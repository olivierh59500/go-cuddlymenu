package screens

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/timeline"
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
	director, err := timeline.NewStageSequence(presets.CuddlyResetDirector())
	if err != nil {
		s.err = err
		return
	}
	rasterOrbit, err := composite.NewPairedRasterOrbit(presets.CuddlyResetRasterOrbit(upA, upB, downA, downB))
	if err != nil {
		s.err = err
		return
	}
	rasterScroll, err := motion.NewWrapBank(presets.CuddlyResetRasterScroll())
	if err != nil {
		s.err = err
		return
	}
	index := 12.0
	offset := 0.0
	s.music("", false)
	drawBack := func() {
		angle := 2 * math.Pi / 384 * index
		x := 335 - 220*math.Sin(angle)*math.Sin(offset)
		y := 250 + 170*math.Sin(angle+math.Pi/3)
		s.transform(back, backdrop, x, y, 1, 1, 0, float64(backdrop.Bounds().Dx()/2), float64(backdrop.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		index++
		offset = math.Mod(offset+math.Pi/32, 2*math.Pi)
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
		s.draw(off, raster, 0, rasterScroll.At(0))
		rasterScroll.Step()
		// These are independent checks: a control can start the next part in
		// this same tick, as in the original screen.
		if director.Active("plain-a") {
			r1.Step()
			r1.DrawAt(plain, 0, 184)
			s.draw(s.Canvas, plain, 64, 70)
			if r1.NextRune() == '\\' {
				s.music("cuddlyreset.ym", true)
				director.Trigger("first-end")
			}
		}
		if director.Active("plain-b") {
			r2.Step()
			r2.DrawAt(plain, 0, 184)
			s.draw(s.Canvas, plain, 64, 70)
			if r2.NextRune() == ']' {
				director.Trigger("second-end")
			}
		}
		if director.Active("masked-a") {
			r3.Step()
			r3.DrawAt(scroll, 0, 184)
			maskScroll()
			if r3.NextRune() == '\\' {
				director.Trigger("third-end")
			}
		}
		if director.Active("showcase") {
			drawBack()
			rasterOrbit.Step()
			rasterOrbit.Draw(bars)
			phase, alpha := director.Window()
			switch phase {
			case 0:
				s.transform(s.Canvas, back, 64, 70, 1, 1, 0, 0, 0, alpha, ebiten.BlendSourceOver)
			case 1:
				s.transform(s.Canvas, bars, 0, 40, 1, .85, 0, 0, 0, alpha, ebiten.BlendSourceOver)
				s.draw(s.Canvas, back, 64, 70)
			case 2:
				s.transform(s.Canvas, bars, 0, 40, 1, .85, 0, 0, 0, 1, ebiten.BlendSourceOver)
				s.draw(s.Canvas, back, 64, 70)
			}
			director.StepWindow()
			if r4.NextRune() == ']' {
				director.Trigger("fourth-end")
			}
			r4.Step()
			r4.DrawAt(scroll, 0, 184)
			maskScroll()
		}
		if director.Active("final") {
			drawBack()
			rasterOrbit.Step()
			rasterOrbit.Draw(bars)
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
