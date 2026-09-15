package menu

import (
	"encoding/binary"
	"testing"

	gameassets "go-cuddlymenu/assets"
)

func TestYMPlayerReadProducesStereoWithoutAllocating(t *testing.T) {
	data, err := gameassets.Files.ReadFile("menu/menu.ym")
	if err != nil {
		t.Fatal(err)
	}

	player, err := NewYMPlayer(data, sampleRate, true)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := player.Close(); err != nil {
			t.Errorf("Close returned an error: %v", err)
		}
	}()

	pcm := make([]byte, 4096*4)
	n, err := player.Read(pcm)
	if err != nil {
		t.Fatalf("Read returned an error for a looping stream: %v", err)
	}
	if n != len(pcm) {
		t.Fatalf("Read returned %d bytes, want %d", n, len(pcm))
	}

	nonSilent := false
	for offset := 0; offset < n; offset += 4 {
		left := int16(binary.LittleEndian.Uint16(pcm[offset : offset+2]))
		right := int16(binary.LittleEndian.Uint16(pcm[offset+2 : offset+4]))
		if left != right {
			t.Fatalf("PCM frame %d differs between channels: left=%d right=%d", offset/4, left, right)
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

func TestYMPlayerReadAfterCloseReturnsSilence(t *testing.T) {
	data, err := gameassets.Files.ReadFile("menu/menu.ym")
	if err != nil {
		t.Fatal(err)
	}

	player, err := NewYMPlayer(data, sampleRate, true)
	if err != nil {
		t.Fatal(err)
	}
	if err := player.Close(); err != nil {
		t.Fatal(err)
	}

	pcm := make([]byte, 64)
	for i := range pcm {
		pcm[i] = 0xff
	}
	n, err := player.Read(pcm)
	if n != len(pcm) {
		t.Fatalf("Read returned %d bytes, want %d", n, len(pcm))
	}
	if err == nil {
		t.Fatal("Read after Close returned a nil error")
	}
	for i, sampleByte := range pcm {
		if sampleByte != 0 {
			t.Fatalf("PCM byte %d after Close is %d, want silence", i, sampleByte)
		}
	}
}
