package screens

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func (s *Scene) ehh() {
	background, main, logo, roller, inner, raster := s.asset("backdrop.png"), s.asset("main.png"), s.asset("toplogo.png"), s.asset("roller.png"), s.asset("roller_inner.png"), s.asset("toptextraster.png")
	font1, font2, font3 := s.asset("bannerfont.png"), s.asset("font2.png"), s.asset("font3.png")
	a, b, c, roll := s.surface(578, 64), s.surface(640, 32), s.surface(446, 10), s.surface(578, 66)
	s.filters[s.Canvas] = ebiten.FilterNearest
	r1, r2, r3 := s.ring(a, font1, "cuddly-ehh-main", s.data.Strings["text1"], 6), s.ring(b, font2, "cuddly-ehh-middle", s.data.Strings["text2"], 6), s.ring(c, font3, "cuddly-ehh-small", s.data.Strings["text3"], 4)
	var bars []*ebiten.Image
	for i := 1; i < 8; i++ {
		bars = append(bars, s.asset(fmt.Sprintf("r%d.png", 8-i)))
	}
	barGroup, err := sprites.NewTrain(sprites.TrainConfig{
		Images: bars, ScaleX: 77, Filter: ebiten.FilterNearest, Blend: ebiten.BlendSourceOver,
		Y: sprites.TrainAxis{Offset: 106, Wave: &motion.Wave{
			Amplitude: 44, Spatial: 6.0 / 20, Speed: 1.2 / 20, Phase: 6.0 / 20, Cos: true,
		}},
	})
	if err != nil {
		s.err = err
		return
	}
	curve, err := presets.CuddlyEhhhProfile()
	if err != nil {
		s.err = err
		return
	}
	rollerClock, err := motion.NewCuedWaveClock(presets.CuddlyEhhhRoller())
	if err != nil {
		s.err = err
		return
	}
	phase, rollX := 0.0, 0.0
	orbit := motion.DefaultNestedOrbit(motion.Point{X: 384, Y: 270}, motion.Point{X: 135, Y: 200})
	index := 0
	s.render = func() {
		clearBlack(s.Canvas)
		a.Clear()
		b.Clear()
		c.Clear()
		phase += .006
		position := orbit.At(phase)
		s.transform(s.Canvas, background, position.X, position.Y, 1, 1, 0, float64(background.Bounds().Dx()/2), float64(background.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		s.draw(s.Canvas, main, 0, 0)
		if err := barGroup.Update(kit.Frame{Time: float64(s.frame)}); err != nil {
			s.err = err
			return
		}
		barGroup.Draw(s.Canvas)
		if rollerClock.State().ShowSecond {
			r2.Step()
			r2.DrawAt(b, 0, 0)
		}
		s.transform(b, raster, 0, 0, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
		s.draw(s.Canvas, b, 64, 86)
		s.transform(s.Canvas, b, 64, 152, 1, -1, 0, 0, 0, 1, ebiten.BlendSourceOver)
		for row := 0; row < 170; row++ {
			s.part(s.Canvas, logo, composite.Region{Y: float64(row), Width: 767, Height: 1}, 384+curve[(index+row)%len(curve)]-383.5, float64(row)-.5, 1, 1)
		}
		index++
		y := rollerClock.At()
		s.draw(s.Canvas, roller, 64, 346-y)
		s.draw(roll, inner, rollX, 0)
		rollX -= 2
		if rollX <= -16 {
			rollX = 0
		}
		s.draw(s.Canvas, roll, 96, 360-y)
		r1.Step()
		r1.DrawAt(a, 0, 0)
		s.draw(s.Canvas, a, 96, 360-y)
		if rollerClock.State().ShowThird {
			r3.Step()
			r3.DrawAt(c, 0, 0)
			s.draw(s.Canvas, c, 160, 348-y)
			s.transform(s.Canvas, c, 160, 436-y, 1, -1, 0, 0, 0, 1, ebiten.BlendSourceOver)
		}
		if err := rollerClock.Step(r1.NextRune()); err != nil {
			s.err = err
			return
		}
	}
}
