package app

import "testing"

func TestIntroMenuAndDoorTransitions(t *testing.T) {
	g, err := New(Config{Screen: "intro", Muted: true})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if g.scene == nil || g.scene.Descriptor.ID != "intro" || g.track != "" {
		t.Fatal("full production must start with silent Calvin phase")
	}
	if err = g.Begin("menu"); err != nil {
		t.Fatal(err)
	}
	if g.transition == nil || g.target != "menu" {
		t.Fatal("intro exit bypassed loader")
	}
	for !g.transition.Done() {
		g.transition.Update()
	}
	if err = g.open(g.target); err != nil {
		t.Fatal(err)
	}
	if g.scene != nil || g.transition != nil {
		t.Fatal("menu not reached")
	}
	if err = g.Begin("starwars"); err != nil {
		t.Fatal(err)
	}
	if g.transition == nil || g.target != "starwars" {
		t.Fatal("screen selection bypassed countdown")
	}
}

func TestTouchNavigationUsesPixelSideMargins(t *testing.T) {
	g, err := New(Config{Screen: "intro", Muted: true, TouchControls: true})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	w, _ := g.Layout(2424, 1080)
	left := (w - 768) / 2
	for _, b := range g.buttons() {
		if b.rect.Max.X > left {
			t.Fatalf("%s overlaps intro artwork", b.label)
		}
	}
}

func TestTouchCRTControlTracksMenuState(t *testing.T) {
	g, err := New(Config{Muted: true, TouchControls: true})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	g.Layout(2424, 1080)
	findCRT := func() string {
		for _, b := range g.buttons() {
			if b.action == "crt" {
				return b.label
			}
		}
		return ""
	}
	if got := findCRT(); got != "CRT OFF" {
		t.Fatalf("initial menu button = %q, want CRT OFF", got)
	}
	if !g.menu.SetCRT(true) {
		t.Fatal("could not enable the CRT material")
	}
	if got := findCRT(); got != "CRT ON" {
		t.Fatalf("enabled menu button = %q, want CRT ON", got)
	}
	if err := g.open("intro"); err != nil {
		t.Fatal(err)
	}
	if got := findCRT(); got != "" {
		t.Fatalf("CRT button is shown over a demo screen: %q", got)
	}
}
