package screens

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/effects"
	"github.com/olivierh59500/democonstructionkit/presets"
)

func (s *Scene) doc() {
	background, mountains, raster, ball := s.asset("backdrop.png"), s.asset("mountains.png"), s.asset("rasters.png"), s.asset("ball.png")
	intro, outer, outerRows, inner, innerRows := s.surface(640, 50), s.surface(760, 50), s.surface(640, 50), s.surface(760, 50), s.surface(640, 120)
	r1 := s.ring(intro, s.asset("kh6.png"), "cuddly-doc", s.data.Strings["text1"], 10)
	r2 := s.ring(inner, s.asset("font_in.png"), "cuddly-doc", s.data.Strings["text2"], 10)
	r3 := s.ring(outer, s.asset("font_out.png"), "cuddly-doc", s.data.Strings["text2"], 10)
	var shadows []*ebiten.Image
	for _, name := range []string{"shadow1.png", "shadow2.png", "shadow3.png", "shadow4.png"} {
		shadows = append(shadows, s.asset(name))
	}
	stage := s.surface(384, 270)
	floor, err := effects.NewPerspectiveCheckerboard(presets.CuddlyDOCCheckerboard())
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, floor.Close)
	train, err := effects.NewProjectedBallTrain(presets.CuddlyDOCProjectedBalls(ball, shadows))
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, train.Close)
	s.filters[s.Canvas] = ebiten.FilterNearest
	white := s.surface(1, 1)
	white.Fill(color.White)
	curve := make([]float64, 389)
	for i := range curve {
		curve[i] = 20*math.Sin(float64(i)*(7.0/180*math.Pi)) + 30*math.Cos(float64(i)*(3.0/180*math.Pi))
	}
	for i := 0; i < 68; i++ {
		curve = append(curve, 30*math.Sin(float64(i)*(8.0/180*math.Pi)))
	}
	for i := 0; i < 189; i++ {
		curve = append(curve, 30*math.Sin(float64(i)*(8.0/180*math.Pi)))
	}
	jump, mainTick, stripTick := false, 0, 0
	vbl4 := 0.0
	s.music("", false)
	s.render = func() {
		clearBlack(s.Canvas)
		if !jump {
			intro.Clear()
			r1.Step()
			r1.DrawAt(intro, 0, 0)
			if r1.NextRune() == '\\' {
				jump = true
			}
			s.draw(s.Canvas, intro, 64, 62)
			return
		}
		if mainTick == 0 {
			s.music("Cuddly - 3D doc.ym", true)
		}
		stage.Clear()
		s.transform(s.Canvas, background, 0, 0, 77, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.draw(s.Canvas, mountains, 0, 0)
		outer.Clear()
		outerRows.Clear()
		inner.Clear()
		innerRows.Clear()
		r2.Step()
		r3.Step()
		r2.DrawAt(inner, 0, 0)
		r3.DrawAt(outer, 0, 0)
		vertical := 30 + 30*math.Cos(vbl4/20)
		for j := 0; j < 25; j++ {
			source := 64 + curve[(stripTick+j)%len(curve)]
			s.part(outerRows, outer, composite.Region{X: source, Y: float64(j * 2), Width: 640, Height: 2}, 0, float64(j*2), 1, 1)
			s.part(innerRows, inner, composite.Region{X: source, Y: float64(j * 2), Width: 640, Height: 2}, 0, float64(j*2)+vertical, 1, 1)
		}
		s.draw(s.Canvas, outerRows, 64, 62+vertical)
		black := ebiten.DrawImageOptions{Blend: ebiten.BlendSourceAtop}
		black.GeoM.Scale(640, 120)
		black.ColorScale.Scale(0, 0, 0, 1)
		innerRows.DrawImage(white, &black)
		s.transform(innerRows, raster, 0, 20, 75, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
		s.draw(s.Canvas, innerRows, 64, 62)
		vbl4 += 1.2
		stripTick++
		if err := floor.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		floor.Draw(stage)
		s.transform(s.Canvas, stage, 0, -128, 2, 2.6, 0, 0, 0, 1, ebiten.BlendSourceOver)
		// Preserve the reference's mixed seconds/frame arithmetic. Playback of
		// these logical ticks uses the application's selected fixed rate.
		if err := train.AdvanceAt(float64(mainTick) / 60); err != nil {
			s.err = err
			return
		}
		train.Draw(s.Canvas)
		mainTick++
	}
}
