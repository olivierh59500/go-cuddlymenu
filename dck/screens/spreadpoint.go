package screens

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"
	"github.com/olivierh59500/democonstructionkit/timeline"
)

func (s *Scene) spreadpoint() {
	var cards []*ebiten.Image
	for i := 1; i <= 4; i++ {
		cards = append(cards, s.asset(fmt.Sprintf("intro%d.png", i)))
	}
	in, out, raster, font, gradient, dnaFont, dnaGradient, ball := s.asset("tcb-in.png"), s.asset("tcb-out.png"), s.asset("tcb-raster.png"), s.asset("font.png"), s.asset("gradient.png"), s.asset("font_dna.png"), s.asset("gradient_dna.png"), s.asset("ball.png")
	main, spread := s.surface(416, 276), s.surface(320, 200)
	s.filters[s.Canvas] = ebiten.FilterNearest
	// Stretch single-pixel ramps without sampling transparent side padding.
	s.filters[spread] = ebiten.FilterNearest
	logoLayer, err := composite.NewSurfaceLayer(presets.CuddlySpreadpointLogoLayer(in, raster, out))
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, logoLayer.Close)
	grid := s.bitmap(font, "cuddly-spreadpoint", ebiten.FilterLinear)
	bandsConfig := presets.CuddlySpreadpointBands(grid, s.data.Strings["text"])
	bands, err := scrolling.New(scrolling.Config{Bands: &bandsConfig})
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, bands.Close)
	feedbackConfig, err := presets.CuddlySpreadpointFeedbackScroll(
		s.bitmap(dnaFont, "cuddly-dna", ebiten.FilterLinear), dnaGradient,
		s.data.Strings["text_dna"], s.data.Numbers["dna_pos"])
	if err != nil {
		s.err = err
		return
	}
	dna, err := scrolling.New(feedbackConfig)
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, dna.Close)
	ballConfig, err := presets.CuddlySpreadpointBallFormation(ball, 20)
	if err != nil {
		s.err = err
		return
	}
	ballTrain, err := sprites.NewGroup(ballConfig)
	if err != nil {
		s.err = err
		return
	}
	phase, err := motion.NewPhaseSequence(presets.CuddlySpreadpointLogoPhases())
	if err != nil {
		s.err = err
		return
	}
	orbit, err := motion.NewHarmonicTransform(presets.CuddlySpreadpointLogoTransform())
	if err != nil {
		s.err = err
		return
	}
	introCards, err := timeline.NewTintedCards(presets.CuddlySpreadpointCards())
	if err != nil {
		s.err = err
		return
	}
	white := s.surface(1, 1)
	white.Fill(color.White)
	introDone, iteration := false, 0
	s.music("intro1.mp3", false)
	s.render = func() {
		if !introDone {
			frame := introCards.Next()
			if frame.Cue >= 0 {
				if frame.Completed {
					s.music("master.mp3", true)
				} else {
					s.music(fmt.Sprintf("intro%d.mp3", frame.Cue+1), false)
				}
			}
			main.Clear()
			if !frame.Completed {
				s.draw(main, cards[frame.Card], 0, 0)
				op := ebiten.DrawImageOptions{Blend: ebiten.BlendSourceAtop}
				op.GeoM.Scale(416, 276)
				op.ColorScale.Scale(float32(frame.Tint.R)/255, float32(frame.Tint.G)/255, float32(frame.Tint.B)/255, float32(frame.Tint.A)/255)
				main.DrawImage(white, &op)
			}
			clearBlack(s.Canvas)
			s.transform(s.Canvas, main, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
			introDone = frame.Completed
			return
		}
		clearBlack(main)
		if iteration >= 1284 {
			spread.Clear()
			if err := bands.Update(kit.Frame{Tick: uint64(iteration - 1284)}); err != nil {
				s.err = err
				return
			}
			bands.Draw(spread)
			s.transform(spread, gradient, 0, 0, 320, 1, 0, 0, 0, 1, ebiten.BlendSourceIn)
			s.draw(main, spread, 52, 29)
		}
		if iteration >= 804 {
			if err := ballTrain.Update(kit.Frame{}); err != nil {
				s.err = err
				return
			}
			ballTrain.Draw(main)
		}
		pose := orbit.At(phase.Current())
		if err := logoLayer.SetPassTransform(0, pose.X, pose.Y, pose.ScaleX, pose.ScaleY, 0); err != nil {
			s.err = err
			return
		}
		if err := logoLayer.SetPassTransform(2, pose.X, pose.Y, pose.ScaleX, pose.ScaleY, 0); err != nil {
			s.err = err
			return
		}
		logoLayer.Draw(main)
		if iteration >= 1952 {
			if err := dna.Update(kit.Frame{Tick: uint64(iteration)}); err != nil {
				s.err = err
				return
			}
			dna.Draw(main)
		}
		s.transform(s.Canvas, main, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		phase.Step()
		iteration++
	}
}
