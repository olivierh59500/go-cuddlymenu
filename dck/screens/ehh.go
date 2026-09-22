package screens

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
)

func (s *Scene) ehh() {
	background, main, logo, roller, inner, raster := s.asset("backdrop.png"), s.asset("main.png"), s.asset("toplogo.png"), s.asset("roller.png"), s.asset("roller_inner.png"), s.asset("toptextraster.png")
	font1, font2, font3 := s.asset("bannerfont.png"), s.asset("font2.png"), s.asset("font3.png")
	a, b, c, roll := s.surface(578, 64), s.surface(640, 32), s.surface(446, 10), s.surface(578, 66)
	s.filters[s.Canvas] = ebiten.FilterNearest
	r1, r2, r3 := s.ring(a, font1, 64, 64, 32, s.data.Strings["text1"], 6), s.ring(b, font2, 32, 32, 32, s.data.Strings["text2"], 6), s.ring(c, font3, 16, 10, 32, s.data.Strings["text3"], 4)
	var bars []*ebiten.Image
	var phases []float64
	for i := 1; i < 8; i++ {
		bars = append(bars, s.asset(fmt.Sprintf("r%d.png", 8-i)))
		phases = append(phases, float64(i*6))
	}
	curve := make([]float64, 250)
	for _, v := range digiCurve()[:763] {
		curve = append(curve, v/2)
	}
	for i := 0; i < 630; i++ {
		curve = append(curve, 20*math.Sin(float64(i)*.01))
	}
	for i := 0; i < 628; i++ {
		v := 20 * math.Sin(float64(i)*.01)
		if i < 100 || i >= 200 && i < 300 {
			v += 2 * math.Sin(float64(i))
		}
		curve = append(curve, v)
	}
	phase, bounce, bounceStep, rollX := 0.0, 0.0, 0.0, 0.0
	orbit := motion.DefaultNestedOrbit(motion.Point{X: 384, Y: 270}, motion.Point{X: 135, Y: 200})
	index := 0
	mode := 0
	show2, show3, stop := false, false, false
	s.render = func() {
		clearBlack(s.Canvas)
		a.Clear()
		b.Clear()
		c.Clear()
		phase += .006
		position := orbit.At(phase)
		s.transform(s.Canvas, background, position.X, position.Y, 1, 1, 0, float64(background.Bounds().Dx()/2), float64(background.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		s.draw(s.Canvas, main, 0, 0)
		for i, img := range bars {
			s.transform(s.Canvas, img, 0, 106+44*math.Cos(phases[i]/20), 77, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
			phases[i] += 1.2
		}
		if show2 {
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
		y := math.Abs(math.Sin(bounce) * 158)
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
		if show3 {
			r3.Step()
			r3.DrawAt(c, 0, 0)
			s.draw(s.Canvas, c, 160, 348-y)
			s.transform(s.Canvas, c, 160, 436-y, 1, -1, 0, 0, 0, 1, ebiten.BlendSourceOver)
		}
		switch r1.NextRune() {
		case '[':
			mode = 1
			show2 = true
		case '\\':
			mode = 2
			show3 = true
		case ']':
			mode = 3
		case '{':
			mode = 4
		}
		switch mode {
		case 1:
			bounceStep = .03
			stop = false
		case 2:
			bounceStep = .02
			stop = false
		case 3:
			bounceStep = .07
			stop = false
		case 4:
			stop = true
		}
		if bounce >= 3.1 {
			bounce = 0
			if stop {
				bounceStep = 0
			}
		}
		bounce += bounceStep
	}
}
