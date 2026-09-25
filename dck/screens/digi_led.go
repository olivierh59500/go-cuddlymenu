package screens

import (
	"image"
	"image/color"
	"math"

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
	curve, err := motion.CompileWaveProgram(presets.CuddlyDigiWaveProgram()...)
	if err != nil {
		s.err = err
		return
	}
	counter := 0
	weave := motion.DefaultWeave(motion.Point{X: 306, Y: 207}, motion.Point{X: 306, Y: 90.5})
	group, err := sprites.NewGroup(sprites.GroupConfig{
		Frames: letters, Count: len(letters), FrameStride: 1, Weave: &weave,
		Phase: -1, PhaseStep: 1, Filter: s.filters[stage],
	})
	if err != nil {
		s.err = err
		return
	}
	logoBounce, err := motion.NewWaveClock(presets.RectifiedSine(200, -160, .03))
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
		bounce := logoBounce.At(0)
		logoBounce.Step()
		for i := 0; i < 170; i++ {
			s.part(stage, logo, composite.Region{Y: float64(i), Width: 335, Height: 1}, 320+curve[(counter+i)%len(curve)]-167.5, 160+bounce/2+float64(i)-.5, 1, 1)
		}
		counter++
		s.transform(stage, union, 320, bounce+14, 1, 1, 0, float64(union.Bounds().Dx()/2), float64(union.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
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
	stage, led, off, warped, tiled := s.surface(640, 400), s.surface(640, 112), s.surface(384, 233), s.surface(384, 233), s.surface(384, 233)
	s.filters[s.Canvas] = ebiten.FilterNearest
	s.filters[stage] = ebiten.FilterNearest
	tiles, err := composite.NewBackground(composite.BackgroundConfig{PeriodX: 32, PeriodY: 33, Filter: ebiten.FilterLinear})
	if err != nil {
		s.err = err
		return
	}
	tiles.DrawAt(tiled, tile, 0, 0)
	bubbles := s.surface(640, 112)
	bubbleTiles, err := composite.NewBackground(composite.BackgroundConfig{PeriodX: 16, PeriodY: 16, Filter: ebiten.FilterLinear})
	if err != nil {
		s.err = err
		return
	}
	// The authored bubble grid ends two pixels before the text surface edge.
	bubbleTiles.DrawAt(bubbles.SubImage(image.Rect(0, 0, 638, 112)).(*ebiten.Image), bubble, -2, 0)
	gradient := ledGradient()
	s.surfaces = append(s.surfaces, gradient)
	wave := composite.WaveStrips{Axis: composite.Rows, Thickness: 1, Filter: ebiten.FilterLinear, Waves: []composite.StripWave{{Amplitude: 6, Spatial: .08, Speed: .2}}}
	r := s.ring(led, font, "cuddly-led", s.data.Strings["text"], 16)
	var letters []*ebiten.Image
	for _, name := range []string{"c", "a", "r", "e", "b", "e", "a", "r", "s"} {
		letters = append(letters, s.asset(name+".png"))
	}
	tileMotion, err := motion.NewWrapBank(motion.WrapBankConfig{
		Start: []float64{0}, Velocity: []float64{-2},
		Lower: &motion.WrapLimit{Boundary: -33, Restart: 0, Inclusive: true},
	})
	if err != nil {
		s.err = err
		return
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
		s.draw(off, tiled, 0, tileMotion.At(0))
		tileMotion.Step()
		wave.DrawAt(warped, off, -32, 0)
		wave.Advance()
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

func ledGradient() *ebiten.Image {
	stops := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {0, 255, 0, 255}, {255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {0, 255, 0, 255}}
	pixels := image.NewRGBA(image.Rect(0, 0, 384, 2000))
	for y := 0; y < 2000; y++ {
		p := (float64(y) + .5) / 2000 * 10
		i := min(9, int(p))
		t := p - float64(i)
		a, b := stops[i], stops[i+1]
		c := color.RGBA{uint8(math.Round(float64(a.R)*(1-t) + float64(b.R)*t)), uint8(math.Round(float64(a.G)*(1-t) + float64(b.G)*t)), uint8(math.Round(float64(a.B)*(1-t) + float64(b.B)*t)), 255}
		for x := 0; x < 384; x++ {
			pixels.SetRGBA(x, y, c)
		}
	}
	return ebiten.NewImageFromImage(pixels)
}
