#!/usr/bin/env bash
set -e

# Ensure we run from repository root
cd "$(dirname "$0")/.."

VERSION="${1:-}"
if [ -z "$VERSION" ]; then
    # Try reading the latest git tag, or fallback to dev
    VERSION=$(git describe --tags --always 2>/dev/null || echo "v1.0.0")
fi

echo "==> Building Wallet release artifacts ($VERSION)..."

DIST_DIR="$(pwd)/dist"
mkdir -p "$DIST_DIR"
rm -f "$DIST_DIR"/*

# 1. Build Linux Desktop (amd64)
echo "--> Building Linux Desktop (x86_64)..."
CGO_CFLAGS="-I$(pwd)/include" go build -ldflags="-s -w" -o "$DIST_DIR/wallet-linux-amd64" ./cmd/wallet
tar -czf "$DIST_DIR/wallet-linux-amd64.tar.gz" -C "$DIST_DIR" wallet-linux-amd64

# 2. Build Windows Desktop (amd64)
echo "--> Building Windows Desktop (x86_64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o "$DIST_DIR/wallet-windows-amd64.exe" ./cmd/wallet
zip -j -q "$DIST_DIR/wallet-windows-amd64.zip" "$DIST_DIR/wallet-windows-amd64.exe"

# 3. Build Android APK (if ANDROID_HOME is set)
GOGIO_BIN=""
if [ -x "./bin/gogio" ]; then
    GOGIO_BIN="./bin/gogio"
elif command -v gogio >/dev/null 2>&1; then
    GOGIO_BIN="gogio"
elif [ -n "$(go env GOBIN 2>/dev/null)" ] && [ -x "$(go env GOBIN)/gogio" ]; then
    GOGIO_BIN="$(go env GOBIN)/gogio"
elif [ -n "$(go env GOPATH 2>/dev/null)" ] && [ -x "$(go env GOPATH)/bin/gogio" ]; then
    GOGIO_BIN="$(go env GOPATH)/bin/gogio"
fi

if [ -n "$ANDROID_HOME" ] && [ -n "$GOGIO_BIN" ]; then
    echo "--> Building Android APK (arm64)..."
    "$GOGIO_BIN" -target android \
        -arch arm64 \
        -appid com.shakilnwz.wallet \
        -version 1.0.0.1 \
        -minsdk 24 \
        -targetsdk 34 \
        -o "$DIST_DIR/wallet-android-arm64.apk" \
        ./cmd/wallet
else
    echo "--> Skipping local Android build (ANDROID_HOME not set; will be built by GitHub Actions on release tag)."
fi

echo ""
echo "✓ Release artifacts generated in dist/:"
ls -lh "$DIST_DIR"/*.tar.gz "$DIST_DIR"/*.zip "$DIST_DIR"/*.apk 2>/dev/null || ls -lh "$DIST_DIR"
