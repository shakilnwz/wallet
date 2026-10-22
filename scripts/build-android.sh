#!/usr/bin/env bash
set -e

# Ensure we run from repository root
cd "$(dirname "$0")/.."

echo "==> Building Wallet Android APK using gogio..."

# Locate gogio tool (project-local bin or environment)
GOGIO_BIN=""
if [ -x "./bin/gogio" ]; then
    GOGIO_BIN="./bin/gogio"
elif command -v gogio >/dev/null 2>&1; then
    GOGIO_BIN="gogio"
elif [ -n "$(go env GOBIN 2>/dev/null)" ] && [ -x "$(go env GOBIN)/gogio" ]; then
    GOGIO_BIN="$(go env GOBIN)/gogio"
elif [ -n "$(go env GOPATH 2>/dev/null)" ] && [ -x "$(go env GOPATH)/bin/gogio" ]; then
    GOGIO_BIN="$(go env GOPATH)/bin/gogio"
elif [ -x "$HOME/go/bin/gogio" ]; then
    GOGIO_BIN="$HOME/go/bin/gogio"
fi

if [ -z "$GOGIO_BIN" ]; then
    echo "Error: 'gogio' tool not found."
    echo "Please install gogio either project-locally:"
    echo "  mkdir -p bin && GOBIN=\$(pwd)/bin go install gioui.org/cmd/gogio@latest"
    echo "or globally:"
    echo "  go install gioui.org/cmd/gogio@latest"
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

OUT="${1:-wallet.apk}"

# Version from the latest tag: v1.2.3 -> 1.2.3, versionCode 10203 (matches CI).
TAG="$(git describe --tags --abbrev=0 2>/dev/null || echo v0.0.0)"
VERSION="${TAG#v}"
IFS=. read -r MAJOR MINOR PATCH <<< "$VERSION"
CODE=$(( MAJOR * 10000 + MINOR * 100 + PATCH ))

# Sign with the release key when available so local builds install over
# GitHub releases (and vice versa). Override with WALLET_KEYSTORE/_PASS_FILE.
KEYSTORE="${WALLET_KEYSTORE:-$HOME/.local/share/wallet-signing/wallet-release.p12}"
PASS_FILE="${WALLET_KEYSTORE_PASS_FILE:-$HOME/.local/share/wallet-signing/wallet-release.pass}"
SIGN_ARGS=()
if [ -f "$KEYSTORE" ] && [ -f "$PASS_FILE" ]; then
    echo "Signing with release key: $KEYSTORE"
    SIGN_ARGS=(-signkey "$KEYSTORE" -signpass "$(cat "$PASS_FILE")")
else
    echo "Warning: release keystore not found; using a debug key. This APK won't install over release builds."
fi

# Build APK for arm64 (Android 64-bit phones)
"$GOGIO_BIN" -target android \
    -arch arm64 \
    -appid com.shakilnwz.wallet \
    -version "$VERSION.$CODE" \
    -minsdk 24 \
    -targetsdk 34 \
    "${SIGN_ARGS[@]}" \
    -o "$OUT" \
    ./cmd/wallet

echo "✓ Successfully built $OUT ($VERSION, versionCode $CODE)"
ls -lh "$OUT"
