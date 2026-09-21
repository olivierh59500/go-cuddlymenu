package screens

import "testing"

func TestCreditsExcludedAndEveryRemainingScreenConstructs(t *testing.T) {
	if _, ok := Find("credits"); ok {
		t.Fatal("Unsupported credits screen must stay outside the native scope")
	}
	if DoorID("CREDITS") != "" {
		t.Fatal("excluded credits door still dispatches")
	}
	if len(Catalog()) != 15 {
		t.Fatal("expected thirteen demos plus introduction and reset")
	}
	for _, d := range Catalog() {
		t.Run(d.ID, func(t *testing.T) {
			if !d.Ready {
				t.Fatal("screen is not enabled")
			}
			s, err := New(d.ID)
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()
			if err = s.Update(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIntroductionStartsMusicAtItsVisualTransition(t *testing.T) {
	s, err := New("intro")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if cue := s.TakeAudioCue(); cue == nil || cue.File != "" {
		t.Fatal("introduction must begin silently")
	}
	for tick := 1; tick <= 250; tick++ {
		if err = s.Update(); err != nil {
			t.Fatal(err)
		}
		cue := s.TakeAudioCue()
		if tick < 250 && cue != nil {
			t.Fatalf("music started early at %d", tick)
		}
		if tick == 250 && (cue == nil || cue.File != "cuddlyintro.mp3" || !cue.Loop) {
			t.Fatal("main intro music cue missing")
		}
	}
}

func TestSpreadpointRetainsAllFourIntroTracks(t *testing.T) {
	s, err := New("spreadpoint")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if cue := s.TakeAudioCue(); cue == nil || cue.File != "intro1.mp3" || cue.Loop {
		t.Fatal("wrong first intro track")
	}
	want := map[int]AudioCue{180: {File: "intro2.mp3"}, 361: {File: "intro3.mp3"}, 542: {File: "intro4.mp3"}, 723: {File: "master.mp3", Loop: true}}
	for tick := 1; tick <= 724; tick++ {
		s.Update()
		cue := s.TakeAudioCue()
		expected, ok := want[tick]
		if ok {
			if cue == nil || *cue != expected {
				t.Fatalf("tick %d: got %+v, want %+v", tick, cue, expected)
			}
		} else if cue != nil {
			t.Fatalf("unexpected music switch at %d", tick)
		}
	}
}
