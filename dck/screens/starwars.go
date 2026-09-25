package screens

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func (s *Scene) starwars() {
	background, green, red, raster, font, sprite := s.asset("bg.png"), s.asset("fontg.png"), s.asset("fontr.png"), s.asset("scrollraster.png"), s.asset("swfont.png"), s.asset("theunionsprite.png")
	main, a, b, masked := s.surface(320, 200), s.surface(320, 25), s.surface(320, 25), s.surface(320, 200)
	s.filters[s.Canvas] = ebiten.FilterNearest
	r1, r2 := s.ring(a, green, "cuddly-starwars-scroll", s.data.Strings["stext"], 8), s.ring(b, red, "cuddly-starwars-scroll", s.data.Strings["stext"], 8)
	words := s.data.Lists["sstext"]
	if len(words) < 30 {
		s.err = fmt.Errorf("missing Starwars text rows")
		return
	}
	grid := s.bitmap(font, "cuddly-starwars-crawl", ebiten.FilterLinear)
	var wave []float64
	absoluteWave := motion.Wave{Amplitude: 50, Speed: 1, Rectify: true}
	for _, segment := range []struct {
		n        int
		absolute bool
	}{{100, false}, {100, false}, {50, false}, {60, true}, {30, false}, {50, false}, {30, false}, {72, false}, {60, false}, {50, false}, {30, false}, {100, false}} {
		phase := 0.0
		step := 2 * math.Pi / float64(segment.n)
		for i := 0; i < segment.n; i++ {
			v := 50*math.Sin(phase)/2 + .5
			if segment.absolute {
				v = absoluteWave.At(0, phase)
			}
			wave = append(wave, roundHalfUp(v))
			phase += step
		}
	}
	// The reference pushes its tail array as one element. It is never sampled,
	// but its presence extends the counter's wrap threshold by exactly one tick.
	wave = append(wave, 0)
	field, err := sprites.NewProjectedField(sprites.ProjectedFieldConfig{
		Field: sprites.FieldConfig{Count: 400, Depth: sprites.DepthWrap, Near: 0, Far: 130,
			Spawn: func(i int, _ bool) sprites.Point {
				return sprites.Point{X: math.Floor(s.rnd()*320) - 160, Y: math.Floor(s.rnd()*200) - 100, Z: float64(i) * (130.0 / 400)}
			}},
		View:             sprites.FieldView{Camera: geometry.Camera{Center: geometry.Vec2{X: 160, Y: 100}, Focal: 128, Near: math.SmallestNonzeroFloat64}},
		RendererCapacity: 400,
	})
	if err != nil {
		s.err = err
		return
	}
	px, py := append([]float64(nil), s.data.Numbers["spriteX"]...), append([]float64(nil), s.data.Numbers["spriteY"]...)
	period := len(px)
	px = append(px, px...)
	py = append(py, py...)
	if period == 0 || len(py) < len(px) {
		s.err = fmt.Errorf("missing Starwars sprite path")
		return
	}
	s.closeEffects = append(s.closeEffects, field.Close)
	fieldStyle := sprites.FieldStyle{Antialias: true, Sample: func(p sprites.FieldSample, a *sprites.FieldAppearance) bool {
		shade := float32(1)
		if p.Z > 130.0/3 {
			shade = float32(170*257) / 65535
		}
		if p.Z > 260.0/3 {
			shade = float32(85*257) / 65535
		}
		a.Tint.Scale(shade, shade, shade, 1)
		return true
	}}
	crawlConfig, err := presets.CuddlyStarwarsCrawl(grid, words)
	if err != nil {
		s.err = err
		return
	}
	crawl, err := scrolling.New(scrolling.Config{Crawl: &crawlConfig})
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, crawl.Close)
	depth, rotation, fade := 0.0, 180.0, 0.0
	rasterFill, err := composite.NewRasterOverlay(composite.RasterOverlayConfig{
		Image: raster, ScaleX: 320, ScaleY: 1, Alpha: 1,
		VelocityY: -.5, WrapY: &composite.RasterWrap{Boundary: -72, Restart: 0},
		Filter: ebiten.FilterNearest, Blend: ebiten.BlendSourceIn,
	})
	if err != nil {
		s.err = err
		return
	}
	flash, waveIndex, spriteIndex := 0, 20, 0
	s.render = func() {
		clearBlack(s.Canvas)
		for _, surface := range []*ebiten.Image{main, a, b, masked} {
			surface.Clear()
		}
		flash++
		if flash > 300 {
			flash = 0
		}
		if flash == 280 {
			fade = 20
		}
		if fade > 0 {
			fade -= .5
		} else {
			fade = 0
		}
		s.transform(main, background, 30, 10, 1, 1, 0, 0, 0, fade*.05, ebiten.BlendSourceOver)
		depth += 1.5
		rotation -= .02
		field.SetView(sprites.FieldView{Camera: geometry.Camera{Center: geometry.Vec2{X: 160, Y: 100}, Focal: 128, Near: math.SmallestNonzeroFloat64}, Offset: geometry.Vec3{Z: -depth}, Angle: rotation})
		if err := field.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		field.DrawStyle(main, fieldStyle)
		if err := crawl.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		crawl.Draw(main)
		rasterFill.Step()
		waveIndex++
		if waveIndex > len(wave)-80 {
			waveIndex = 20
		}
		r2.Step()
		r2.DrawAt(b, 0, 0)
		for i := 0; i < 20; i++ {
			s.part(masked, b, composite.Region{X: float64(i * 16), Width: 16, Height: 26}, float64(i*16), wave[waveIndex+i]+26, 1, 1)
		}
		s.filters[masked] = ebiten.FilterNearest
		rasterFill.Draw(masked)
		s.filters[masked] = ebiten.FilterLinear
		s.draw(main, masked, 0, 0)
		r1.Step()
		r1.DrawAt(a, 0, 0)
		for i := 0; i < 20; i++ {
			s.part(main, a, composite.Region{X: float64(i * 16), Width: 16, Height: 26}, float64(i*16), wave[waveIndex+i]+26, 1, 1)
		}
		spriteIndex++
		if spriteIndex > period {
			spriteIndex = 0
		}
		for i := 0; i < 8; i++ {
			extra := 0
			if i >= 3 {
				extra = 5
			}
			index := spriteIndex + i*5 + extra
			s.part(main, sprite, composite.Region{X: float64(112 - i*16), Width: 16, Height: 10}, px[index]+10*math.Sin(float64(spriteIndex)*.07+float64(i)), py[index]+5*math.Cos(float64(spriteIndex)*.09+float64(i)), 1, 1)
		}
		s.transform(s.Canvas, main, 64, 64, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
	}
}
