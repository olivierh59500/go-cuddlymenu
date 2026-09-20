// Command capture records deterministic frames from one native Cuddly screen.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	capture "github.com/olivierh59500/democonstructionkit/fidelity/ebiten"
	"go-cuddlymenu/dck/loader"
	"go-cuddlymenu/dck/screens"
)

func main() {
	id := flag.String("screen", "colorshock", "native screen ID")
	loading := flag.String("loader", "", "capture an original loader by door name, e.g. BIG_SPRITE")
	out := flag.String("out", "captures/native", "capture directory")
	framesFlag := flag.String("frames", "0,60,240,600", "capture ticks")
	flag.Parse()
	var frames []int
	for _, part := range strings.Split(*framesFlag, ",") {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 || (len(frames) > 0 && n <= frames[len(frames)-1]) {
			fmt.Fprintln(os.Stderr, "invalid frames")
			os.Exit(1)
		}
		frames = append(frames, n)
	}
	if *loading != "" {
		err := capture.Run(capture.Config{Directory: filepath.Join(*out, "cuddly_loader_"+*loading), Frames: frames, Width: loader.Width, Height: loader.Height}, func() (ebiten.Game, error) { return loader.New(*loading) })
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	d, ok := screens.Find(*id)
	if !ok || !d.Ready {
		fmt.Fprintln(os.Stderr, "screen not available:", *id)
		os.Exit(1)
	}
	err := capture.Run(capture.Config{Directory: filepath.Join(*out, d.Page), Frames: frames, Width: d.Width, Height: d.Height}, func() (ebiten.Game, error) { return screens.New(*id) })
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
