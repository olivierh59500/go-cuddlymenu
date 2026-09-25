package screens

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func (s *Scene) colorshock() {
	background, logo, font, accent := s.asset("fond_tcb.png"), s.asset("logo_tcb.png"), s.asset("font_tcb3.png"), s.asset("scroll_acc.png")
	back, scroll := s.surface(768, 540), s.surface(600, 500)
	r := s.ring(scroll, font, "cuddly-colorshock", s.data.Strings["text"], 6)
	positions := s.data.Numbers["movescroll"]
	if len(positions) < 1873 {
		s.err = fmt.Errorf("missing Colorshock motion table")
		return
	}
	backgroundLayer, err := composite.NewBackground(composite.BackgroundConfig{PeriodY: 1152, Filter: ebiten.FilterLinear})
	if err != nil {
		s.err = err
		return
	}
	vbl := 0.0
	position := 0
	s.render = func() {
		clearBlack(s.Canvas)
		back.Clear()
		scroll.Clear()
		x := 400 - math.Sin(vbl*math.Pi/100)*400
		y := 180 - math.Cos(vbl*math.Pi/200)*300
		hx, hy := float64(background.Bounds().Dx()/2), float64(background.Bounds().Dy()/2)
		backgroundLayer.DrawAt(back, background, x-hx, y-hy)
		s.draw(s.Canvas, back, 0, 0)
		r.Step()
		r.DrawAt(scroll, 0, 0)
		s.draw(scroll, accent, 0, 0)
		if position > 1872 {
			position = 0
		}
		s.draw(s.Canvas, scroll, 90, 110+positions[position])
		vbl += .8
		position += 2
		s.transform(s.Canvas, logo, 370, 120, 1, 1, 0, float64(logo.Bounds().Dx()/2), float64(logo.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		s.Canvas.SubImage(image.Rect(0, 0, 768, 70)).(*ebiten.Image).Fill(color.Black)
	}
}

func (s *Scene) megaScroller() {
	font, tile, bars := s.asset("font.png"), s.asset("bg.png"), s.asset("whitebars.png")
	stage, background, mask, merge := s.surface(640, 400), s.surface(500, 240), s.surface(320, 240), s.surface(320, 240)
	s.filters[stage] = ebiten.FilterNearest
	s.filters[s.Canvas] = ebiten.FilterNearest
	tiles, err := composite.NewBackground(composite.BackgroundConfig{PeriodX: 8, PeriodY: 8, Filter: ebiten.FilterLinear})
	if err != nil {
		s.err = err
		return
	}
	tiles.DrawAt(background, tile, 0, 0)
	for i := 0; i < 48; i++ {
		s.draw(mask, bars, float64(i*8), 0)
	}
	r := s.ring(merge, font, "cuddly-megascroller", s.data.Strings["text"], 10)
	wave := composite.WaveStrips{Axis: composite.Rows, Thickness: 1, Filter: ebiten.FilterNearest, Waves: []composite.StripWave{{Amplitude: 30, Spatial: .03, Speed: -.05}, {Amplitude: 30, Spatial: .01, Speed: .08}}}
	y, dy := 45.0, -2.0
	s.render = func() {
		clearBlack(s.Canvas)
		clearBlack(stage)
		merge.Clear()
		wave.DrawAt(stage, background, -110, -9)
		wave.Advance()
		r.Step()
		r.DrawAt(merge, 0, 0)
		y += dy
		if y >= 20 {
			dy = -2
		}
		if y <= -70 {
			dy = 2
		}
		s.transform(merge, mask, 0, 0, 1, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
		s.draw(stage, merge, 0, y)
		s.part(s.Canvas, stage, composite.Region{Width: 320, Height: 200}, 64, 70, 2, 2)
	}
}

func (s *Scene) bigSprite() {
	front, back, edge, fontIn, fontOut, raster := s.asset("tcb_f.png"), s.asset("tcb_r.png"), s.asset("edge.png"), s.asset("font_in.png"), s.asset("font_out.png"), s.asset("raster.png")
	stage, a, b, starCanvas := s.surface(640, 400), s.surface(568, 41), s.surface(568, 41), s.surface(320, 200)
	s.filters[s.Canvas] = ebiten.FilterNearest
	s.filters[starCanvas] = ebiten.FilterNearest
	r1, r2 := s.ring(a, fontIn, "cuddly-bigsprite", s.data.Strings["text"], 8), s.ring(b, fontOut, "cuddly-bigsprite", s.data.Strings["text"], 8)
	rasterFill, err := composite.NewRasterOverlay(composite.RasterOverlayConfig{
		Image: raster, ScaleX: 85, ScaleY: 1, Alpha: 1,
		VelocityY: -2, WrapY: &composite.RasterWrap{Boundary: -177, Restart: 0, Inclusive: true},
		Filter: s.filters[a], Blend: ebiten.BlendSourceAtop,
	})
	if err != nil {
		s.err = err
		return
	}
	fieldConfig, err := sprites.StreakField(sprites.StreakConfig{Width: 320, Height: 200, Count: 80, Speed: 4, Focal: 100, CenterX: 160, CenterY: 100, Color: color.RGBA{170, 170, 170, 255}, Random: s.rnd})
	if err != nil {
		s.err = err
		return
	}
	field, err := sprites.NewProjectedField(fieldConfig)
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, field.Close)
	// The word is stored in reverse order by the reference choreography.
	var letters []*ebiten.Image
	for _, name := range []string{"s", "r", "a", "e", "b", "e", "r", "a", "c", "e", "h", "t"} {
		letters = append(letters, s.asset(name+".png"))
	}
	// Start the movement cycle at the authored nine-step phase offset.
	phase, flip, flipStep := 9.0, 1.0, -.02
	orbit := motion.DefaultNestedOrbit(motion.Point{X: 320, Y: 200}, motion.Point{X: 160, Y: 400 / 3.7})
	weave := motion.DefaultWeave(motion.Point{X: 308, Y: 190}, motion.Point{X: 308, Y: 500.0 / 6})
	letterGroup, err := sprites.NewGroup(sprites.GroupConfig{
		Frames: letters, Count: len(letters), FrameStride: 1, Weave: &weave,
		PhaseStep: 1.25, Filter: s.filters[stage],
	})
	if err != nil {
		s.err = err
		return
	}
	s.render = func() {
		clearBlack(s.Canvas)
		stage.Clear()
		a.Clear()
		b.Clear()
		starCanvas.Clear()
		r1.Step()
		r2.Step()
		r1.DrawAt(a, 0, 0)
		r2.DrawAt(b, 0, 0)
		rasterFill.Draw(a)
		rasterFill.Step()
		if err := field.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		field.Draw(starCanvas)
		s.transform(stage, starCanvas, 0, 0, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
		phase += .008
		position := orbit.At(phase)
		x, y := position.X, position.Y
		img, angle := front, 0.0
		if flip <= .01 {
			img, angle = back, 180
		}
		s.transform(stage, img, x, y, 1, flip, angle, float64(img.Bounds().Dx()/2), float64(img.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		flip += flipStep
		if flip <= -1 {
			flipStep = .02
		}
		if flip >= 1 {
			flipStep = -.02
		}
		letterGroup.Draw(stage)
		if err := letterGroup.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		s.draw(s.Canvas, edge, 84, 481)
		s.transform(s.Canvas, edge, 706, 481, -1, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.draw(s.Canvas, a, 110, 479)
		s.draw(s.Canvas, b, 110, 479)
		s.draw(s.Canvas, stage, 64, 70)
	}
}

func (s *Scene) fullscreen() {
	background, logo, raster, font := s.asset("backdrop.png"), s.asset("tcblogo.png"), s.asset("colourbar.png"), s.asset("font-fullscreen.png")
	off, scroll := s.surface(768, 52), s.surface(768, 536)
	s.filters[s.Canvas] = ebiten.FilterNearest
	decor, err := composite.NewBackground(composite.BackgroundConfig{Source: image.Rect(0, 0, 16, background.Bounds().Dy()), PeriodX: 16, Filter: ebiten.FilterNearest})
	if err != nil {
		s.err = err
		return
	}
	backgroundScroll := &composite.BackgroundLayer{
		Renderer:  decor,
		Image:     background,
		Pose:      composite.BackgroundPose{Y: 52},
		VelocityX: -4 * TicksPerSecond,
	}
	var rs []*scrolling.Ring
	for i := 0; i < 7; i++ {
		rs = append(rs, s.ring(scroll, font, "cuddly-fullscreen", s.data.Strings[fmt.Sprintf("text%d", i)], 6))
	}
	var letters []*ebiten.Image
	for _, name := range []string{"N", "O", "I", "N", "U", "E", "H", "T"} {
		letters = append(letters, s.asset(name+".png"))
	}
	ys := []float64{-28, 52, 132, 212, 292, 372, 452}
	loop := 0
	weave := motion.DefaultWeave(motion.Point{X: 380, Y: 277}, motion.Point{X: 380, Y: 125.5})
	letterGroup, err := sprites.NewGroup(sprites.GroupConfig{
		Frames: letters, Count: len(letters), FrameStride: 1, Weave: &weave,
		PhaseStep: 1, Filter: s.filters[s.Canvas],
	})
	if err != nil {
		s.err = err
		return
	}
	s.render = func() {
		clearBlack(s.Canvas)
		scroll.Clear()
		off.Clear()
		s.err = backgroundScroll.Update(kit.Frame{Tick: s.frame, Time: float64(s.frame) / TicksPerSecond})
		if s.err != nil {
			return
		}
		backgroundScroll.Draw(s.Canvas)
		for i, r := range rs {
			r.Step()
			r.DrawAt(scroll, 0, ys[i])
		}
		s.draw(s.Canvas, scroll, 0, 0)
		letterGroup.Draw(s.Canvas)
		if err := letterGroup.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		s.Canvas.SubImage(image.Rect(0, 0, 768, 52)).(*ebiten.Image).Fill(color.Black)
		s.transform(off, logo, 95, 5, 1.3, 1.3, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.transform(off, raster, 0, 0, 1, 1.2, 0, 0, 0, 1, ebiten.BlendSourceAtop)
		s.draw(s.Canvas, off, 0, 0)
		if loop >= 128 {
			for i := range ys {
				ys[i] += 2
				if ys[i] >= 540 {
					ys[i] = -28
				}
			}
			loop = 0
		}
		loop++
	}
}

func (s *Scene) knucklebuster() {
	logo, body, head, bass, left, right, font := s.asset("logo.png"), s.asset("batteur.png"), s.asset("tete2.png"), s.asset("bdrum.png"), s.asset("ldrum.png"), s.asset("rdrum.png"), s.asset("fonts.png")
	stage, scroll := s.surface(640, 400), s.surface(640, 34)
	r := s.ring(scroll, font, "cuddly-knucklebuster", s.data.Strings["text"], 8)
	timer := 0
	var on [3]bool
	var triggers [3]float64
	s.render = func() {
		if timer == 0 {
			for i := range on {
				triggers[i] = math.Floor(s.rnd()*750) + 1
			}
			timer = 5
		}
		// The original trigger values reassert the overlay each tick; a new
		// low value does not clear an already latched hit until the timer does.
		for i, v := range triggers {
			if v >= 725 {
				on[i] = true
			}
		}
		clearBlack(stage)
		s.draw(stage, logo, 0, 0)
		s.draw(stage, body, 0, 110)
		scroll.Clear()
		r.Step()
		r.DrawAt(scroll, 0, 0)
		s.draw(stage, scroll, 0, 364)
		if on[0] {
			s.draw(stage, head, 316, 110)
			s.draw(stage, bass, 298, 210)
		}
		if on[1] {
			s.draw(stage, left, 196, 138)
		}
		if on[2] {
			s.draw(stage, right, 381, 131)
		}
		clearBlack(s.Canvas)
		s.draw(s.Canvas, stage, 64, 70)
		timer--
		if timer == 1 {
			on = [3]bool{}
		}
	}
}
