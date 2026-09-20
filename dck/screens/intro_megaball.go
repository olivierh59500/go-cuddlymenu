package screens

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func (s *Scene) intro() {
	calvin, main, union, logo := s.asset("calvin.png"), s.asset("main.png"), s.asset("unionlogo.png"), s.asset("tcblogodist.png")
	stars := []*ebiten.Image{s.asset("star1.png"), s.asset("star2.png")}
	a, dist := s.surface(400, 60), s.surface(256, 122)
	s.filters[s.Canvas] = ebiten.FilterNearest
	s.draw(a, logo, 46, 0)
	s.draw(dist, union, 0, 0)
	xwave := composite.WaveStrips{Axis: composite.Rows, Thickness: 1, CenterStrips: true, Filter: ebiten.FilterLinear, Waves: []composite.StripWave{{Amplitude: 7, Spatial: .03, Speed: -.035}, {Amplitude: 7, Spatial: .01, Speed: .05}}}
	ywave := composite.WaveStrips{Axis: composite.Columns, Thickness: 1, CenterStrips: true, PixelSnap: true, Filter: ebiten.FilterNearest, Waves: []composite.StripWave{{Amplitude: 4, Spatial: .02, Speed: -.035}, {Amplitude: 4, Spatial: .005, Speed: .05}}}
	chain, err := composite.NewWaveChain(
		composite.WavePass{Size: image.Pt(460, 120), X: 230, Y: 40, Wave: xwave},
		composite.WavePass{Size: image.Pt(460, 80), X: 0, Y: 35, Wave: ywave},
	)
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, chain.Close)
	curve := make([]float64, 252)
	for i := 0; i < 252; i++ {
		curve = append(curve, 80*math.Sin(float64(i)*.05))
	}
	for repeat := 0; repeat < 3; repeat++ {
		part := make([]float64, 259)
		for i := 0; i < 30; i++ {
			part[i] = 200 * math.Sin(float64(i)*.05)
		}
		for i := 30; i < 80; i++ {
			part[i] = 200
		}
		for i := 30; i < 90; i++ {
			part[50+i] = 200 * math.Sin(float64(i)*.05)
		}
		for i := 88; i < 150; i++ {
			part[50+i] = -200
		}
		for i := 90; i < 149; i++ {
			part[110+i] = 220 * math.Sin(9.6+float64(i)*.02)
		}
		curve = append(curve, part...)
	}
	for i := 0; i < 252; i++ {
		curve = append(curve, 40*math.Sin(float64(i)*.05))
	}
	px, py := s.data.Numbers["starposX"], s.data.Numbers["starposY"]
	profile := composite.ProfileStrips{Offsets: curve, Speed: 2, Thickness: 1, Filter: ebiten.FilterNearest}
	points := make([]sprites.SparklePoint, 9)
	for i := range points {
		points[i] = sprites.SparklePoint{X: px[i], Y: py[i]}
	}
	sparkles, err := sprites.NewSparkles(sprites.SparkleConfig{Images: []sprites.SparkleImage{{Image: stars[0], Spin: 10}, {Image: stars[1], Angles: []float64{0, 45}}}, Positions: points, StartScale: 1, EndScale: 0, ScaleStep: -.025, PauseTicks: 20, Filter: ebiten.FilterNearest})
	if err != nil {
		s.err = err
		return
	}
	timer := 500
	s.music("", false)
	s.render = func() {
		clearBlack(s.Canvas)
		if timer > 0 {
			if timer > 200 {
				s.draw(s.Canvas, calvin, 288, 124)
			}
			timer -= 2
			return
		}
		if timer == 0 {
			s.music("cuddlyintro.mp3", true)
			timer = -1
		}
		s.draw(s.Canvas, main, 0, 0)
		profile.DrawAt(s.Canvas, dist, 256, 204)
		profile.Advance()
		warped := chain.Render(a)
		chain.Advance()
		s.transform(s.Canvas, warped, 384, 130, 2, 2, 0, 230, 40, 1, ebiten.BlendSourceOver)
		sparkles.DrawAt(s.Canvas, 0, 0)
		sparkles.Advance()
	}
}

func (s *Scene) megaball() {
	font, ball, big := s.asset("font1.png"), s.asset("ball.png"), s.asset("bigfont.png")
	half, scroll := s.surface(384, 270), s.surface(640, 400)
	s.filters[half] = ebiten.FilterNearest
	s.filters[s.Canvas] = ebiten.FilterNearest
	r := s.ring(scroll, big, 320, 288, 32, s.data.Strings["text"], 10)
	grid := scrolling.BitmapGrid{Image: font, Width: 8, Height: 8, Columns: font.Bounds().Dx() / 8, ColumnSpan: float64(font.Bounds().Dx()) / 8, First: 32, Filter: ebiten.FilterNearest}
	params := []int{251, 246, 1, 4, 5, -10, -1, -2, 1}
	selected, blink := 8, 0
	locations := [][2]float64{{272, 53}, {272, 45}, {272, 37}, {176, 53}, {176, 45}, {176, 37}, {80, 53}, {80, 45}, {80, 37}}
	orbit := motion.CoupledOrbit{CenterX: 192, CenterY: 135, Radius: 60, DepthRadius: 30, PhaseStep: .00025}
	s.input = func(in Input) {
		if in.Up {
			selected = (selected + 1) % 9
		}
		if in.Down {
			selected = (selected + 8) % 9
		}
		if in.Left {
			params[selected] = max(-254, params[selected]-1)
		}
		if in.Right {
			params[selected] = min(254, params[selected]+1)
		}
	}
	phase := 0.0
	s.render = func() {
		clearBlack(s.Canvas)
		half.Clear()
		scroll.Clear()
		r.Step()
		r.DrawAt(scroll, 0, 112-math.Abs(math.Sin(phase)*66))
		phase += .04
		s.draw(s.Canvas, scroll, 64, 70)
		blink = (blink + 1) % 2
		p := locations[selected]
		s.part(half, font, composite.Region{Y: float64(49 + blink), Width: 24, Height: 1}, p[0], p[1], 1, 1)
		labels := []string{"QMUL", "QOFF", "ZOFF", "YOFF", "XOFF", "QINC", "ZINC", "YINC", "XINC"}
		for i, value := range params {
			loc := locations[i]
			sign := " "
			if value < 0 {
				sign = "-"
			}
			grid.Print(half, fmt.Sprintf("%s %s%03d", labels[i], sign, int(math.Abs(float64(value)))), loc[0]-48, loc[1]-7, 1, 1)
		}
		orbit.XIncrement, orbit.YIncrement, orbit.ZIncrement, orbit.QIncrement = float64(params[8]), float64(params[7]), float64(params[6]), float64(params[5])
		orbit.XOffset, orbit.YOffset, orbit.ZOffset, orbit.QOffset, orbit.QScale = float64(params[4]), float64(params[3]), float64(params[2]), float64(params[1]), float64(params[0])
		for _, start := range []int{0, 40} {
			for i := start; i < start+19; i++ {
				x, y := orbit.Next(float64(i))
				s.transform(half, ball, x, y, 1, 1, 0, float64(ball.Bounds().Dx()/2), float64(ball.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
			}
		}
		s.transform(s.Canvas, half, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
	}
}
