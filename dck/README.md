# DCK version

This directory contains the construction-kit version of go-cuddlymenu. The original Go sources are preserved at their original paths (revision `75f1c185be4d999e38708bc724ee9ad4a6e4d828`), with small asset accessors so both versions use the same embedded resources.

Run the original with `go run ./cmd/cuddlymenu` and this version with `go run ./dck/cmd/cuddlymenu` from the repository root.

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
