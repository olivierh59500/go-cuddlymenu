// Command checkcuddly validates embedded images and decodes every active music asset.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/olivierh59500/ym-player/pkg/stsound"
	media "go-cuddlymenu/assets/cuddly"
	"go-cuddlymenu/dck/screens"
)

type asset struct {
	Path          string
	Width, Height int
}
type music struct {
	Screen, Path, SHA256, Format, Title string
	DurationMS                          uint32
	NonzeroSamples                      int
}
type report struct {
	AssetFiles    int
	Images        []asset
	Music         []music
	NativeHasSNDH bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	r := report{}
	err := fs.WalkDir(media.Files, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		r.AssetFiles++
		switch strings.ToLower(path.Ext(p)) {
		case ".png", ".jpg", ".jpeg", ".gif":
			b, err := media.Files.ReadFile(p)
			if err != nil {
				return err
			}
			cfg, _, err := image.DecodeConfig(bytes.NewReader(b))
			if err != nil {
				return fmt.Errorf("%s: %w", p, err)
			}
			r.Images = append(r.Images, asset{p, cfg.Width, cfg.Height})
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err = fs.WalkDir(media.Files, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(strings.ToLower(p), ".sndh") {
			r.NativeHasSNDH = true
		}
		return nil
	}); err != nil {
		return err
	}
	if r.NativeHasSNDH {
		return fmt.Errorf("SNDH leaked into native assets")
	}
	selections := screens.Catalog()
	selections = append(selections, screens.Descriptor{ID: "loader", Directory: "menu/resources", Music: "loader.wav", Ready: true})
	for i := 1; i <= 4; i++ {
		d, _ := screens.Find("spreadpoint")
		d.ID = fmt.Sprintf("spreadpoint-intro-%d", i)
		d.Music = fmt.Sprintf("intro%d.mp3", i)
		selections = append(selections, d)
	}
	for _, d := range selections {
		if !d.Ready {
			continue
		}
		name := d.Directory + "/" + d.Music
		if strings.HasPrefix(d.Music, "@") {
			name = d.Music[1:]
		}
		b, err := media.Files.ReadFile(name)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		m := music{Screen: d.ID, Path: name, SHA256: hex.EncodeToString(sum[:])}
		if path.Ext(name) == ".ym" {
			p := stsound.CreateWithRate(48000)
			if err = p.LoadMemory(b); err != nil {
				p.Destroy()
				return fmt.Errorf("%s: %w", name, err)
			}
			info := p.GetInfo()
			m.Format, m.Title, m.DurationMS = info.SongType, info.SongName, uint32(info.MusicTimeInMs)
			samples := make([]int16, 96000)
			p.Compute(samples, len(samples))
			for _, v := range samples {
				if v != 0 {
					m.NonzeroSamples++
				}
			}
			p.Destroy()
		} else {
			var decoded io.Reader
			if path.Ext(name) == ".wav" {
				decoded, err = wav.DecodeWithSampleRate(48000, bytes.NewReader(b))
				m.Format = "WAV decoded to stereo PCM"
			} else {
				decoded, err = mp3.DecodeWithSampleRate(48000, bytes.NewReader(b))
				m.Format = "MP3 decoded to stereo PCM"
			}
			if err != nil {
				return err
			}
			pcm := make([]byte, 96000*4)
			n, err := io.ReadFull(decoded, pcm)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return err
			}
			for i := 0; i+1 < n; i += 2 {
				if pcm[i] != 0 || pcm[i+1] != 0 {
					m.NonzeroSamples++
				}
			}
		}
		if m.NonzeroSamples == 0 {
			return fmt.Errorf("silent music: %s", name)
		}
		r.Music = append(r.Music, m)
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll("captures", 0755); err != nil {
		return err
	}
	if err = os.WriteFile("captures/assets.json", append(b, '\n'), 0644); err != nil {
		return err
	}
	fmt.Printf("Verified %d embedded files, %d images and %d native music selections; no SNDH in native assets.\n", r.AssetFiles, len(r.Images), len(r.Music))
	return nil
}
