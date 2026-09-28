package screens

import (
	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func (s *Scene) dna() {
	main := s.surface(416, 276)
	s.filters[s.Canvas] = ebiten.FilterNearest
	introText, front, back, sineText, sineWave := s.surface(320, 25), s.surface(320, 25), s.surface(320, 25), s.surface(320, 16), s.surface(320, 100)
	intro := s.ring(introText, s.asset("font_intro.png"), "cuddly-dna", s.data.Strings["text_intro"], 6)
	topFront := s.ring(front, s.asset("font_top_front.png"), "cuddly-dna", s.data.Strings["text_top"], 6)
	topBack := s.ring(back, s.asset("font_top_back.png"), "cuddly-dna", s.data.Strings["text_top"], 6)
	sine := s.ring(sineText, s.asset("font_sin.png"), "cuddly-dna-sine", s.data.Strings["text_sin"], 3)
	// Preserve the authored texture upload order for fractional logo sampling.
	dnaFontFront, dnaFontBack := s.asset("font_dna_front.png"), s.asset("font_dna_back.png")
	ribbon, err := composite.NewTwistingRibbon(presets.CuddlyDNARibbon(front, back))
	if err != nil {
		s.err = err
		return
	}
	gradientFront, gradientBack, logo := s.asset("gradient_dna_front.png"), s.asset("gradient_dna_back.png"), s.asset("tcb.png")
	feedback := func(font *ebiten.Image, gradient *ebiten.Image, direction int) *scrolling.Scrolling {
		config, err := presets.CuddlyDNAFeedbackScroll(
			s.bitmap(font, "cuddly-dna", ebiten.FilterLinear), gradient,
			s.data.Strings["text_dna"], s.data.Numbers["dna_pos"], direction)
		if err != nil {
			s.err = err
			return nil
		}
		scroll, err := scrolling.New(config)
		if err != nil {
			s.err = err
			return nil
		}
		s.closeEffects = append(s.closeEffects, scroll.Close)
		return scroll
	}
	fFront := feedback(dnaFontFront, gradientFront, 1)
	fBack := feedback(dnaFontBack, gradientBack, -1)
	outerConfig, err := presets.CuddlyDNAOuterLogoRows(logo)
	if err != nil {
		s.err = err
		return
	}
	outerRows, err := composite.NewSampledRows(outerConfig)
	if err != nil {
		s.err = err
		return
	}
	centerConfig, err := presets.CuddlyDNACenterLogoRows(logo)
	if err != nil {
		s.err = err
		return
	}
	centerRows, err := composite.NewSampledRows(centerConfig)
	if err != nil {
		s.err = err
		return
	}
	wave := composite.WaveStrips{Axis: composite.Columns, Thickness: 1, Filter: ebiten.FilterLinear, Waves: []composite.StripWave{{Amplitude: 30, Spatial: .004, Speed: .04}}}
	points, err := geometry.SphereCloud(geometry.SphereCloudConfig{
		Count: 125, Radius: 100, Sampling: geometry.SphereRandomAngles,
		NextFloat: s.rnd,
	})
	if err != nil {
		s.err = err
		return
	}
	cloud, err := sprites.NewRotatingDiscCloud(presets.CuddlyDNADiscCloud(points))
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, cloud.Close)
	iteration := 0
	inIntro := true
	s.music("", false)
	s.render = func() {
		clearBlack(main)
		if inIntro {
			introText.Clear()
			s.advanceScroll(intro)
			s.drawScroll(intro, introText, 0, 0)
			s.draw(main, introText, 52, 180)
			if intro.RecycledController().Cursor() == 0 {
				inIntro = false
				s.music("bankok-knights-1.ym", true)
			}
			s.transform(s.Canvas, main, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
			iteration++
			return
		}
		front.Clear()
		back.Clear()
		s.advanceScroll(topFront)
		s.advanceScroll(topBack)
		s.drawScroll(topFront, front, 0, 0)
		s.drawScroll(topBack, back, 0, 0)
		if err := ribbon.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		ribbon.DrawAt(main, 52, 0)
		sineText.Clear()
		sineWave.Clear()
		s.advanceScroll(sine)
		s.drawScroll(sine, sineText, 0, 0)
		wave.DrawAt(sineWave, sineText, 0, 50)
		wave.Advance()
		s.draw(main, sineWave, 52, 72)
		if err := cloud.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		cloud.Draw(main)
		rowFrame := kit.Frame{Time: float64(iteration)}
		if err := outerRows.Update(rowFrame); err != nil {
			s.err = err
			return
		}
		outerRows.Draw(main)
		if err := centerRows.Update(rowFrame); err != nil {
			s.err = err
			return
		}
		centerRows.Draw(main)
		if err := fBack.Update(kit.Frame{Tick: uint64(iteration)}); err != nil {
			s.err = err
			return
		}
		fBack.Draw(main)
		if err := fFront.Update(kit.Frame{Tick: uint64(iteration)}); err != nil {
			s.err = err
			return
		}
		fFront.Draw(main)
		s.transform(s.Canvas, main, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		iteration++
	}
}
