package screens

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/scrolling"
)

func (s *Scene) feedback(width, height, speedX, speedY, direction int, insertY float64, profile []float64) *scrolling.FeedbackDNA {
	rows := make([]int, len(profile))
	for i, v := range profile {
		rows[i] = int(v)
	}
	f, err := scrolling.NewFeedbackDNA(scrolling.FeedbackDNAConfig{Width: width, Height: height, HorizontalSpeed: speedX, VerticalSpeed: speedY, ColumnWidth: 2, Direction: direction, InsertY: insertY, Profile: rows, Filter: ebiten.FilterLinear})
	if err != nil {
		s.err = err
		return nil
	}
	s.closeEffects = append(s.closeEffects, f.Close)
	return f
}
func (s *Scene) spreadpoint() {
	var cards []*ebiten.Image
	for i := 1; i <= 4; i++ {
		cards = append(cards, s.asset(fmt.Sprintf("intro%d.png", i)))
	}
	in, out, raster, font, gradient, dnaFont, dnaGradient, ball := s.asset("tcb-in.png"), s.asset("tcb-out.png"), s.asset("tcb-raster.png"), s.asset("font.png"), s.asset("gradient.png"), s.asset("font_dna.png"), s.asset("gradient_dna.png"), s.asset("ball.png")
	main, logo, spread, textSurface, dnaText := s.surface(416, 276), s.surface(128, 128), s.surface(320, 200), s.surface(len([]rune(s.data.Strings["text"]))*8+320, 6), s.surface(320, 25)
	s.filters[s.Canvas] = ebiten.FilterNearest
	s.filters[spread] = ebiten.FilterNearest
	s.filters[logo] = ebiten.FilterNearest
	grid := scrolling.BitmapGrid{Image: font, Width: 8, Height: 6, Columns: font.Bounds().Dx() / 8, ColumnSpan: float64(font.Bounds().Dx()) / 8, First: 32, Filter: ebiten.FilterLinear}
	grid.Print(textSurface, s.data.Strings["text"], 0, 0, 1, 1)
	head := s.surface(320, 6)
	s.part(head, textSurface, composite.Region{Width: 320, Height: 6}, 0, 0, 1, 1)
	s.draw(textSurface, head, float64(textSurface.Bounds().Dx()-320), 0)
	r := s.ring(dnaText, dnaFont, 32, 25, 32, s.data.Strings["text_dna"], 4)
	dna := s.feedback(320, 64, 4, 2, 1, 4, s.data.Numbers["dna_pos"])
	var angles []float64
	a := math.Pi
	for _, segment := range []struct {
		count int
		step  float64
		reset *float64
	}{{300, 0, number(math.Pi)}, {40, -math.Pi / 40, nil}, {300, 0, number(0)}, {640, -math.Pi / 40, nil}, {480, -math.Pi / 32, nil}, {512, math.Pi / 32, nil}} {
		for i := 0; i < segment.count; i++ {
			if segment.reset != nil {
				a = *segment.reset
			} else {
				a += segment.step
			}
			angles = append(angles, math.Mod(a, 2*math.Pi))
		}
	}
	white := s.surface(1, 1)
	white.Fill(color.White)
	intro, iteration := 0, 0
	s.music("intro1.mp3", false)
	s.render = func() {
		if intro < 4 {
			red, green := 0, 0
			switch {
			case iteration < 20:
				red = int(math.Floor(float64(iteration) * 255 / 19))
				green = red
				iteration++
			case iteration < 160:
				red = 255
				green = int(math.Floor(34 + float64(159-iteration)*221/139))
				iteration++
			case iteration < 180:
				n := iteration - 160
				red = int(math.Floor(float64(19-n) * 255 / 19))
				green = int(math.Floor(float64(19-n) * 34 / 19))
				iteration++
			default:
				iteration = 0
				intro++
				if intro < 4 {
					s.music(fmt.Sprintf("intro%d.mp3", intro+1), false)
				} else {
					s.music("master.mp3", true)
				}
			}
			main.Clear()
			if intro < 4 {
				s.draw(main, cards[intro], 0, 0)
				op := ebiten.DrawImageOptions{Blend: ebiten.BlendSourceAtop}
				op.GeoM.Scale(416, 276)
				op.ColorScale.Scale(float32(red)/255, float32(green)/255, float32(green)/255, 1)
				main.DrawImage(white, &op)
			}
			clearBlack(s.Canvas)
			s.transform(s.Canvas, main, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
			return
		}
		clearBlack(main)
		if iteration >= 1284 {
			spread.Clear()
			speeds := []int{12, 11, 10, 9, 8, 7, 6, 5, 4, 5, 6, 7, 8, 9, 10, 11, 12, 11, 10, 9, 8, 7, 6, 5, 4, 5, 6, 7, 8, 9, 10, 11, 12}
			for row, speed := range speeds {
				travel := speed * (iteration - 1284)
				x, source := 0, 0
				if travel < 320 {
					x = 320 - travel
				} else {
					source = (travel - 320) % (textSurface.Bounds().Dx() - 320)
				}
				s.part(spread, textSurface, composite.Region{X: float64(source), Width: 320, Height: 6}, float64(x), float64(row*6), 1, 1)
			}
			s.transform(spread, gradient, 0, 0, 320, 1, 0, 0, 0, 1, ebiten.BlendSourceIn)
			s.draw(main, spread, 52, 29)
		}
		if iteration >= 804 {
			t := float64(iteration - 804)
			for j := 0; j < 20; j++ {
				p := t/71 + float64(j)/47
				x := 151 + roundHalfUp(151*math.Sin(5*p))
				y := 50 + roundHalfUp(50*math.Sin(8*p))
				s.draw(main, ball, x+52, y+29)
			}
		}
		angle := angles[iteration%len(angles)]
		y, z := math.Sin(angle), math.Cos(angle)/4+.75
		x := 64 - 64*z
		logo.Clear()
		s.transform(logo, in, x, 45+40*y, z, z, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.transform(logo, raster, 0, 0, 128, 1, 0, 0, 0, 1, ebiten.BlendSourceIn)
		s.transform(logo, out, x, 45+40*y, z, z, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.draw(main, logo, 148, 29)
		if iteration >= 1952 {
			dnaText.Clear()
			r.Step()
			r.DrawAt(dnaText, 0, 0)
			dna.Step(dnaText)
			dna.DrawAt(main, 52, 150, 0, dnaGradient)
		}
		s.transform(s.Canvas, main, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		iteration++
	}
}
func number(v float64) *float64 { return &v }
func roundHalfUp(v float64) float64 { return math.Floor(v + .5) }
