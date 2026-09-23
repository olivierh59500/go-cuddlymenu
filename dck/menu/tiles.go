package menu

import (
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Animation struct {
	Duration float64
	Indices  []int
	Loop     bool
}

func (a Animation) Current(t float64) int {
	if len(a.Indices) == 0 {
		return 0
	}
	if a.Duration <= 0 {
		return a.Indices[0]
	}
	ct := math.Max(0, t)
	if !a.Loop && ct >= a.Duration {
		return a.Indices[len(a.Indices)-1]
	}
	if a.Loop {
		ct = math.Mod(ct, a.Duration)
	}
	cp := math.Min(ct/a.Duration, 1)
	frame := int(math.Floor(float64(len(a.Indices)) * cp))
	if frame >= len(a.Indices) {
		frame = len(a.Indices) - 1
	}
	return a.Indices[frame]
}

type TileSet struct {
	Image   *ebiten.Image
	Tiles   []*ebiten.Image
	TileW   int
	TileH   int
	Columns int
}

func NewTileSet(img *ebiten.Image, tileW, tileH int) *TileSet {
	if tileW <= 0 {
		tileW = 1
	}
	if tileH <= 0 {
		tileH = 1
	}
	bounds := img.Bounds()
	columns := bounds.Dx() / tileW
	rows := bounds.Dy() / tileH
	if columns < 1 {
		columns = 1
	}
	if rows < 1 {
		rows = 1
	}
	tiles, err := scrolling.GridImages(img, image.Pt(tileW, tileH), columns, columns*rows)
	if err != nil {
		tiles = []*ebiten.Image{img}
		columns = 1
	}

	return &TileSet{
		Image:   img,
		Tiles:   tiles,
		TileW:   tileW,
		TileH:   tileH,
		Columns: columns,
	}
}

func (t *TileSet) Tile(index int) *ebiten.Image {
	if len(t.Tiles) == 0 {
		return t.Image
	}
	if index < 0 {
		index = 0
	}
	if index >= len(t.Tiles) {
		index = index % len(t.Tiles)
	}
	return t.Tiles[index]
}

type TileMap struct {
	scrollRenderer *scrolling.Scrolling
	Data           [][]int
	Tiles          *TileSet
	WidthPx        int
	HeightPx       int
}

func NewTileMap(data [][]int, tiles *TileSet) *TileMap {
	w := 0
	if len(data) > 0 {
		w = len(data[0]) * tiles.TileW
	}
	h := len(data) * tiles.TileH
	return &TileMap{
		Data:     data,
		Tiles:    tiles,
		WidthPx:  w,
		HeightPx: h,
	}
}

func (m *TileMap) Draw(dst *ebiten.Image, offsetX, offsetY, dstX, dstY, viewW, viewH int) {
	if len(m.Data) == 0 || m.Tiles == nil {
		return
	}
	offsetX, offsetY = m.clip(offsetX, offsetY, viewW, viewH)
	if len(m.Data) == 1 {
		if m.scrollRenderer == nil {
			images := make([]*ebiten.Image, len(m.Data[0]))
			for i, index := range m.Data[0] {
				images[i] = m.Tiles.Tile(index)
			}
			var err error
			m.scrollRenderer, err = scrolling.FromImages(images, float64(m.Tiles.TileW))
			if err != nil {
				panic(err)
			}
		}
		state := scrolling.IdentityState()
		state.X = float64(dstX - offsetX)
		state.Y = float64(dstY - offsetY)
		state.First = offsetX / m.Tiles.TileW
		state.End = state.First + viewW/m.Tiles.TileW + 2
		m.scrollRenderer.DrawAt(dst, state)
		return
	}
	composite.Grid{Columns: len(m.Data[0]), Rows: len(m.Data), Cell: image.Pt(m.Tiles.TileW, m.Tiles.TileH), Overscan: 2, Tile: func(x, y int) *ebiten.Image { return m.Tiles.Tile(m.Data[y][x]) }}.Draw(dst, image.Pt(offsetX, offsetY), image.Pt(dstX, dstY), image.Pt(viewW, viewH))
}

func (m *TileMap) clip(offsetX, offsetY, viewW, viewH int) (int, int) {
	maxX := m.WidthPx - viewW
	if maxX < 0 {
		maxX = 0
	}
	maxY := m.HeightPx - viewH
	if maxY < 0 {
		maxY = 0
	}
	if offsetX < 0 {
		offsetX = 0
	} else if offsetX > maxX {
		offsetX = maxX
	}
	if offsetY < 0 {
		offsetY = 0
	} else if offsetY > maxY {
		offsetY = maxY
	}
	return offsetX, offsetY
}
