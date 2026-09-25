package menu

import (
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/scrolling"
	"github.com/olivierh59500/democonstructionkit/sprites"

	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Animation = sprites.FrameSequence
type TileSet = sprites.Atlas

func NewTileSet(img *ebiten.Image, tileW, tileH int) *TileSet {
	if tileW <= 0 {
		tileW = 1
	}
	if tileH <= 0 {
		tileH = 1
	}
	atlas, err := sprites.NewAtlas(sprites.AtlasConfig{
		Image: img, TileW: tileW, TileH: tileH, FallbackWhole: true,
	})
	if err != nil {
		panic(err)
	}
	return atlas
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
