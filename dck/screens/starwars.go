package screens

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

func (s *Scene) starwars() {
	background, green, red, raster, font, sprite := s.asset("bg.png"), s.asset("fontg.png"), s.asset("fontr.png"), s.asset("scrollraster.png"), s.asset("swfont.png"), s.asset("theunionsprite.png")
	main, a, b, masked, text, projected := s.surface(320, 200), s.surface(320, 25), s.surface(320, 25), s.surface(320, 200), s.surface(320, 420), s.surface(320, 400)
	s.filters[s.Canvas] = ebiten.FilterNearest
	r1, r2 := s.ring(a, green, 32, 26, 32, s.data.Strings["stext"], 8), s.ring(b, red, 32, 26, 32, s.data.Strings["stext"], 8)
	words := s.data.Lists["sstext"]
	if len(words) < 30 {
		s.err = fmt.Errorf("missing Starwars text rows")
		return
	}
	grid := scrolling.BitmapGrid{Image: font, Width: 17, Height: 11, Columns: font.Bounds().Dx() / 17, ColumnSpan: float64(font.Bounds().Dx()) / 17, First: 32, Filter: ebiten.FilterLinear}
	var wave []float64
	for _, segment := range []struct {
		n        int
		absolute bool
	}{{100, false}, {100, false}, {50, false}, {60, true}, {30, false}, {50, false}, {30, false}, {72, false}, {60, false}, {50, false}, {30, false}, {100, false}} {
		phase := 0.0
		step := 2 * math.Pi / float64(segment.n)
		for i := 0; i < segment.n; i++ {
			v := 50*math.Sin(phase)/2 + .5
			if segment.absolute {
				v = 50 * math.Abs(math.Sin(phase))
			}
			wave = append(wave, roundHalfUp(v))
			phase += step
		}
	}
	// The reference pushes its tail array as one element. It is never sampled,
	// but its presence extends the counter's wrap threshold by exactly one tick.
	wave = append(wave, 0)
	type point struct{ x, y, z float64 }
	stars := make([]point, 400)
	for i := range stars {
		stars[i] = point{math.Floor(s.rnd() * 320), math.Floor(s.rnd() * 200), float64(i) * (130.0 / 400)}
	}
	px, py := append([]float64(nil), s.data.Numbers["spriteX"]...), append([]float64(nil), s.data.Numbers["spriteY"]...)
	period := len(px)
	px = append(px, px...)
	py = append(py, py...)
	if period == 0 || len(py) < len(px) {
		s.err = fmt.Errorf("missing Starwars sprite path")
		return
	}
	white := s.surface(1, 1)
	white.Fill(color.White)
	batch := render.NewBatch(800)
	batch.Options.AntiAlias = true
	type row struct{ y, scale float64 }
	rows := make([]row, 360)
	for i := -160; i < 200; i++ {
		scale := 250 / (250 + float64(i))
		rows[i+160] = row{roundHalfUp((float64(i) + 350) * scale), scale}
	}
	depth, rotation, fade, rasterY := 0.0, 180.0, 0.0, 0.0
	flash, waveIndex, word, scrollY, spriteIndex := 0, 20, 0, 0, 0
	s.render = func() {
		clearBlack(s.Canvas)
		for _, surface := range []*ebiten.Image{main, a, b, masked, text, projected} {
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
		sin, cos := math.Sincos(rotation)
		batch.Begin(main, white)
		for _, p := range stars {
			z := p.z - depth
			if z > 130 || z < 0 {
				z -= 130 * math.Floor(z/130)
			}
			if z == 0 {
				continue
			}
			x := (p.x-160)*cos - (p.y-100)*sin
			y := (p.x-160)*sin + (p.y-100)*cos
			shade := uint8(255)
			if z > 130.0/3 {
				shade = 170
			}
			if z > 260.0/3 {
				shade = 85
			}
			scale := 128 / z
			batch.Rect(x*scale+160, y*scale+100, 1, 1, white.Bounds(), color.RGBA{shade, shade, shade, 255})
		}
		batch.Flush()
		scrollY--
		if scrollY < -16 {
			scrollY = 0
			word++
			if word > len(words)-30 {
				word = 0
			}
		}
		for line := 0; line < 30; line++ {
			chars := []rune(words[line+word])
			x := float64(300-len(chars)*16)/2 - 8
			for i, ch := range chars {
				if i >= 12 {
					break
				}
				grid.Print(text, string(ch), x+float64(i*20), float64(line*17+scrollY), 1, 1)
			}
		}
		previousY := 0.0
		for i, p := range rows {
			if p.y != previousY {
				h := 1 - float64(i)/float64(len(rows))
				srcY := float64(310 - i)
				if srcY >= 0 {
					scale := p.scale * .8
					composite.Row{Source: composite.Region{Y: srcY, Width: 320, Height: h}, X: 160 - 160*scale, Y: p.y - h/2, Width: 320 * scale, Height: h, Filter: ebiten.FilterLinear}.Draw(projected, text)
				}
			}
			previousY = p.y
		}
		s.part(main, projected, composite.Region{Y: 300, Width: 320, Height: 100}, 0, 100, 1, 1)
		rasterY -= .5
		if rasterY < -72 {
			rasterY = 0
		}
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
		s.transform(masked, raster, 0, rasterY, 320, 1, 0, 0, 0, 1, ebiten.BlendSourceIn)
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
