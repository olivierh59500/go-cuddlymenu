package menu

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewGameDefersCRTResources(t *testing.T) {
	game := NewGame()
	if game.crtShader != nil || game.screenCanvas != nil || game.crtCanvas != nil {
		t.Fatal("CRT resources must not be allocated before the effect is enabled")
	}

	game.initCRT()
	if game.crtShader == nil || game.screenCanvas == nil || game.crtCanvas == nil {
		t.Fatal("initCRT did not create every CRT resource")
	}

	shader := game.crtShader
	screenCanvas := game.screenCanvas
	crtCanvas := game.crtCanvas
	game.initCRT()
	if game.crtShader != shader || game.screenCanvas != screenCanvas || game.crtCanvas != crtCanvas {
		t.Fatal("initCRT must be idempotent")
	}

	game.useCRT = true
	game.Draw(ebiten.NewImage(screenWidth, screenHeight))
}

func TestBackgroundOnlyCoversVisibleViewport(t *testing.T) {
	game := NewGame()
	wantWidth := gameWidth + tileSize
	wantHeight := gameHeight + tileSize
	if got := game.background.Bounds().Dx(); got != wantWidth {
		t.Fatalf("background width = %d, want %d", got, wantWidth)
	}
	if got := game.background.Bounds().Dy(); got != wantHeight {
		t.Fatalf("background height = %d, want %d", got, wantHeight)
	}
}
