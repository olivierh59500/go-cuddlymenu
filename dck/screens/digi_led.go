package screens

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func (s *Scene) digi() {
	union, logo, font := s.asset("unionlogo.png"), s.asset("logo.png"), s.asset("syncfont.png")
	stage := s.surface(640, 400)
	s.filters[stage] = ebiten.FilterNearest
	s.filters[s.Canvas] = ebiten.FilterNearest
	r := s.ring(stage, font, "cuddly-digi", s.data.Strings["text"], 8)
	var letters []*ebiten.Image
	for _, ch := range "HAEY" {
		letters = append(letters, s.asset(string(ch)+".png"))
	}
	logoConfig, err := presets.CuddlyDigiLogo(logo)
	if err != nil {
		s.err = err
		return
	}
	logoRows, err := composite.NewTableWarpLogo(logoConfig)
	if err != nil {
		s.err = err
		return
	}
	weave := motion.DefaultWeave(motion.Point{X: 306, Y: 207}, motion.Point{X: 306, Y: 90.5})
	group, err := sprites.NewGroup(sprites.GroupConfig{
		Frames: letters, Count: len(letters), FrameStride: 1, Weave: &weave,
		Phase: -1, PhaseStep: 1, Filter: s.filters[stage],
	})
	if err != nil {
		s.err = err
		return
	}
	scrollBounce, err := motion.NewWaveClock(presets.RectifiedSine(340, -40, .06))
	if err != nil {
		s.err = err
		return
	}
	s.render = func() {
		clearBlack(s.Canvas)
		clearBlack(stage)
		if err := logoRows.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		logoRows.Draw(stage)
		s.transform(stage, union, 320, logoRows.Bounce()+14, 1, 1, 0, float64(union.Bounds().Dx()/2), float64(union.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		r.Step()
		r.DrawAt(stage, 0, scrollBounce.At(0))
		scrollBounce.Step()
		s.err = group.Update(kit.Frame{Tick: s.frame, Time: float64(s.frame) / TicksPerSecond})
		group.Draw(stage)
		s.draw(s.Canvas, stage, 64, 70)
	}
}

func (s *Scene) led() {
	font, tile, the, bubble := s.asset("ledfont.png"), s.asset("tcbtile.png"), s.asset("the.png"), s.asset("bubble.png")
	stage, led, warped := s.surface(640, 400), s.surface(640, 112), s.surface(384, 233)
	s.filters[s.Canvas] = ebiten.FilterNearest
	s.filters[stage] = ebiten.FilterNearest
	backdrop, err := composite.NewTiledWaveBackdrop(presets.CuddlyLEDBackdrop(tile))
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, backdrop.Close)
	bubbles := s.surface(640, 112)
	bubbleTiles, err := composite.NewBackground(composite.BackgroundConfig{PeriodX: 16, PeriodY: 16, Filter: ebiten.FilterLinear})
	if err != nil {
		s.err = err
		return
	}
	// The authored bubble grid ends two pixels before the text surface edge.
	bubbleTiles.DrawAt(bubbles.SubImage(image.Rect(0, 0, 638, 112)).(*ebiten.Image), bubble, -2, 0)
	gradient, err := composite.NewUniformGradientImage(presets.CuddlyLEDGradient())
	if err != nil {
		s.err = err
		return
	}
	s.surfaces = append(s.surfaces, gradient)
	r := s.ring(led, font, "cuddly-led", s.data.Strings["text"], 16)
	var letters []*ebiten.Image
	for _, name := range []string{"c", "a", "r", "e", "b", "e", "a", "r", "s"} {
		letters = append(letters, s.asset(name+".png"))
	}
	gradientMotion, err := motion.NewWrapBank(motion.WrapBankConfig{
		Start: []float64{0}, Velocity: []float64{-1},
		Lower: &motion.WrapLimit{Boundary: -1400, Restart: 0, Inclusive: true},
	})
	if err != nil {
		s.err = err
		return
	}
	formation, err := motion.NewHarmonicFormation(presets.CuddlyLEDLetterFormationConfig())
	if err != nil {
		s.err = err
		return
	}
	letterGroup, err := sprites.NewGroup(sprites.GroupConfig{
		Frames: letters, Count: len(letters), FrameStride: 1, Harmonic: formation,
		HarmonicClockStart: [2]float64{0, 1.2}, HarmonicClockStep: [2]float64{0, 1.2},
		HarmonicEnvelope: &motion.BounceBankConfig{
			Start: []float64{75}, Velocity: []float64{-1}, Min: 20, Max: 75,
			Inclusive: true, Clamp: true,
		},
	})
	if err != nil {
		s.err = err
		return
	}
	theBounce, err := motion.NewWaveClock(presets.RectifiedSine(95, -80, .045))
	if err != nil {
		s.err = err
		return
	}
	ledBounce, err := motion.NewWaveClock(presets.RectifiedSine(340, -60, .08))
	if err != nil {
		s.err = err
		return
	}
	load := func() {
		warped.Clear()
		backdrop.Draw(warped)
		backdrop.Step()
		s.transform(warped, the, 160, theBounce.At(0), .5, .5, 0, float64(the.Bounds().Dx()/2), float64(the.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		theBounce.Step()
	}
	load() // The reference preloads one background frame before its intro wait.
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		load()
		s.transform(warped, gradient, 0, gradientMotion.At(0), 1, 1, 0, 0, 0, .8, ebiten.BlendSourceOver)
		gradientMotion.Step()
		s.transform(stage, warped, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		letterGroup.Draw(stage)
		if err := letterGroup.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		s.draw(s.Canvas, stage, 64, 70)
		led.Clear()
		s.draw(led, bubbles, 0, 0)
		r.Step()
		r.DrawAt(led, 0, 2)
		s.draw(s.Canvas, led, 64, ledBounce.At(0))
		ledBounce.Step()
	}
}
