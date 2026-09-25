# Cuddly Demo / DCK

The complete production uses shared Go/Ebitengine effects and embedded assets.

Go downloads the published DCK module pinned in `go.mod` and
`github.com/olivierh59500/ym-player v1.0.0` automatically.

```sh
go run ./dck/cmd/cuddlydemo        # Calvin, animated intro, then the menu.
go run ./dck/cmd/cuddlymenu        # Start directly in the menu.
go run ./dck/cmd/cuddlydemo -list
go run ./dck/cmd/cuddlydemo -screen dna
```

The original Go implementation remains available with `go run ./cmd/cuddlymenu`.

F1 opens the selector. Space or Escape returns through the loading screen, and R
opens Reset. Playback defaults to 60 fixed updates per second. F3 or `-hz 50`
selects the comparison rate without changing the music tempo.

On Android the existing pad moves the character. ENTER opens the door, MENU
leaves a screen, RESET starts the reset demo and SCREENS opens the selector. The
60 HZ / 50 HZ button changes animation speed without restarting the scene.

```sh
./scripts/run-cuddlydemo-android.sh
```

The separate application is **Cuddly Demo (DCK)**, package
`com.olivierh.cuddlydemo`. Add `--build-only` to build without installing.

Assets live directly under `assets/cuddly/<screen>/`; common replacement tracks
are in `assets/cuddly/ym/`. `data.json` contains native screen text and movement
tables, and `loader.json` contains loading-screen captions and counters.

The effects include configurable wave chains, profile-based distortion, sparkle
overlays, feedback ribbons, projected rows and batched circular particles.
Images, fonts, phases, sizes and layer order remain production parameters.
Big Sprite colors its scrolling text with `composite.RasterOverlay` using
source-atop blending; Starwars uses the same effect with source-in blending and
a different wrap rule. Their 15-second captures match all 900 original frames
per screen.
Ehhh's seven moving raster images now use one `sprites.Train` with a
phase-spaced cosine wave. Five checkpoints match the previous Ehhh renderer
pixel for pixel.
LED's background and gradient also use DCK's `motion.WrapBank`; ten
checkpoints, including the pre-render and both wrap boundaries, match the
previous screen pixel for pixel.
Its nine letters now use `motion.HarmonicFormation` through `sprites.Group`.
The authored phase table and a bouncing amplitude remain editable DCK
parameters, and the group keeps the original draw-before-step timing.

The unattended tour visits the introduction, every menu door, Reset and the menu
again. Scene durations exclude loading and walking time.

```sh
go run ./dck/cmd/cuddlydemo -tour -screen-duration 1m
go run ./dck/cmd/video -output /path/to/cuddly-demo.mp4
```

Video export requires FFmpeg. It records the game canvas and its own audio at
60 FPS with a minute per screen, and writes chapters, a PNG poster and a JSON
timing report. `-intro-duration`, `-screen-duration` and `-menu-duration` customize
the route. No desktop pixels or other applications' audio are captured.

```sh
go test -race ./...
go vet ./...
go run ./dck/cmd/checkcuddly
go run ./dck/cmd/checktiming -hz 60 -out captures/timing.json
```
