# Cuddly Demos in Go

Native Go/Ebitengine implementation with YM playback and embedded artwork,
fonts and audio. Desktop and Android share the same rendering and input code.

```sh
go run ./dck/cmd/cuddlydemo   # Complete production, starting with the intro.
go run ./dck/cmd/cuddlymenu   # Start directly in the menu.
go run ./cmd/cuddlymenu       # Original standalone Go menu.
```

The complete version contains the introduction, menu, thirteen demo screens and
Reset. Loading transitions include artwork, counters, scrolling text and audio.

Playback defaults to 60 fixed updates per second. Use F3, `-hz 50`, or the mobile
rate button to compare animation speeds. Music retains its normal sample clock.
F1 selects a screen; Space/Escape returns to the menu; R opens Reset.

Build and install the complete Android application on one authorized USB device:

```sh
./scripts/run-cuddlydemo-android.sh
```

The app is **Cuddly Demo (DCK)** (`com.olivierh.cuddlydemo`). Its touch controls
provide movement/thrust, ENTER, MENU, RESET, SCREENS and the rate selector.
The original standalone Android menu uses `./scripts/run-android.sh`.

See [the native DCK entry points](dck/README.md) for screen selection, asset layout
and validation commands. Local module replacements select sibling DCK and
YM-player checkouts in the workspace.
