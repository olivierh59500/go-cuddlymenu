package app

import (
	"go-cuddlymenu/dck/screens"
	"testing"
)

func TestTourVisitsEveryScreenOnce(t *testing.T) {
	seen := map[string]bool{"intro": true, "reset": true}
	for _, id := range tourScreens {
		if seen[id] {
			t.Fatalf("duplicate tour screen %s", id)
		}
		seen[id] = true
		name := screens.DoorName(id)
		if name == "" || screens.DoorID(name) != id {
			t.Fatalf("screen %s has no reversible menu door", id)
		}
	}
	for _, screen := range screens.Catalog() {
		if !seen[screen.ID] {
			t.Fatalf("tour omits %s", screen.ID)
		}
	}
}
