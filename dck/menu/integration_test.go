package menu

import "testing"

func TestDoorHandoffDelegatesTransitionToController(t *testing.T) {
	g := &Game{}
	var received []string
	g.SetScreenHandler(func(name string) { received = append(received, name) })
	g.startLoading("BIG_SPRITE")
	if len(received) != 1 || received[0] != "BIG_SPRITE" || g.IsLoading() {
		t.Fatal("door handoff failed", received)
	}
}
