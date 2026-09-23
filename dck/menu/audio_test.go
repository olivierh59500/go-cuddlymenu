package menu

import (
	"encoding/binary"
	"github.com/olivierh59500/democonstructionkit/sound"
	"io"
	"testing"

	gameassets "go-cuddlymenu/dck/assets"
)

func TestMusicReadProducesStereoWithoutAllocating(t *testing.T) {
	data, err := gameassets.Files.ReadFile("menu/menu.ym")
	if err != nil {
		t.Fatal(err)
	}

	player, err := sound.Open("menu.ym", data, sound.Options{SampleRate: sampleRate, Loop: true})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := player.Close(); err != nil {
			t.Errorf("Close returned an error: %v", err)
		}
	}()

	pcm := make([]byte, 4096*8)
	n, err := player.Read(pcm)
	if err != nil {
		t.Fatalf("Read returned an error for a looping stream: %v", err)
	}
	if n != len(pcm) {
		t.Fatalf("Read returned %d bytes, want %d", n, len(pcm))
	}

	nonSilent := false
	for offset := 0; offset < n; offset += 8 {
		left := binary.LittleEndian.Uint32(pcm[offset : offset+4])
		right := binary.LittleEndian.Uint32(pcm[offset+4 : offset+8])
		if left != right {
			t.Fatalf("PCM frame %d differs between channels: left=%d right=%d", offset/8, left, right)
		}
		if left != 0 {
			nonSilent = true
		}
	}
	if !nonSilent {
		t.Fatal("generated PCM block is silent")
	}

	var readN int
	var readErr error
	allocations := testing.AllocsPerRun(100, func() {
		readN, readErr = player.Read(pcm)
	})
	if readErr != nil {
		t.Fatalf("Read returned an error while measuring allocations: %v", readErr)
	}
	if readN != len(pcm) {
		t.Fatalf("Read returned %d bytes while measuring allocations, want %d", readN, len(pcm))
	}
	if allocations != 0 {
		t.Fatalf("Read allocated %.2f times per call, want 0", allocations)
	}
}

func TestMusicReadAfterCloseFails(t *testing.T) {
	data, err := gameassets.Files.ReadFile("menu/menu.ym")
	if err != nil {
		t.Fatal(err)
	}
	player, err := sound.Open("menu.ym", data, sound.Options{SampleRate: sampleRate, Loop: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := player.Close(); err != nil {
		t.Fatal(err)
	}
	if n, err := player.Read(make([]byte, 64)); n != 0 || err != io.ErrClosedPipe {
		t.Fatalf("closed music read = %d, %v; want 0, ErrClosedPipe", n, err)
	}
}
