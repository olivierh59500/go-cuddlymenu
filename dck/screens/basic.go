package screens

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/motion"
	"github.com/olivierh59500/democonstructionkit/presets"
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
	orbit, err := motion.NewFormulaFormation(presets.CuddlyColorshockOrbit())
	if err != nil {
		s.err = err
		return
	}
	tableClock, err := motion.NewWrapBank(presets.CuddlyColorshockTableClock())
	if err != nil {
		s.err = err
		return
	}
	vbl := 0.0
	s.render = func() {
		clearBlack(s.Canvas)
		back.Clear()
		scroll.Clear()
		point := orbit.At(vbl, 0, 0, 0, 1)
		hx, hy := float64(background.Bounds().Dx()/2), float64(background.Bounds().Dy()/2)
		backgroundLayer.DrawAt(back, background, point.X-hx, point.Y-hy)
		s.draw(s.Canvas, back, 0, 0)
		r.Step()
		r.DrawAt(scroll, 0, 0)
		s.draw(scroll, accent, 0, 0)
		s.draw(s.Canvas, scroll, 90, 110+positions[int(tableClock.At(0))])
		vbl += .8
		tableClock.Step()
		s.transform(s.Canvas, logo, 370, 120, 1, 1, 0, float64(logo.Bounds().Dx()/2), float64(logo.Bounds().Dy()/2), 1, ebiten.BlendSourceOver)
		s.Canvas.SubImage(image.Rect(0, 0, 768, 70)).(*ebiten.Image).Fill(color.Black)
	}
}

func (s *Scene) megaScroller() {
	font, tile, bars := s.asset("font.png"), s.asset("bg.png"), s.asset("whitebars.png")
	stage, mask, merge := s.surface(640, 400), s.surface(320, 240), s.surface(320, 240)
	s.filters[stage] = ebiten.FilterNearest
	s.filters[s.Canvas] = ebiten.FilterNearest
	backdrop, err := composite.NewTiledWaveBackdrop(presets.CuddlyMegaScrollerBackdrop(tile))
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, backdrop.Close)
	maskTiles, err := composite.NewBackground(presets.CuddlyMegaScrollerMaskTiles())
	if err != nil {
		s.err = err
		return
	}
	maskTiles.DrawAt(mask, bars, 0, 0)
	if err := maskTiles.Err(); err != nil {
		s.err = err
		return
	}
	maskMaterial, err := composite.NewRasterOverlay(presets.CuddlyMegaScrollerMask(mask))
	if err != nil {
		s.err = err
		return
	}
	r := s.ring(merge, font, "cuddly-megascroller", s.data.Strings["text"], 10)
	verticalMotion, err := motion.NewBounceBank(motion.BounceBankConfig{
		Start: []float64{45}, Velocity: []float64{-2}, Min: -70, Max: 20,
		Inclusive: true, Directional: true, AllowOutsideStart: true,
	})
	if err != nil {
		s.err = err
		return
	}
	s.render = func() {
		clearBlack(s.Canvas)
		clearBlack(stage)
		merge.Clear()
		backdrop.Draw(stage)
		backdrop.Step()
		r.Step()
		r.DrawAt(merge, 0, 0)
		verticalMotion.Step()
		maskMaterial.Draw(merge)
		s.draw(stage, merge, 0, verticalMotion.At(0))
		s.part(s.Canvas, stage, composite.Region{Width: 320, Height: 200}, 64, 70, 2, 2)
	}
}

func (s *Scene) bigSprite() {
	front, back, edge, fontIn, fontOut, raster := s.asset("tcb_f.png"), s.asset("tcb_r.png"), s.asset("edge.png"), s.asset("font_in.png"), s.asset("font_out.png"), s.asset("raster.png")
	stage, a, b, starCanvas := s.surface(640, 400), s.surface(568, 41), s.surface(568, 41), s.surface(320, 200)
	s.filters[s.Canvas] = ebiten.FilterNearest
	s.filters[starCanvas] = ebiten.FilterNearest
	dualScroll, err := scrolling.NewRingLanes(scrolling.RingLanesConfig{
		Rings: []scrolling.RingConfig{
			{Text: s.data.Strings["text"], Font: s.bitmap(fontIn, "cuddly-bigsprite", s.filters[a]), Viewport: float64(a.Bounds().Dx()), Speed: 8, Controls: true},
			{Text: s.data.Strings["text"], Font: s.bitmap(fontOut, "cuddly-bigsprite", s.filters[b]), Viewport: float64(b.Bounds().Dx()), Speed: 8, Controls: true},
		},
		Y: []float64{0, 0},
	})
	if err != nil {
		s.err = err
		return
	}
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
	phase := 9.0
	orbit := motion.DefaultNestedOrbit(motion.Point{X: 320, Y: 200}, motion.Point{X: 160, Y: 400 / 3.7})
	face, err := sprites.NewAxisFlip(sprites.AxisFlipConfig{
		Front: front, Back: back, SwitchAt: .01, BackAngle: 180,
		Motion:     motion.BounceBankConfig{Start: []float64{1}, Velocity: []float64{-.02}, Min: -1, Max: 1, Inclusive: true, Directional: true},
		SnapCenter: true,
		Filter:     s.filters[stage], Blend: ebiten.BlendSourceOver,
	})
	if err != nil {
		s.err = err
		return
	}
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
		dualScroll.Step()
		dualScroll.DrawLaneAt(0, a, 0, 0)
		dualScroll.DrawLaneAt(1, b, 0, 0)
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
		face.DrawAt(stage, position.X, position.Y)
		face.Step()
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
	scroll := s.surface(768, 536)
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
	logoLayer, err := composite.NewSurfaceLayer(presets.CuddlyFullscreenLogoLayer(logo, raster))
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, logoLayer.Close)
	fontGrid := s.bitmap(font, "cuddly-fullscreen", s.filters[scroll])
	rings := make([]scrolling.RingConfig, 7)
	for i := 0; i < 7; i++ {
		rings[i] = scrolling.RingConfig{
			Text: s.data.Strings[fmt.Sprintf("text%d", i)], Font: fontGrid,
			Viewport: float64(scroll.Bounds().Dx()), Speed: 6, Controls: true,
		}
	}
	laneConfig := scrolling.RingLanesConfig{
		Rings: rings, Y: []float64{-28, 52, 132, 212, 292, 372, 452},
		ShiftEvery: 128, ShiftFirst: 130, ShiftVelocity: []float64{2},
		ShiftUpper: &motion.WrapLimit{Boundary: 540, Restart: -28, Inclusive: true},
	}
	lanes, err := scrolling.New(scrolling.Config{RingLanes: &laneConfig})
	if err != nil {
		s.err = err
		return
	}
	var letters []*ebiten.Image
	for _, name := range []string{"N", "O", "I", "N", "U", "E", "H", "T"} {
		letters = append(letters, s.asset(name+".png"))
	}
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
		s.err = backgroundScroll.Update(kit.Frame{Tick: s.frame, Time: float64(s.frame) / TicksPerSecond})
		if s.err != nil {
			return
		}
		backgroundScroll.Draw(s.Canvas)
		if err := lanes.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		lanes.Draw(scroll)
		s.draw(s.Canvas, scroll, 0, 0)
		letterGroup.Draw(s.Canvas)
		if err := letterGroup.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		s.Canvas.SubImage(image.Rect(0, 0, 768, 52)).(*ebiten.Image).Fill(color.Black)
		logoLayer.Draw(s.Canvas)
	}
}

func (s *Scene) knucklebuster() {
	logo, body, head, bass, left, right, font := s.asset("logo.png"), s.asset("batteur.png"), s.asset("tete2.png"), s.asset("bdrum.png"), s.asset("ldrum.png"), s.asset("rdrum.png"), s.asset("fonts.png")
	stage, scroll := s.surface(640, 400), s.surface(640, 34)
	r := s.ring(scroll, font, "cuddly-knucklebuster", s.data.Strings["text"], 8)
	hits, err := sprites.NewLatchedOverlay(presets.CuddlyKnucklebusterHits(head, bass, left, right, s.rnd))
	if err != nil {
		s.err = err
		return
	}
	s.render = func() {
		if err := hits.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		clearBlack(stage)
		s.draw(stage, logo, 0, 0)
		s.draw(stage, body, 0, 110)
		scroll.Clear()
		r.Step()
		r.DrawAt(scroll, 0, 0)
		s.draw(stage, scroll, 0, 364)
		hits.Draw(stage)
		clearBlack(s.Canvas)
		s.draw(s.Canvas, stage, 64, 70)
	}
}
