# DCK version

This directory contains the construction-kit version of go-cuddlymenu. The original Go sources are preserved at their original paths (revision `75f1c185be4d999e38708bc724ee9ad4a6e4d828`), with small asset accessors so both versions use the same embedded resources.

Run the original with `go run ./cmd/cuddlymenu` and this version with `go run ./dck/cmd/cuddlymenu` from the repository root.

The DCK menu now opens eight native Cuddly screens: Big Sprite, Colorshock II,
Ehhh, Mega Scroller, Digi Sound, LED Scroller, Fullscreen and Knucklebuster.
Press F1 for selection and Esc or Space to return from a screen.

```sh
go run ./dck/cmd/cuddlymenu -list
go run ./dck/cmd/cuddlymenu -screen knucklebuster
```

The complete source audit and porting guide describes all
17 entries, downloaded assets, selected YM files and the remaining ports. Runtime
code is Go/Ebitengine; the original implementation is retained solely as a local visual reference.
