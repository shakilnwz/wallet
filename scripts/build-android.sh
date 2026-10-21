#!/usr/bin/env bash
set -e

# Ensure we run from repository root
cd "$(dirname "$0")/.."

echo "==> Building Wallet Android APK using gogio..."

# Ensure gogio binary is available
if [ ! -f "./bin/gogio" ]; then
    echo "gogio not found in ./bin, installing..."
    mkdir -p bin
    GOBIN="$(pwd)/bin" go install gioui.org/cmd/gogio@latest
fi

# Detect ANDROID_HOME if not set
if [ -z "$ANDROID_HOME" ]; then
    if [ -d "$HOME/Android/Sdk" ]; then
        export ANDROID_HOME="$HOME/Android/Sdk"
    elif [ -d "/opt/android-sdk" ]; then
        export ANDROID_HOME="/opt/android-sdk"
    elif [ -d "$HOME/.android-sdk" ]; then
        export ANDROID_HOME="$HOME/.android-sdk"
    fi
fi

if [ -z "$ANDROID_HOME" ]; then
    echo "Error: ANDROID_HOME environment variable is not set."
    echo "Please install the Android SDK and export ANDROID_HOME, for example:"
    echo "  export ANDROID_HOME=\$HOME/Android/Sdk"
    exit 1
fi

echo "Using ANDROID_HOME: $ANDROID_HOME"

# Build APK for arm64 (Android 64-bit phones)
./bin/gogio -target android \
    -arch arm64 \
    -appid com.shakilnwz.wallet \
    -version 1.0.0.1 \
    -minsdk 24 \
    -targetsdk 34 \
    -o wallet.apk \
    ./cmd/wallet

echo "✓ Successfully built wallet.apk!"
ls -lh wallet.apk
