package screens

import (
	"compress/gzip"
	"encoding/binary"
	"io"
	"os"
	"testing"
	"time"

	"github.com/olivierh59500/democonstructionkit/motion"
	media "go-cuddlymenu/assets/cuddly"
)

func knucklebusterSignals(t *testing.T) *motion.SampledSignals {
	t.Helper()
	bank, err := media.Files.ReadFile("tex/drummer-triggers.bin")
	if err != nil {
		t.Fatal(err)
	}
	signals, err := motion.NewSampledSignals(motion.SampledSignalsConfig{
		Channels: 4, Rate: 50, Masks: bank, Loop: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return signals
}

// The words are captured from the original 68000 music/animation routine,
// including its six-frame drum holds and independent twenty-frame head hold.
func TestKnucklebusterMatchesEveryNativeAnimationFlag(t *testing.T) {
	signals := knucklebusterSignals(t)
	file, err := os.Open("testdata/knucklebuster-native-flags.bin.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 54750*8 || signals.Duration() != 1095*time.Second {
		t.Fatal("incomplete native animation or soundtrack duration")
	}
	for frame := 0; frame < 54750; frame++ {
		for channel := 0; channel < 4; channel++ {
			want := binary.BigEndian.Uint16(data[frame*8+channel*2:]) == 0
			if got := signals.At(time.Duration(frame)*time.Second/50, channel); got != want {
				t.Fatalf("native frame %d channel %d: got %v, want %v", frame, channel, got, want)
			}
		}
	}
}

func TestKnucklebusterKeepsHeadAndDrumsIndependent(t *testing.T) {
	signals := knucklebusterSignals(t)
	for _, sample := range []struct {
		frame int
		on    [4]bool
	}{
		{388, [4]bool{}},
		{389, [4]bool{false, false, true, true}},
		{395, [4]bool{false, false, false, true}},
		{409, [4]bool{}},
		{437, [4]bool{false, true, false, false}},
		{443, [4]bool{}},
		{8477, [4]bool{true, true, false, false}},
	} {
		for channel, want := range sample.on {
			if got := signals.At(time.Duration(sample.frame)*time.Second/50, channel); got != want {
				t.Fatalf("frame %d channel %d: got %v, want %v", sample.frame, channel, got, want)
			}
		}
	}
}

func TestSceneMusicClockCanPauseAndSeekWithoutChangingAnimationRate(t *testing.T) {
	s := &Scene{}
	if err := s.SetAnimationRate(50); err != nil {
		t.Fatal(err)
	}
	position := 170 * time.Second
	s.SetMusicClock(func() time.Duration { return position })
	for _, want := range []time.Duration{170 * time.Second, 170 * time.Second, 8 * time.Second} {
		position = want
		if got := s.MusicPosition(); got != want || s.AnimationRate() != 50 {
			t.Fatalf("audio clock got %v, want %v", got, want)
		}
	}
	s.SetMusicClock(nil)
	s.render = func() {}
	for i := 0; i < 50; i++ {
		if err := s.Update(); err != nil {
			t.Fatal(err)
		}
	}
	if s.MusicPosition() != time.Second {
		t.Fatal("muted fallback accelerated the native clock")
	}
	if err := s.SetAnimationRate(60); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 60; i++ {
		s.Update()
	}
	if s.MusicPosition() != 2*time.Second {
		t.Fatal("changing display rate reset elapsed music time")
	}
}
