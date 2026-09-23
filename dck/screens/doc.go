package screens

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/render"
)

func (s *Scene) doc() {
	background, mountains, raster, ball := s.asset("backdrop.png"), s.asset("mountains.png"), s.asset("rasters.png"), s.asset("ball.png")
	intro, outer, outerRows, inner, innerRows := s.surface(640, 50), s.surface(760, 50), s.surface(640, 50), s.surface(760, 50), s.surface(640, 120)
	r1 := s.ring(intro, s.asset("kh6.png"), "cuddly-doc", s.data.Strings["text1"], 10)
	r2 := s.ring(inner, s.asset("font_in.png"), "cuddly-doc", s.data.Strings["text2"], 10)
	r3 := s.ring(outer, s.asset("font_out.png"), "cuddly-doc", s.data.Strings["text2"], 10)
	var shadows []*ebiten.Image
	for _, name := range []string{"shadow1.png", "shadow2.png", "shadow3.png", "shadow4.png"} {
		shadows = append(shadows, s.asset(name))
	}
	stage, board := s.surface(384, 270), s.surface(320, 80)
	s.filters[s.Canvas] = ebiten.FilterNearest
	white := s.surface(1, 1)
	white.Fill(color.White)
	batch := render.NewBatch(2)
	batch.Options.AntiAlias = true
	quad := func(points [4][2]float64, blend ebiten.Blend) {
		batch.Options.Blend = blend
		batch.Begin(board, white)
		var vertices [4]ebiten.Vertex
		for i, p := range points {
			vertices[i] = render.Vertex(p[0], p[1], .5, .5, color.RGBA{136, 0, 136, 255})
		}
		batch.Quad(vertices)
		batch.Flush()
	}
	curve := make([]float64, 389)
	for i := range curve {
		curve[i] = 20*math.Sin(float64(i)*(7.0/180*math.Pi)) + 30*math.Cos(float64(i)*(3.0/180*math.Pi))
	}
	for i := 0; i < 68; i++ {
		curve = append(curve, 30*math.Sin(float64(i)*(8.0/180*math.Pi)))
	}
	for i := 0; i < 189; i++ {
		curve = append(curve, 30*math.Sin(float64(i)*(8.0/180*math.Pi)))
	}
	jump, mainTick, stripTick := false, 0, 0
	phase, xmove, ymove, xm, speed, vbl, vbl2, vbl4 := 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0
	type projected struct{ x, y, z, scale float64 }
	balls, shade := make([]projected, 4), make([]projected, 4)
	s.music("", false)
	s.render = func() {
		clearBlack(s.Canvas)
		if !jump {
			intro.Clear()
			r1.Step()
			r1.DrawAt(intro, 0, 0)
			if r1.NextRune() == '\\' {
				jump = true
			}
			s.draw(s.Canvas, intro, 64, 62)
			return
		}
		if mainTick == 0 {
			s.music("Cuddly - 3D doc.ym", true)
		}
		stage.Clear()
		s.transform(s.Canvas, background, 0, 0, 77, 1, 0, 0, 0, 1, ebiten.BlendSourceOver)
		s.draw(s.Canvas, mountains, 0, 0)
		outer.Clear()
		outerRows.Clear()
		inner.Clear()
		innerRows.Clear()
		r2.Step()
		r3.Step()
		r2.DrawAt(inner, 0, 0)
		r3.DrawAt(outer, 0, 0)
		vertical := 30 + 30*math.Cos(vbl4/20)
		for j := 0; j < 25; j++ {
			source := 64 + curve[(stripTick+j)%len(curve)]
			s.part(outerRows, outer, composite.Region{X: source, Y: float64(j * 2), Width: 640, Height: 2}, 0, float64(j*2), 1, 1)
			s.part(innerRows, inner, composite.Region{X: source, Y: float64(j * 2), Width: 640, Height: 2}, 0, float64(j*2)+vertical, 1, 1)
		}
		s.draw(s.Canvas, outerRows, 64, 62+vertical)
		black := ebiten.DrawImageOptions{Blend: ebiten.BlendSourceAtop}
		black.GeoM.Scale(640, 120)
		black.ColorScale.Scale(0, 0, 0, 1)
		innerRows.DrawImage(white, &black)
		s.transform(innerRows, raster, 0, 20, 75, 1, 0, 0, 0, 1, ebiten.BlendSourceAtop)
		s.draw(s.Canvas, innerRows, 64, 62)
		vbl4 += 1.2
		stripTick++
		board.Clear()
		xmove += xm * speed * .01
		if xmove > 32 {
			xmove -= 32
		}
		if xmove < 0 {
			xmove += 32
		}
		for i := 0; i < 11; i++ {
			x := float64(i)
			quad([4][2]float64{{-8 + x*32 + xmove, 0}, {8 + x*32 + xmove, 0}, {-752 + x*192 + xmove*6, 80}, {-848 + x*192 + xmove*6, 80}}, ebiten.BlendSourceOver)
		}
		ymove += 315 * speed * .032
		if ymove > 64 {
			ymove -= 64
		}
		if ymove < 0 {
			ymove += 64
		}
		for i := -2; i < 8; i++ {
			y1 := -20 + 250/(250+float64(2*i*32)-ymove)*50
			y2 := -20 + 250/(250+float64(2*i*32+32)-ymove)*50
			quad([4][2]float64{{0, y1}, {320, y1}, {320, y2}, {0, y2}}, ebiten.BlendXor)
		}
		s.transform(stage, board, 32, 149, 1, 1, 0, 0, 0, .3, ebiten.BlendSourceOver)
		speed = -math.Cos(vbl / 40)
		vbl += .16
		xm = 128 * math.Cos(vbl2/40)
		vbl2 += .8
		s.transform(s.Canvas, stage, 0, -128, 2, 2.6, 0, 0, 0, 1, ebiten.BlendSourceOver)
		// Preserve the reference's mixed seconds/frame arithmetic. Playback of
		// these logical ticks uses the application's selected fixed rate.
		t := float64(mainTick) / 60
		segment := int(math.Floor(math.Mod(t/7, 7)))
		alpha := math.Min(1, math.Mod(t/7, 1)*7*1.3)
		for i := 0; i < 4; i++ {
			a, b := docMovement(segment, t, float64(i)), docMovement(segment+1, t, float64(i))
			var animation [4]float64
			for j := range a {
				animation[j] = a[j]*(1-alpha) + b[j]*alpha
			}
			spin, y, displace, radius := animation[0], animation[1], animation[2], animation[3]
			angle := math.Pi * 2 / 360 * (displace * float64(i))
			x, z := radius*math.Cos(angle), -radius*math.Sin(angle)
			phase += math.Pi * 2 / 360 * spin * .2
			phase = math.Mod(phase, math.Pi*2)
			x, z = x*math.Cos(phase)+z*math.Sin(phase), z*math.Cos(phase)-x*math.Sin(phase)
			scale := 400 / (400 + z)
			balls[i] = projected{x: 384 + x*scale, y: 310 + y*scale, z: z, scale: scale * .7}
			shade[i] = projected{x: 384 + x*scale, y: 310 + 60*scale, z: z, scale: scale * .7}
		}
		sort.SliceStable(balls, func(i, j int) bool { return balls[i].z > balls[j].z })
		for _, p := range shade {
			which := 3 - max(0, min(3, int(math.Floor((p.scale-.5)*10/2))))
			dy := math.Min(1, math.Max(0, 1-p.scale)) * 26
			s.transform(s.Canvas, shadows[which], p.x-32, p.y-8-dy, p.scale, p.scale, 0, 0, 0, 1, ebiten.BlendSourceOver)
		}
		for _, p := range balls {
			s.transform(s.Canvas, ball, p.x-32, p.y-32, p.scale, p.scale, 0, 0, 0, 1, ebiten.BlendSourceOver)
		}
		mainTick++
	}
}

func docMovement(index int, t, i float64) [4]float64 {
	if index < 2 && t > 21 {
		index = 7
	}
	switch index {
	case 0, 1:
		return [4]float64{-5, 40, 0, 0}
	case 2:
		return [4]float64{-5, -60 - math.Sin(t*7)*95, 35, 150}
	case 3:
		return [4]float64{5, math.Sin((t+i)*.5*13)*90 - 50, 16, 150}
	case 4:
		q := (t + i) * .125 * 13.5
		return [4]float64{5, 80 - math.Abs(math.Sin(q)*8*math.Cos(q)*42) - 50, 20, 150}
	case 5, 6:
		q := (t + i) * .25 * 13.5
		spin := 5.0
		if index == 6 {
			spin = -7
		}
		return [4]float64{spin, math.Sin(q)*8*math.Cos(q)*22 - 50, 20, 150}
	default:
		return [4]float64{-8, 10 - math.Abs(math.Sin((t*.6+i*.05)*1.75)*70)*2.3, 20, 150}
	}
}
