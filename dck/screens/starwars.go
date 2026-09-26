package screens

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/geometry"
	"github.com/olivierh59500/democonstructionkit/modulation"
	"github.com/olivierh59500/democonstructionkit/presets"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"
)

func (s *Scene) starwars() {
	background, green, red, raster, font, sprite := s.asset("bg.png"), s.asset("fontg.png"), s.asset("fontr.png"), s.asset("scrollraster.png"), s.asset("swfont.png"), s.asset("theunionsprite.png")
	main := s.surface(320, 200)
	s.filters[s.Canvas] = ebiten.FilterNearest
	dualConfig, err := presets.CuddlyStarwarsDualScroll(green, red, raster, s.data.Strings["stext"])
	if err != nil {
		s.err = err
		return
	}
	dualScroll, err := scrolling.NewDualProfiledRing(dualConfig)
	if err != nil {
		s.err = err
		return
	}
	s.closeEffects = append(s.closeEffects, dualScroll.Close)
	words := s.data.Lists["sstext"]
	if len(words) < 30 {
		s.err = fmt.Errorf("missing Starwars text rows")
		return
	}
	grid := s.bitmap(font, "cuddly-starwars-crawl", ebiten.FilterLinear)
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
	spriteTrain, err := sprites.NewSampledSpriteTrain(presets.CuddlyStarwarsSpriteTrain(
		sprite, s.data.Numbers["spriteX"], s.data.Numbers["spriteY"]))
	if err != nil {
		s.err = fmt.Errorf("Starwars sprite path: %w", err)
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
	depth, rotation := 0.0, 180.0
	flashClock, err := modulation.NewPeriodicDecay(presets.CuddlyStarwarsFlash())
	if err != nil {
		s.err = err
		return
	}
	s.render = func() {
		clearBlack(s.Canvas)
		main.Clear()
		s.transform(main, background, 30, 10, 1, 1, 0, 0, 0, flashClock.Step(), ebiten.BlendSourceOver)
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
		if err := dualScroll.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		dualScroll.Draw(main)
		if err := spriteTrain.Update(kit.Frame{}); err != nil {
			s.err = err
			return
		}
		spriteTrain.Draw(main)
		s.transform(s.Canvas, main, 64, 64, 2, 2, 0, 0, 0, 1, ebiten.BlendSourceOver)
	}
}
