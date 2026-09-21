# Cuddly Demo / DCK

The complete production uses shared Go/Ebitengine effects and embedded assets.

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

```sh
go test -race ./...
go vet ./...
go run ./dck/cmd/checkcuddly
go run ./dck/cmd/checktiming -hz 60 -out captures/timing.json
```
