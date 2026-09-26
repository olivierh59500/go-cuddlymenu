package screens

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
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
	curve, err := motion.CompileWaveProgram(presets.CuddlyIntroWaveProgram()...)
	if err != nil {
		s.err = err
		return
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
	r := s.ring(scroll, big, "cuddly-megaball", s.data.Strings["text"], 10)
	grid := s.bitmap(font, "cuddly-values", ebiten.FilterNearest)
	params := []int{251, 246, 1, 4, 5, -10, -1, -2, 1}
	selected, blink := 8, 0
	locations := [][2]float64{{272, 53}, {272, 45}, {272, 37}, {176, 53}, {176, 45}, {176, 37}, {80, 53}, {80, 45}, {80, 37}}
	ballGroup, err := sprites.NewGroup(presets.CuddlyMegaballFormation(ball))
	if err != nil {
		s.err = err
		return
	}
	scrollBounce, err := motion.NewWaveClock(presets.RectifiedSine(112, -66, .04))
	if err != nil {
		s.err = err
		return
	}
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
	s.render = func() {
		clearBlack(s.Canvas)
		half.Clear()
		scroll.Clear()
		r.Step()
		r.DrawAt(scroll, 0, scrollBounce.At(0))
		scrollBounce.Step()
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
		orbit := ballGroup.CoupledOrbitController().Orbit()
		orbit.XIncrement, orbit.YIncrement, orbit.ZIncrement, orbit.QIncrement = float64(params[8]), float64(params[7]), float64(params[6]), float64(params[5])
		orbit.XOffset, orbit.YOffset, orbit.ZOffset, orbit.QOffset, orbit.QScale = float64(params[4]), float64(params[3]), float64(params[2]), float64(params[1]), float64(params[0])
		if err := ballGroup.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		ballGroup.Draw(half)
		s.transform(s.Canvas, half, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
	}
}
