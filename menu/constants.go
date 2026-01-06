package main

const (
	screenWidth  = 768
	screenHeight = 536

	gameWidth   = screenWidth
	gameHeight  = 400
	gameOffsetX = 0
	gameOffsetY = 0

	scrollOffsetX = 0
	scrollOffsetY = gameHeight + (screenHeight-gameHeight-scrollHeight)/2
	scrollWidth   = screenWidth
	scrollHeight  = 80

	tileSize      = 32
	dudeSize      = 64
	carebearTileW = 32
	carebearTileH = 20
	scrollTileW   = 32
	scrollTileH   = 80

	bounceSpeed               = 7
	scrollSpeed               = 8
	autoPilotActivateDuration = 60 * 60 * 2

	sampleRate = 44100
)
