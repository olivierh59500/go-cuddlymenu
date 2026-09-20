package menu

// UseExternalAudio lets a screen controller own the single application player.
// Call before the first Update. Standalone menu behavior remains the default.
func (g *Game) UseExternalAudio()                     { g.audioReady = true }
func (g *Game) SetScreenHandler(handler func(string)) { g.screenHandler = handler }
func (g *Game) IsLoading() bool                       { return g.loading.Active }

func (g *Game) Close() error {
	if g.audioPlayer != nil {
		g.audioPlayer.Close()
		g.audioPlayer = nil
	}
	if g.ymPlayer != nil {
		g.ymPlayer.Close()
		g.ymPlayer = nil
	}
	return nil
}
