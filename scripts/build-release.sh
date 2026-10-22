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
[ -n "$DIST_DIR" ] || exit 1
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

# 3. Build Android APK (if the Android toolchain is available)
if [ -n "$ANDROID_HOME" ]; then
    echo "--> Building Android APK (arm64)..."
    ./scripts/build-android.sh "$DIST_DIR/wallet-android-arm64.apk"
else
    echo "--> Skipping local Android build (ANDROID_HOME not set; will be built by GitHub Actions on release tag)."
fi

echo ""
echo "✓ Release artifacts generated in dist/:"
ls -lh "$DIST_DIR"/*.tar.gz "$DIST_DIR"/*.zip "$DIST_DIR"/*.apk 2>/dev/null || ls -lh "$DIST_DIR"
