#!/bin/sh
set -eu

project_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
android_sdk=${ANDROID_HOME:-${ANDROID_SDK_ROOT:-/opt/homebrew/share/android-commandlinetools}}
java_path=${JAVA_HOME:-/opt/homebrew/opt/openjdk@17}
export ANDROID_HOME="$android_sdk" ANDROID_SDK_ROOT="$android_sdk" JAVA_HOME="$java_path"
export PATH="$java_path/bin:$android_sdk/platform-tools:$PATH"

if [ ! -x "$android_sdk/platform-tools/adb" ] || [ ! -x "$java_path/bin/java" ]; then
    echo "Set ANDROID_HOME and JAVA_HOME to the installed Android SDK and JDK." >&2
    exit 1
fi
mkdir -p "$project_root/android/cuddlydemo/libs"
cd "$project_root"
go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@v2.9.11 bind \
    -target android/arm64 -androidapi 23 -javapkg com.olivierh.cuddlydemo \
    -o android/cuddlydemo/libs/cuddlydemo.aar ./dck/mobile/cuddlydemo
"$project_root/android/gradlew" -p "$project_root/android" --console=plain --offline :cuddlydemo:assembleDebug
if [ "${1:-}" = "--build-only" ]; then exit 0; fi

device_count=$(adb devices | awk 'NR > 1 && $2 == "device" { count++ } END { print count + 0 }')
if [ "$device_count" -ne 1 ]; then
    echo "Exactly one authorized Android device is required; found $device_count." >&2
    adb devices -l >&2
    exit 1
fi
adb install -r "$project_root/android/cuddlydemo/build/outputs/apk/debug/cuddlydemo-debug.apk"
adb shell am start -S -W -n com.olivierh.cuddlydemo/.MainActivity
