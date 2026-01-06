package main

import (
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
	total := columns * rows
	tiles := make([]*ebiten.Image, total)
	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			rect := image.Rect(x*tileW, y*tileH, (x+1)*tileW, (y+1)*tileH)
			tiles[y*columns+x] = img.SubImage(rect).(*ebiten.Image)
		}
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
	Data     [][]int
	Tiles    *TileSet
	WidthPx  int
	HeightPx int
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
	startX := offsetX / m.Tiles.TileW
	startY := offsetY / m.Tiles.TileH
	offX := offsetX % m.Tiles.TileW
	offY := offsetY % m.Tiles.TileH
	tilesX := viewW/m.Tiles.TileW + 2
	tilesY := viewH/m.Tiles.TileH + 2
	maxY := len(m.Data)
	maxX := len(m.Data[0])

	var op ebiten.DrawImageOptions
	for y := 0; y < tilesY; y++ {
		mapY := startY + y
		if mapY < 0 || mapY >= maxY {
			continue
		}
		for x := 0; x < tilesX; x++ {
			mapX := startX + x
			if mapX < 0 || mapX >= maxX {
				continue
			}
			tile := m.Tiles.Tile(m.Data[mapY][mapX])
			op.GeoM.Reset()
			op.GeoM.Translate(float64(dstX+x*m.Tiles.TileW-offX), float64(dstY+y*m.Tiles.TileH-offY))
			dst.DrawImage(tile, &op)
		}
	}
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
