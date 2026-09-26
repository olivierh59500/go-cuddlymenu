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
Spreadpoint's twenty balls now use `sprites.Group` with a compiled X/Y formula.
The time and index divisors, sine frequencies, amplitudes and optional pixel
snapping are editable; a pure DCK test compares all positions over 400 ticks.
Megaball's two interleaved nineteen-ball series also use `sprites.Group`.
`motion.CoupledOrbitFormation` retains the exact per-ball phase stepping and
range order, while the nine original controls edit its owned orbit each tick.
Spreadpoint and Fullscreen now compose their raster-filled logos with
`composite.SurfaceLayer` using their existing 128×128 and 768×52 surfaces.
Spreadpoint edits the two logo passes' position and scale each tick and keeps
the authored first-frame filter change.
Cuddly DNA's front/back text ribbon now uses `composite.TwistingRibbon` for its
twenty ordered strips and strict phase wrap. A pure controller test compares
every strip pose for 1,000 frames; the logo row warp and scene layers retain
their independent timing.
Its three logo copies now use two `composite.SampledRows` programs with the
same absolute scene clock. Their source-row crops, bounce, zoom and draw order
are configurable DCK effects.
Digi's 170-row logo now uses `composite.TableWarpLogo`, which also supplies the
bounce position for its companion Union logo without a screen-local phase loop.
Ehhh's roller now uses `motion.CuedWaveClock` for lookahead characters, lane
visibility, phase speeds and its stop at the next landing boundary.
Its logo now uses the same sampled-row effect as Digi with an independent curve;
its inner strip and middle-text raster use `WrapBank` and `RasterOverlay`.
Mega Scroller's bar mask and 3D DOC's ordered inner-text color passes now use
configurable `composite.RasterOverlay` source-atop materials.
Mega Scroller and LED Scroller now share `composite.TiledWaveBackdrop` for their
cached tile fields and row waves. LED selects a retained moving source and
keeps its explicit one-frame preload; Mega Scroller draws its two waves directly.
Their source sizes and output surfaces are unchanged.
LED Scroller's 11-color, 2,000-row raster now comes from DCK's configurable
uniform gradient material with the original center sampling and rounding.
The 125 orange discs now use `sprites.RotatingDiscCloud` for their shared
rotation, projection, stable depth order and batched material. A pure test
matches the prior model poses over 1,000 ticks; screen composition stays here.
Knucklebuster's head and drums now use `sprites.LatchedOverlay` for its seeded
five-tick hit windows. A separate DCK preset can replace the random source
with a live music/event signal without changing the sprite positions.
Starwars' eight small sprites now use `sprites.SampledSpriteTrain` with its
authored XY tables, staggered samples and independent X/Y waves. The shared
controller keeps one copy of each path table and retains the strict wrap tick.
Its green and red text now uses `scrolling.DualProfiledRing`: both fonts, the
segmented vertical profile, raster fill and first-frame filter are DCK settings.
The recurring background flash uses `modulation.PeriodicDecay`, including its
same-tick fade and optional live-trigger configuration.
Its projected stars now use a DCK camera-velocity and depth-shade preset,
without a screen-local rotation or depth counter.
Big Sprite colors its scrolling text with `composite.RasterOverlay` using
source-atop blending; Starwars uses the same effect with source-in blending and
a different wrap rule. Their 15-second captures match all 900 original frames
per screen.
Ehhh's seven moving raster images now use one `sprites.Train` with a
phase-spaced cosine wave. Five checkpoints match the previous Ehhh renderer
pixel for pixel.
LED's gradient uses `motion.WrapBank`, while the tile-source clock now belongs
to `TiledWaveBackdrop`. Ten earlier pixel checkpoints covered the pre-render
and both wraps; the shared backdrop has an opt-in GPU comparison ready for a
graphical session.
Its nine letters now use `motion.HarmonicFormation` through `sprites.Group`.
The authored phase table and a bouncing amplitude remain editable DCK
parameters, and the group keeps the original draw-before-step timing.
Big Sprite and Fullscreen also use `sprites.Group` for their phase-spaced
`motion.Weave` letters. Digi's group now owns its phase step. Fifteen captures
across the three screens match their previous rendered frames pixel for pixel.
Mega Scroller's masked text surface uses a directional `motion.BounceBank` with
an offscreen entrance; eleven captures around its entry and rebound match the
previous frames pixel for pixel.
Fullscreen's seven recycled bitmap text lanes now use the shared
`scrolling.Config.RingLanes` transport, including its 128-tick vertical cadence.
Big Sprite uses the same ring-bank controller with two independent fonts and
draws each lane into its own mask. Twelve captures match the previous images;
the lane-wrap recurrence is also checked over 6,000 updates in DCK.
The paired inner and outer fonts of 3D DOC now use `composite.RowWarp` with
fractional source-column sampling. Seven captures across intro and main screen
match the previous images pixel for pixel; the checkerboard and ball train keep
their independent DCK controllers.
Its `timeline.IntroHandoff` retains the last intro frame and starts music on
the first main tick. Eleven captures around that boundary and during the main
screen match the previous renderer pixel for pixel.
The loader's sector/blipp countdown, overlapping fade and black hold now use
`timeline.Countdown` and `timeline.CueClock`. Rate changes keep elapsed fade
time. Twelve captures across the entry, fade and blank boundaries match the
previous loader pixel for pixel.
Reset now uses one DCK `timeline.StageSequence` for same-tick text-cursor
handoffs and a three-window showcase with two reusable fades. Its raster
images and layer order remain scene data; eight captures through the final
part match the previous renderer pixel for pixel.
The menu now borrows its map, character, logo and text tiles from
`sprites.Atlas`, while `sprites.FrameSequence` selects the authored movement
and thrust frames. Its seven sprite trajectories now use the editable DCK
`sprites.FormationCarousel`: ordered sine/cosine formulas, viewport-relative
radii, per-sprite spacing, timed slides and pixel-snapped anchors are compiled
once. The local `sine.go` formulas have been removed from this DCK version.
The menu backdrop now uses `composite.CachedTileParallax`: one 800 × 432
unmanaged surface is built at startup, and its wrapped half-speed camera
offset is drawn once per frame. The former local tiling and offset functions
are gone from the DCK menu.
The character/map camera now uses `motion.CameraFollow` with independent world
bounds, viewport size and sprite anchor. Its previous three-branch X/Y camera
calculation is removed from the DCK menu.
The rectified `motion.WaveClock` now drives the Digi logo/text, LED logo/text,
Megaball text and Ehhh roller. Ehhh keeps its text-controlled tempo and cycle
reset. Starwars uses the same `motion.Wave` form in its sampled text profile.
Thirty-six checkpoints across the five screens match the previous images pixel
for pixel, including Ehhh's later control changes.
Big Sprite's front/back emblem now uses `sprites.AxisFlip` with the original
threshold, signed scale and face angle. Sixteen checkpoints around the face
changes and bounce boundaries remain pixel-identical after snapping each
face's odd-sized anchor to its original integer center.
The Digi, Ehhh and introduction row-profile tables now come from editable DCK
wave programs. Their overlapping writes, zero lead-ins and source-phase offsets
are compiled once at setup. Thirty checkpoints, including late wraps,
remain pixel-identical.
Colorshock II now uses a DCK formula for its two-frequency backdrop orbit and
`motion.WrapBank` for the indexed scroll position. Ten checkpoints around the
strict table wrap remain pixel-identical.

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
