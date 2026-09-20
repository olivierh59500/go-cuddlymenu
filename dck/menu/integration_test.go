package menu

import "testing"

func TestDoorHandoffWaitsForLoaderAndFiresOnce(t *testing.T) {
	g := &Game{}
	var received []string
	g.SetScreenHandler(func(name string) { received = append(received, name) })
	g.startLoading("BIG_SPRITE")
	for i := 0; i < 119; i++ {
		g.updateLoading()
	}
	if len(received) != 0 || !g.IsLoading() {
		t.Fatal("door completed before its loading interval")
	}
	g.updateLoading()
	if len(received) != 1 || received[0] != "BIG_SPRITE" || g.IsLoading() {
		t.Fatal("door handoff failed", received)
	}
	g.updateLoading()
	if len(received) != 1 {
		t.Fatal("door dispatched twice")
	}
}
