#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."

if [ ! -f "wallet.apk" ]; then
    ./scripts/build-android.sh
fi

echo "Installing wallet.apk on connected device via adb..."
adb install -r wallet.apk
echo "✓ Installed successfully! Launching app..."
adb shell monkey -p com.shakilnwz.wallet -c android.intent.category.LAUNCHER 1
