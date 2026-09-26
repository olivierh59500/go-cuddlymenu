package screens

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func (s *Scene) dna() {
	main := s.surface(416, 276)
	s.filters[s.Canvas] = ebiten.FilterNearest
	introText, front, back, sineText, sineWave, dnaFront, dnaBack := s.surface(320, 25), s.surface(320, 25), s.surface(320, 25), s.surface(320, 16), s.surface(320, 100), s.surface(320, 25), s.surface(320, 25)
	intro := s.ring(introText, s.asset("font_intro.png"), "cuddly-dna", s.data.Strings["text_intro"], 6)
	topFront := s.ring(front, s.asset("font_top_front.png"), "cuddly-dna", s.data.Strings["text_top"], 6)
	topBack := s.ring(back, s.asset("font_top_back.png"), "cuddly-dna", s.data.Strings["text_top"], 6)
	sine := s.ring(sineText, s.asset("font_sin.png"), "cuddly-dna-sine", s.data.Strings["text_sin"], 3)
	rFront := s.ring(dnaFront, s.asset("font_dna_front.png"), "cuddly-dna", s.data.Strings["text_dna"], 4)
	rBack := s.ring(dnaBack, s.asset("font_dna_back.png"), "cuddly-dna", s.data.Strings["text_dna"], 4)
	ribbon, err := composite.NewTwistingRibbon(presets.CuddlyDNARibbon(front, back))
	if err != nil {
		s.err = err
		return
	}
	fFront, fBack := s.feedback(320, 64, 4, 2, 1, 25, s.data.Numbers["dna_pos"]), s.feedback(320, 64, 4, 2, -1, 37, s.data.Numbers["dna_pos"])
	gradientFront, gradientBack, logo := s.asset("gradient_dna_front.png"), s.asset("gradient_dna_back.png"), s.asset("tcb.png")
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
	points := make([]geometry.Vec3, 125)
	for i := range points {
		a, b := math.Pi*s.rnd(), 2*math.Pi*s.rnd()
		points[i] = geometry.Vec3{X: 100 * math.Sin(a) * math.Cos(b), Y: 100 * math.Sin(a) * math.Sin(b), Z: 100 * math.Cos(a)}
	}
	type dot struct{ x, y, z, r float64 }
	dots := make([]dot, len(points))
	discs, err := sprites.NewDiscs(len(points))
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, discs.Close)
	particles := make([]sprites.Disc, len(points))
	var tint ebiten.ColorScale
	tint.ScaleWithColor(color.RGBA{238, 136, 0, 255})
	rotation := 0.0
	iteration := 0
	inIntro := true
	s.music("", false)
	s.render = func() {
		clearBlack(main)
		if inIntro {
			introText.Clear()
			intro.Step()
			intro.DrawAt(introText, 0, 0)
			s.draw(main, introText, 52, 180)
			if intro.Cursor() == 0 {
				inIntro = false
				s.music("bankok-knights-1.ym", true)
			}
			s.transform(s.Canvas, main, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
			iteration++
			return
		}
		front.Clear()
		back.Clear()
		topFront.Step()
		topBack.Step()
		topFront.DrawAt(front, 0, 0)
		topBack.DrawAt(back, 0, 0)
		if err := ribbon.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		ribbon.DrawAt(main, 52, 0)
		sineText.Clear()
		sineWave.Clear()
		sine.Step()
		sine.DrawAt(sineText, 0, 0)
		wave.DrawAt(sineWave, sineText, 0, 50)
		wave.Advance()
		s.draw(main, sineWave, 52, 72)
		rotation += .03
		sin, cos := math.Sincos(rotation)
		focal := 138 / math.Tan(20*math.Pi/180)
		for i, p := range points {
			x, z := p.X*cos+p.Z*sin, p.Z*cos-p.X*sin
			scale := focal / (900 - z)
			dots[i] = dot{x: 208 + x*scale, y: 138 - (p.Y+16)*scale, z: z, r: scale}
		}
		sort.SliceStable(dots, func(i, j int) bool { return dots[i].z < dots[j].z })
		for i, p := range dots {
			particles[i] = sprites.Disc{X: p.x, Y: p.y, Radius: p.r, ColorScale: tint}
		}
		discs.DrawAt(main, particles, 0, 0)
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
		dnaBack.Clear()
		rBack.Step()
		rBack.DrawAt(dnaBack, 0, 0)
		fBack.Step(dnaBack)
		fBack.DrawAt(main, 52, 175, -iteration/2, gradientBack)
		dnaFront.Clear()
		rFront.Step()
		rFront.DrawAt(dnaFront, 0, 0)
		fFront.Step(dnaFront)
		fFront.DrawAt(main, 52, 175, iteration/2, gradientFront)
		s.transform(s.Canvas, main, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		iteration++
	}
}
