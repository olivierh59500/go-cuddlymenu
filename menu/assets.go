package menu

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"path"

	"github.com/hajimehoshi/ebiten/v2"

	gameassets "go-cuddlymenu/assets"
)

type Assets struct {
	Tiles     *ebiten.Image
	Dude      *ebiten.Image
	Carebears *ebiten.Image
	Chrome    *ebiten.Image
	MenuYM    []byte
}

func LoadAssets(maxTileIndex int) *Assets {
	assets := &Assets{}
	const root = "menu"

	tilesPath := path.Join(root, "tiles.png")
	assets.Tiles = loadImage(tilesPath, func() *ebiten.Image {
		return makePlaceholderTiles(tileSize, tileSize, maxTileIndex+1)
	})

	dudePath := path.Join(root, "dude.png")
	assets.Dude = loadImage(dudePath, func() *ebiten.Image {
		return makePlaceholderSheet(640, 128, color.RGBA{220, 80, 80, 255})
	})

	carebearsPath := path.Join(root, "carebears.png")
	assets.Carebears = loadImage(carebearsPath, func() *ebiten.Image {
		return makePlaceholderCarebears()
	})

	chromePath := path.Join(root, "chrome.png")
	assets.Chrome = loadImageWithOptions(chromePath, func() *ebiten.Image {
		return makePlaceholderScrollFont()
	}, &ebiten.NewImageFromImageOptions{Unmanaged: true})

	ymPath := path.Join(root, "menu.ym")
	data, err := gameassets.Files.ReadFile(ymPath)
	if err != nil {
		log.Printf("menu.ym missing (%v): YM playback disabled", err)
	} else {
		assets.MenuYM = data
	}

	return assets
}

func loadImage(path string, fallback func() *ebiten.Image) *ebiten.Image {
	return loadImageWithOptions(path, fallback, nil)
}

func loadImageWithOptions(path string, fallback func() *ebiten.Image, options *ebiten.NewImageFromImageOptions) *ebiten.Image {
	data, err := gameassets.Files.ReadFile(path)
	if err != nil {
		log.Printf("missing asset %s (%v), using placeholder", path, err)
		return fallback()
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("failed to decode %s (%v), using placeholder", path, err)
		return fallback()
	}
	return ebiten.NewImageFromImageWithOptions(img, options)
}

func makePlaceholderTiles(tileW, tileH, total int) *ebiten.Image {
	columns := 16
	rows := int(math.Ceil(float64(total) / float64(columns)))
	w := columns * tileW
	h := rows * tileH
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < total; i++ {
		x := (i % columns) * tileW
		y := (i / columns) * tileH
		col := color.RGBA{uint8(40 + (i*7)%180), uint8(60 + (i*11)%160), uint8(90 + (i*13)%120), 255}
		draw.Draw(img, image.Rect(x, y, x+tileW, y+tileH), &image.Uniform{C: col}, image.Point{}, draw.Src)
	}
	return ebiten.NewImageFromImage(img)
}

func makePlaceholderSheet(width, height int, base color.Color) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: base}, image.Point{}, draw.Src)
	return ebiten.NewImageFromImage(img)
}

func makePlaceholderCarebears() *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 384, 20))
	for i := 0; i < 12; i++ {
		x := i * carebearTileW
		col := color.RGBA{uint8(80 + i*12), uint8(160 - i*5), uint8(120 + i*8), 255}
		draw.Draw(img, image.Rect(x, 0, x+carebearTileW, carebearTileH), &image.Uniform{C: col}, image.Point{}, draw.Src)
	}
	return ebiten.NewImageFromImage(img)
}

func makePlaceholderScrollFont() *ebiten.Image {
	totalTiles := len(scrollerCharWidth) * 3
	columns := 16
	rows := int(math.Ceil(float64(totalTiles) / float64(columns)))
	w := columns * scrollTileW
	h := rows * scrollTileH
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < totalTiles; i++ {
		x := (i % columns) * scrollTileW
		y := (i / columns) * scrollTileH
		col := color.RGBA{uint8(30 + (i*9)%180), uint8(30 + (i*5)%140), uint8(30 + (i*7)%200), 255}
		draw.Draw(img, image.Rect(x, y, x+scrollTileW, y+scrollTileH), &image.Uniform{C: col}, image.Point{}, draw.Src)
	}
	return ebiten.NewImageFromImage(img)
}
