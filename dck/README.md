# DCK version

This directory contains the construction-kit version of go-cuddlymenu. The original Go sources are preserved at their original paths (revision `75f1c185be4d999e38708bc724ee9ad4a6e4d828`), with small asset accessors so both versions use the same embedded resources.

From the repository root:

```sh
go run ./dck/cmd/cuddlydemo   # Complete production: Calvin, animated intro, menu.
go run ./dck/cmd/cuddlymenu   # Same controller, starting directly in the menu.
go run ./cmd/cuddlymenu       # Preserved original Go menu.
```

Press Space or the touch **MENU** button to leave the intro. Both DCK launchers
use the original picture, sector/decrunch countdown, scrolling loader text and
loader audio on transitions. The MP3 intro track is embedded and starts with the
logo animation.

The DCK menu opens all thirteen main Cuddly screens. Introduction and Reset are
also available from the selector: fifteen screens plus the existing menu. NATIVE
Credits is intentionally excluded.

```sh
go run ./dck/cmd/cuddlymenu -list
go run ./dck/cmd/cuddlymenu -screen intro
go run ./dck/cmd/cuddlymenu -screen spreadpoint
go run ./dck/cmd/cuddlymenu -screen dna
```

F1 opens the selector; arrows and Enter select a screen. Esc or Space returns to
the menu, and R opens Reset from an active screen. In Megaball, Up/Down selects one
of the nine parameters and Left/Right changes it. `-mute` disables device audio.

Playback uses **50 logical updates per second**, independently of display refresh.
Drawing never advances the animation. YM and recorded audio keep their normal
sample clock; the speed of the original implementation in a modern browser is not reproduced.

The source audit and porting guide documents each screen,
shared effects, music, timing and the 94 visual checkpoints. Runtime code is
Go/Ebitengine; the original implementation remains solely an archived visual reference.

## Android target

`./scripts/run-cuddlydemo-android.sh` builds, installs and launches the separate
**Cuddly Demo (DCK)** application (`com.olivierh.cuddlydemo`). The original Android
application remains installed independently. Add `--build-only` to build an APK
without touching a connected device.

The existing directional/thrust pad is preserved. **ENTER** enters the door in
front of the character; **MENU** exits a screen; **RESET** opens the reset demo.
**SCREENS** opens a touch-selectable list. In Megaball the right-hand buttons edit
its movement parameters. Buttons use the Pixel's side margins.

See Android build and validation and
reusable image effects.
