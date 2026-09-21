# Cuddly Demos — menu

La démo tourne sur ordinateur et sur Android. Sur téléphone, elle utilise un pad tactile gauche/droite et un bouton **VOLER**. Les commandes acceptent le multitouch : on peut donc avancer et voler en même temps.

## Lancer sur le Pixel en USB

Prérequis : Go 1.25 ou plus récent, Java 17, le SDK/NDK Android, et le débogage USB activé sur le téléphone.

```sh
./scripts/run-android.sh
```

Le script :

1. génère la bibliothèque Android arm64 avec Ebitengine ;
2. compile `android/app/build/outputs/apk/debug/app-debug.apk` ;
3. installe l’APK sur l’unique appareil USB autorisé ;
4. lance l’application en plein écran paysage.

Au premier branchement, déverrouillez le Pixel et acceptez la demande d’autorisation RSA. Si le SDK ou Java ne sont pas installés aux emplacements Homebrew habituels, définissez `ANDROID_HOME` et `JAVA_HOME` avant de lancer le script.

## Lancer sur ordinateur

```sh
go run ./cmd/cuddlymenu
```

Les contrôles clavier existants restent disponibles : flèches gauche/droite, flèche haut ou Entrée pour voler, Espace devant une porte, `R` pour recommencer et `C` pour l’effet CRT.

## Guide de portage réutilisable

Le document Adapter un projet Go + Ebitengine vers Android / Pixel décrit en détail l’installation des outils, l’architecture du portage, le tactile, le build, l’installation USB et le diagnostic des problèmes rencontrés.

La migration du synthétiseur, la suppression des allocations dans le flux PCM et le passage à 48 kHz sont documentés dans [GUIDE_MIGRATION_YM_PLAYER_48KHZ.md](GUIDE_MIGRATION_YM_PLAYER_48KHZ.md).

## Optional DCK version

The original implementation remains at its original paths. Run it with `go run ./cmd/cuddlymenu`.

The construction-kit version is in [dck/](dck/README.md). Run
`go run ./dck/cmd/cuddlydemo` for the complete introduction and menu, or
`go run ./dck/cmd/cuddlymenu` to start directly in the menu. Both DCK entry points
use the original image/countdown loading transition. Original assets are shared.

The separate Android target is built and installed with
`./scripts/run-cuddlydemo-android.sh`. Its launcher name is **Cuddly Demo (DCK)**;
see the Android guide.

The DCK version connects the menu to all thirteen main Cuddly screens, plus
Introduction and Reset. NATIVE Credits is excluded. Playback runs at a fixed
cadence of 60 updates per second by default, independent of monitor refresh, with
normal music playback speed. Use `-hz 50`, F3 or the mobile rate button to compare. F1 opens the selector; Esc or Space returns from a screen;
R opens Reset. `-list` lists the ports and `-screen dna` launches one directly.

The Cuddly porting guide covers effect reuse, controls,
media provenance and visual/timing validation. Original desktop/Android commands
continue to launch the preserved original menu.
