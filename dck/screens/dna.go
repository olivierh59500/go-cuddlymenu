package screens

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/geometry"
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
	fFront, fBack := s.feedback(320, 64, 4, 2, 1, 25, s.data.Numbers["dna_pos"]), s.feedback(320, 64, 4, 2, -1, 37, s.data.Numbers["dna_pos"])
	gradientFront, gradientBack, logo := s.asset("gradient_dna_front.png"), s.asset("gradient_dna_back.png"), s.asset("tcb.png")
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
	iteration, curve := 0, 0
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
		position := curve
		for x := 0; x < 320; x += 16 {
			amp, scale, decal := 15.0, 1.5, float64(position)*8*math.Pi/1280
			if position >= 1280 {
				amp, scale, decal = 30, 1, float64(position-1280)*8*math.Pi/2560
			}
			position += 6
			if position > 3840 {
				position -= 3840
			}
			a1 := math.Mod(math.Pi+.4+decal, 2*math.Pi)
			a2 := math.Mod(a1+1.12, 2*math.Pi)
			y1, y2 := roundHalfUp(amp*math.Sin(a1)), roundHalfUp(amp*math.Sin(a2))
			if a1 > 3.6 || a1 < 1.5 {
				from, to := y1, y2
				if a1 > 3.6 && a1 < 4.6 {
					from = -amp
				}
				if a1 > .5 && a1 < 1.5 {
					to = amp
				}
				h := (to - from) / amp / scale
				if h > .075 {
					s.part(main, back, composite.Region{X: float64(x), Width: 16, Height: 25}, float64(x+52), amp+30+to, 1, -h)
				}
			}
			if a1 > .5 && a1 < 4.6 {
				from, to := y2, y1
				if a1 > .5 && a1 < 1.5 {
					to = amp
				}
				if a1 > 3.6 && a1 < 4.6 {
					from = -amp
				}
				h := (to - from) / amp / scale
				if h > .075 {
					s.part(main, front, composite.Region{X: float64(x), Width: 16, Height: 25}, float64(x+52), amp+30+from, 1, h)
				}
			}
		}
		curve += 8
		if curve > 3840 {
			curve -= 3840
		}
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
		t := float64(iteration)
		for _, x := range []float64{60, 268} {
			decal := 0.0
			if x >= 160 {
				decal = 10
			}
			height := .5 + (1+math.Sin(t/15))*.75
			bounce := float64(int(10 * math.Cos(decal+t/15)))
			for j := 0; j < 56; j++ {
				offset := float64(int(4 * math.Sin(decal+(t+float64(j)*height)/10)))
				s.part(main, logo, composite.Region{Y: float64(int(float64(j) * height)), Width: 96, Height: 1}, x+offset, 100+bounce+float64(j), 1, 1)
			}
		}
		for j := 0; j < 56; j++ {
			z := math.Sin(4*(t/67+float64(j)/131))/4 + .75
			y := 20 * math.Sin(5*(t/61+float64(j)*z/127))
			s.part(main, logo, composite.Region{Y: float64(int(float64(55-j) * .5)), Width: 96, Height: 2}, 212-z*48, 120+y, z, z)
		}
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
