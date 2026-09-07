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

Le document [Adapter un projet Go + Ebitengine vers Android / Pixel](docs/GUIDE_ANDROID_EBITENGINE_PIXEL.md) décrit en détail l’installation des outils, l’architecture du portage, le tactile, le build, l’installation USB et le diagnostic des problèmes rencontrés.
