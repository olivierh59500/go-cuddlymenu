package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Vec2 struct {
	X float64
	Y float64
}

type Model struct {
	Position           Vec2
	Direction          int
	Moving             bool
	Thrusting          bool
	JustLanded         int
	BounceDisplacement int
	ThrustSpeed        float64
	FallingSpeed       float64
	CurrentFrame       int
	ScrollerPosition   int
}

type AutoPilot struct {
	ActivateIn           int
	DontThrustForAwhile  int
	TimeSinceLastThrust  int
	WaitToLoad           int
	NowLoadScreen        bool
	Stuck                bool
	KeepWalkingForAwhile int
	NextScreen           int
}

type DemoScreen struct {
	X    int
	Y    int
	Name string
}

type LoaderState struct {
	Active     bool
	ScreenName string
	Timer      int
}

type DudeAnimations struct {
	MoveRight   Animation
	MoveLeft    Animation
	ThrustRight Animation
	ThrustLeft  Animation
}

type Game struct {
	assets       *Assets
	audioContext *audio.Context
	audioPlayer  *audio.Player
	ymPlayer     *YMPlayer

	mapTiles      *TileSet
	dudeTiles     *TileSet
	carebearTiles *TileSet
	scrollTiles   *TileSet

	mapLevel       *TileMap
	scrollerLevel  *TileMap
	scrollerLength int

	background   *ebiten.Image
	gameCanvas   *ebiten.Image
	screenCanvas *ebiten.Image

	sineSprites *SineSprites
	animations  DudeAnimations

	model        Model
	autoPilot    AutoPilot
	loading      LoaderState
	thrustOff    int
	simTime      float64
	carebearTime float64

	crtShader *ebiten.Shader
	useCRT    bool
}

var bouncingAnimation = []int{0, 3, 5, 6, 5, 3, 0, 1, 2, 3, 2, 1, 0}

func NewGame() *Game {
	maxTile := maxTileIndex(cuddlyMap)
	assets := LoadAssets(filepath.Join("assets", "menu"), maxTile)

	g := &Game{
		assets:       assets,
		useCRT:       false,
		gameCanvas:   ebiten.NewImage(gameWidth, gameHeight),
		screenCanvas: ebiten.NewImage(screenWidth, screenHeight),
	}

	g.mapTiles = NewTileSet(assets.Tiles, tileSize, tileSize)
	g.dudeTiles = NewTileSet(assets.Dude, dudeSize, dudeSize)
	g.carebearTiles = NewTileSet(assets.Carebears, carebearTileW, carebearTileH)
	g.scrollTiles = NewTileSet(assets.Chrome, scrollTileW, scrollTileH)

	g.mapLevel = NewTileMap(cuddlyMap, g.mapTiles)
	scrollMap := BuildScrollMap(scrollTextData)
	g.scrollerLevel = NewTileMap([][]int{scrollMap}, g.scrollTiles)
	g.scrollerLength = len(scrollMap) * scrollTileW

	g.background = g.buildBackground()
	g.sineSprites = &SineSprites{Tiles: g.carebearTiles}

	g.animations = DudeAnimations{
		MoveRight:   Animation{Duration: 0.35, Indices: []int{2, 3, 4, 5, 6, 7, 8, 9}, Loop: true},
		MoveLeft:    Animation{Duration: 0.35, Indices: []int{12, 13, 14, 15, 16, 17, 18, 19}, Loop: true},
		ThrustRight: Animation{Duration: 0.075, Indices: []int{0, 1}, Loop: true},
		ThrustLeft:  Animation{Duration: 0.075, Indices: []int{10, 11}, Loop: true},
	}

	g.Reset()
	g.initAudio()
	g.initShader()

	return g
}

func (g *Game) Reset() {
	g.model = Model{
		Position:     Vec2{X: 320, Y: 450},
		Direction:    1,
		CurrentFrame: 6,
	}
	g.autoPilot = AutoPilot{
		ActivateIn:          autoPilotActivateDuration,
		TimeSinceLastThrust: 200,
		WaitToLoad:          80,
	}
	g.loading = LoaderState{}
	g.thrustOff = 0
	g.simTime = 0
	g.carebearTime = 0
}

func (g *Game) initAudio() {
	g.audioContext = audio.NewContext(sampleRate)
	if len(g.assets.MenuYM) == 0 {
		return
	}
	var err error
	g.ymPlayer, err = NewYMPlayer(g.assets.MenuYM, sampleRate, true)
	if err != nil {
		log.Printf("failed to create YM player: %v", err)
		return
	}
	g.audioPlayer, err = g.audioContext.NewPlayer(g.ymPlayer)
	if err != nil {
		log.Printf("failed to create audio player: %v", err)
		g.ymPlayer.Close()
		g.ymPlayer = nil
		return
	}
	g.audioPlayer.Play()
}

func (g *Game) initShader() {
	shader, err := ebiten.NewShader([]byte(crtShaderSrc))
	if err != nil {
		log.Printf("failed to compile CRT shader: %v", err)
		return
	}
	g.crtShader = shader
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.useCRT = !g.useCRT
	}

	if g.loading.Active {
		g.updateLoading()
		return nil
	}

	left, right, thrust, load, anyKey := g.readInput()

	if anyKey {
		g.autoPilot.ActivateIn = autoPilotActivateDuration
	} else {
		g.autoPilot.ActivateIn--
	}

	if g.autoPilot.ActivateIn <= 0 {
		move := g.autoPilotMovement()
		left = move.left
		right = move.right
		thrust = move.thrust
		if g.autoPilot.NowLoadScreen {
			load = true
		}
	}

	g.integrate(left, right, thrust)
	g.carebearTime += 1.0 / 60.0
	g.simTime += 1.0 / 60.0
	g.handleLoad(load)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawScene(g.screenCanvas)
	if g.useCRT && g.crtShader != nil {
		op := &ebiten.DrawRectShaderOptions{}
		op.Images[0] = g.screenCanvas
		screen.DrawRectShader(screenWidth, screenHeight, g.crtShader, op)
		return
	}
	screen.DrawImage(g.screenCanvas, nil)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) readInput() (left, right, thrust, load, anyKey bool) {
	left = ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyZ)
	right = ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyX)
	thrust = ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyEnter)
	load = ebiten.IsKeyPressed(ebiten.KeySpace)
	anyKey = len(inpututil.AppendPressedKeys(nil)) > 0
	return
}

type movement struct {
	left   bool
	right  bool
	thrust bool
}

func (g *Game) autoPilotMovement() movement {
	left := false
	right := false
	thrust := false

	screen := g.autoPilot.NextScreen
	if screen < 0 || screen >= len(demoScreens) {
		screen = 0
	}
	// Target is the center of the 4x3 door area.
	target := Vec2{
		X: float64(demoScreens[screen].X*tileSize + 64),
		Y: float64(demoScreens[screen].Y*tileSize + 64),
	}
	current := Vec2{
		X: g.model.Position.X,
		Y: g.model.Position.Y,
	}

	currentX := int(current.X) &^ 0xf
	targetX := int(target.X) &^ 0xf
	currentY := int(current.Y)
	onLevel := g.haveLanded(&g.model)
	if onLevel && (currentY&^0x1f) < (int(target.Y)&^0xf) && currentY < 550 {
		g.autoPilot.Stuck = true
		g.autoPilot.KeepWalkingForAwhile = 30
	}
	if currentY >= 550 {
		g.autoPilot.Stuck = false
	}

	if g.autoPilot.Stuck {
		if g.model.Direction == 0 {
			left = true
		} else {
			right = true
		}
		if !onLevel {
			g.autoPilot.KeepWalkingForAwhile--
			if g.autoPilot.KeepWalkingForAwhile <= 0 {
				g.autoPilot.Stuck = false
			}
		}
	} else {
		if currentX == targetX {
			g.autoPilot.WaitToLoad--
		}
		if g.autoPilot.WaitToLoad <= 0 {
			g.autoPilot.NowLoadScreen = true
		}
		if currentX < targetX {
			right = true
		}
		if currentX > targetX {
			left = true
		}
		if current.Y > target.Y && g.autoPilot.DontThrustForAwhile < 70 {
			thrust = true
		}
	}

	if g.autoPilot.TimeSinceLastThrust > 0 {
		g.autoPilot.TimeSinceLastThrust--
	}
	if g.autoPilot.TimeSinceLastThrust <= 0 {
		g.autoPilot.TimeSinceLastThrust = rand.Intn(100) + 100
	}
	if g.autoPilot.DontThrustForAwhile > 0 {
		g.autoPilot.DontThrustForAwhile--
	}
	if g.autoPilot.DontThrustForAwhile <= 0 {
		g.autoPilot.DontThrustForAwhile = 100
	}

	return movement{left: left, right: right, thrust: thrust}
}

func (g *Game) integrate(left, right, thrust bool) {
	m := &g.model

	if m.BounceDisplacement > 0 {
		m.BounceDisplacement--
	}

	if left && !right {
		m.Position.X -= bounceSpeed
		m.Direction = 0
		m.Moving = true
	}
	if right && !left {
		m.Position.X += bounceSpeed
		m.Direction = 1
		m.Moving = true
	}
	if !right && !left {
		m.Moving = false
	}

	if thrust {
		m.Thrusting = true
		if m.ThrustSpeed <= 0 {
			m.ThrustSpeed = 3
		}
		if m.ThrustSpeed < 8 {
			m.ThrustSpeed += 1
		}
		m.FallingSpeed = 0
	} else {
		m.Thrusting = false
		if m.Position.Y < tileSize {
			m.ThrustSpeed = 0
		}
	}

	if m.ThrustSpeed <= 0 && m.JustLanded != 1 {
		if m.FallingSpeed < bounceSpeed {
			m.FallingSpeed += 0.5
		}
	}

	for i := 0; i < int(math.Floor(m.FallingSpeed))+1; i++ {
		if !g.haveLanded(m) {
			m.Position.Y += 1
			m.JustLanded = 0
		} else if m.JustLanded < 1 {
			m.JustLanded = 1
			m.BounceDisplacement = len(bouncingAnimation)
			m.ThrustSpeed = 0
		}
	}

	if m.ThrustSpeed > 0 {
		m.ThrustSpeed -= 0.5
	}
	if m.ThrustSpeed > 0 && g.thrustOff == 0 {
		m.Position.Y -= math.Floor(m.ThrustSpeed)
	}

	mapHeight := g.mapLevel.HeightPx - (tileSize - 4)
	if m.Position.X < float64(9*tileSize) {
		m.Position.X = float64(9 * tileSize)
	}
	if m.Position.Y < 0 {
		m.Position.Y = 0
		g.thrustOff = 10
	}
	if g.thrustOff > 0 {
		g.thrustOff--
		m.Position.Y += 2
	}
	if m.Position.X >= float64(455*tileSize) {
		m.Position.X = float64(455 * tileSize)
	}
	if m.Position.Y >= float64(mapHeight-dudeSize) {
		m.Position.Y = float64(mapHeight - dudeSize)
	}

	m.ScrollerPosition += scrollSpeed
}

func (g *Game) haveLanded(m *Model) bool {
	const floorTile = 69
	if len(cuddlyMap) == 0 {
		return false
	}
	x := int(m.Position.X) / tileSize
	y := int(m.Position.Y) / tileSize
	if y+2 >= len(cuddlyMap) || y < 0 {
		return false
	}
	if x < 0 || x+1 >= len(cuddlyMap[0]) {
		return false
	}
	if int(m.Position.Y)%tileSize != 0 {
		return false
	}

	row := cuddlyMap[y+2]
	if row[x] == floorTile || row[x+1] == floorTile {
		return true
	}
	if x < len(row)-2 && row[x+2] == floorTile {
		return true
	}
	return false
}

func (g *Game) handleLoad(load bool) {
	if !load || g.loading.Active {
		return
	}
	pX := int(g.model.Position.X) / tileSize
	pY := int(g.model.Position.Y) / tileSize
	for _, t := range demoScreens {
		if pX >= t.X && pX <= t.X+3 && pY >= t.Y && pY < t.Y+3 {
			g.startLoading(t.Name)
			return
		}
	}
}

func (g *Game) startLoading(name string) {
	g.loading = LoaderState{
		Active:     true,
		ScreenName: name,
		Timer:      120,
	}
	g.autoPilot.NowLoadScreen = false
	g.autoPilot.WaitToLoad = 80
	g.advanceAutoPilot()
	if g.audioPlayer != nil {
		g.audioPlayer.Pause()
	}
}

func (g *Game) advanceAutoPilot() {
	next := g.autoPilot.NextScreen + 1
	if next >= len(demoScreens)-1 {
		next = 0
	}
	g.autoPilot.NextScreen = next
}

func (g *Game) updateLoading() {
	g.loading.Timer--
	if g.loading.Timer > 0 {
		return
	}
	g.loading.Active = false
	g.autoPilot.NowLoadScreen = false
	g.autoPilot.WaitToLoad = 80
	if g.audioPlayer != nil && !g.audioPlayer.IsPlaying() {
		g.audioPlayer.Play()
	}
}

func (g *Game) drawScene(dst *ebiten.Image) {
	dst.Fill(color.Black)

	if g.scrollerLength > 0 {
		scrollX := g.model.ScrollerPosition % g.scrollerLength
		g.scrollerLevel.Draw(dst, scrollX, 0, scrollOffsetX, scrollOffsetY, scrollWidth, scrollHeight)
	}

	g.gameCanvas.Fill(color.Black)
	mapX, mapY, dudeX, dudeY, bounce := g.calculatePositions()
	g.drawBackground(g.gameCanvas, mapX, mapY)
	g.mapLevel.Draw(g.gameCanvas, mapX, mapY, 0, 0, gameWidth, gameHeight)
	frame := g.calculateFrame()
	g.drawDude(g.gameCanvas, dudeX, dudeY-bounce, frame)
	g.sineSprites.Draw(g.gameCanvas, g.carebearTime)

	var op ebiten.DrawImageOptions
	op.GeoM.Translate(gameOffsetX, gameOffsetY)
	dst.DrawImage(g.gameCanvas, &op)

	if g.loading.Active {
		g.drawLoading(dst)
	}
}

func (g *Game) calculatePositions() (int, int, float64, float64, float64) {
	m := g.model
	dudePosX := int(math.Round(m.Position.X))
	dudePosY := int(math.Round(m.Position.Y))

	mapWidth := g.mapLevel.WidthPx
	mapHeight := g.mapLevel.HeightPx - (tileSize - 4)

	dudeScreenX := 0
	mapX := 0
	if dudePosX <= gameWidth/2-dudeSize/2 {
		dudeScreenX = dudePosX
		mapX = 0
	} else if dudePosX > mapWidth-gameWidth/2-dudeSize/2 {
		dudeScreenX = gameWidth - (mapWidth - dudePosX)
		mapX = mapWidth - gameWidth
	} else {
		dudeScreenX = gameWidth/2 - dudeSize/2
		mapX = dudePosX - (gameWidth/2 - dudeSize/2)
	}

	dudeScreenY := 0
	mapY := 0
	if dudePosY <= gameHeight/2-dudeSize/2 {
		dudeScreenY = dudePosY
		mapY = 0
	} else if dudePosY > mapHeight-gameHeight/2-dudeSize/2 {
		dudeScreenY = gameHeight - (mapHeight - dudePosY)
		mapY = mapHeight - gameHeight
	} else {
		dudeScreenY = gameHeight/2 - dudeSize/2
		mapY = dudePosY - (gameHeight/2 - dudeSize/2)
	}

	bounce := g.calculateBounce()
	return mapX, mapY, float64(dudeScreenX), float64(dudeScreenY), bounce
}

func (g *Game) calculateBounce() float64 {
	if len(bouncingAnimation) == 0 {
		return 0
	}
	idx := len(bouncingAnimation) - g.model.BounceDisplacement - 1
	if idx < 0 || idx >= len(bouncingAnimation) {
		return 0
	}
	bounce := float64(bouncingAnimation[idx])
	falling := g.model.FallingSpeed
	if falling < 3 {
		falling = 3
	}
	bounce *= (falling / bounceSpeed) * 0.9
	return bounce
}

func (g *Game) calculateFrame() int {
	frame := g.model.CurrentFrame
	if g.model.Thrusting {
		if g.model.Direction < 1 {
			frame = g.animations.ThrustLeft.Current(g.simTime)
		} else {
			frame = g.animations.ThrustRight.Current(g.simTime)
		}
		return frame
	}
	if g.model.Moving {
		if g.model.Direction < 1 {
			frame = g.animations.MoveLeft.Current(g.simTime)
		} else {
			frame = g.animations.MoveRight.Current(g.simTime)
		}
		g.model.CurrentFrame = frame
		return frame
	}

	idle := frame % 10
	if g.model.Direction < 1 {
		return idle + 10
	}
	return idle
}

func (g *Game) drawDude(dst *ebiten.Image, posX, posY float64, frame int) {
	sprite := g.dudeTiles.Tile(frame)
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(posX, posY)
	dst.DrawImage(sprite, &op)
}

func (g *Game) drawBackground(dst *ebiten.Image, posX, posY int) {
	if g.background == nil {
		return
	}
	if posX < 0 {
		posX = 0
	}
	if posY < 0 {
		posY = 0
	}

	deltaX := int(math.Mod(float64(posX)*0.5, float64(tileSize)))
	deltaY := int(math.Mod(float64(posY)*0.5, float64(tileSize)))
	if deltaX < 0 {
		deltaX += tileSize
	}
	if deltaY < 0 {
		deltaY += tileSize
	}

	rect := image.Rect(deltaX, deltaY, deltaX+gameWidth, deltaY+g.mapLevel.HeightPx)
	sub := g.background.SubImage(rect).(*ebiten.Image)
	dst.DrawImage(sub, nil)
}

func (g *Game) buildBackground() *ebiten.Image {
	bgW := gameWidth + tileSize
	bgH := g.mapLevel.HeightPx + tileSize
	img := ebiten.NewImage(bgW, bgH)
	tile := g.mapTiles.Tile(1)
	for y := 0; y < bgH; y += tileSize {
		for x := 0; x < bgW; x += tileSize {
			var op ebiten.DrawImageOptions
			op.GeoM.Translate(float64(x), float64(y))
			img.DrawImage(tile, &op)
		}
	}
	return img
}

func (g *Game) drawLoading(dst *ebiten.Image) {
	overlay := color.RGBA{0, 0, 0, 200}
	ebitenutil.DrawRect(dst, 0, 0, screenWidth, screenHeight, overlay)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("LOADING %s", g.loading.ScreenName), 20, 20)
}

func maxTileIndex(mapData [][]int) int {
	max := 0
	for _, row := range mapData {
		for _, v := range row {
			if v > max {
				max = v
			}
		}
	}
	return max
}

func BuildScrollMap(text string) []int {
	clean := make([]rune, 0, len(text))
	for _, r := range text {
		if r == '\n' || r == '\r' {
			continue
		}
		clean = append(clean, r)
	}

	result := make([]int, 0, len(clean)*3)
	for _, r := range clean {
		p := int(r) - 32
		if p < 0 || p >= len(scrollerCharWidth) {
			p = 0
		}
		blocks := scrollerCharWidth[p]
		for j := 0; j < blocks; j++ {
			result = append(result, p*3+j)
		}
	}
	return result
}

const crtShaderSrc = `
package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	var uv vec2
	uv = texCoord

	// Barrel distortion
	var dc vec2
	dc = uv - 0.5
	dc = dc * (1.0 + dot(dc, dc) * 0.15)
	uv = dc + 0.5

	if uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0 {
		return vec4(0.0, 0.0, 0.0, 1.0)
	}

	var col vec4
	col = imageSrc0At(uv)

	// Scanlines
	var scanline float
	scanline = sin(uv.y * 800.0) * 0.04
	col.rgb = col.rgb - scanline

	// RGB shift
	var rShift float
	var bShift float
	rShift = imageSrc0At(uv + vec2(0.002, 0.0)).r
	bShift = imageSrc0At(uv - vec2(0.002, 0.0)).b
	col.r = rShift
	col.b = bShift

	// Vignette
	var vignette float
	vignette = 1.0 - dot(dc, dc) * 0.5
	col.rgb = col.rgb * vignette

	return col * color
}
`

func main() {
	rand.Seed(time.Now().UnixNano())
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cuddly Demos - Menu")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
