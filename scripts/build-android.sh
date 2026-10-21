#!/usr/bin/env bash
set -e

# Ensure we run from repository root
cd "$(dirname "$0")/.."

echo "==> Building Wallet Android APK using gogio..."

# Locate gogio tool (manually installed by user)
GOGIO_BIN=""
if command -v gogio >/dev/null 2>&1; then
    GOGIO_BIN="gogio"
elif [ -n "$(go env GOBIN 2>/dev/null)" ] && [ -x "$(go env GOBIN)/gogio" ]; then
    GOGIO_BIN="$(go env GOBIN)/gogio"
elif [ -n "$(go env GOPATH 2>/dev/null)" ] && [ -x "$(go env GOPATH)/bin/gogio" ]; then
    GOGIO_BIN="$(go env GOPATH)/bin/gogio"
elif [ -x "$HOME/go/bin/gogio" ]; then
    GOGIO_BIN="$HOME/go/bin/gogio"
fi

if [ -z "$GOGIO_BIN" ]; then
    echo "Error: 'gogio' tool not found in PATH or ~/go/bin."
    echo "Please install gogio manually using:"
    echo "  go install gioui.org/cmd/gogio@latest"
    echo "and ensure \$(go env GOPATH)/bin is in your PATH."
    exit 1
fi

echo "Using gogio: $GOGIO_BIN"

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
"$GOGIO_BIN" -target android \
    -arch arm64 \
    -appid com.shakilnwz.wallet \
    -version 1.0.0.1 \
    -minsdk 24 \
    -targetsdk 34 \
    -o wallet.apk \
    ./cmd/wallet

echo "✓ Successfully built wallet.apk!"
ls -lh wallet.apk
