# Wallet — Native Android App (Pure Go & Gio UI)

A high-performance, mobile-first personal expense tracker written in **100% pure Go** using **Gio UI (`gioui.org`)**, compiled directly to Android using **`gogio`**.

> **Note:** This is the `native-android-go` branch containing the pure Go native mobile and desktop application.

---

## 🌟 Key Features

- **📱 Pure Go Mobile-First UI**:
  - Built with Gio UI (`gioui.org`) immediate-mode vector graphics.
  - Custom on-screen tactile numeric keypad (`[1] [2] [3]...`) for instant 1-thumb entry without software keyboard resize delay.
  - OLED dark theme (`#090A0F`), smooth touch fling physics, and responsive layouts.
  - Virtualized transaction ledger and category spending progress bars.

- **💾 Zero-CGO Pure-Go SQLite**:
  - Powered by `modernc.org/sqlite`, allowing effortless cross-compilation for Android (`arm64-v8a`) without needing an Android NDK CGO toolchain.
  - Full ACID compliance, WAL mode, and sub-millisecond local queries.

- **🔄 Peer-to-Peer Wi-Fi / Local Network Sync**:
  - Embedded local sync HTTP server with Last-Write-Wins (LWW) conflict resolution.
  - Every transaction is tracked via UUID v4, timestamps, and soft-delete tombstones.
  - Seamless two-way sync between your Android phone and desktop client over local Wi-Fi without any third-party cloud.

- **📦 Unified Codebase (Android & Desktop)**:
  - Zero Java, Kotlin, or Gradle required.
  - The exact same Go code runs as a native Android APK and as a Desktop application (Linux/macOS/Windows).

---

## 🚀 Quick Start

### 1. Run on Desktop (Linux / macOS / Windows)
For instant local testing on your computer:
```bash
go run ./cmd/wallet
```

### 2. Build for Android (`.apk`)
Build the release-ready Android APK using `gogio`:
```bash
./scripts/build-android.sh
```
This generates `wallet.apk` in the project root.

### 3. Install on Connected Phone (via ADB)
Plug in your Android phone via USB with USB Debugging enabled:
```bash
./scripts/install-phone.sh
```

---

## 📦 Releases & CI/CD

### Automatic GitHub Release (Recommended)
Every time a git release tag (e.g., `v0.1.0`) is pushed, GitHub Actions automatically builds and publishes:
- **`wallet-android-arm64.apk`**: Release APK for Android devices
- **`wallet-linux-amd64.tar.gz`**: Standalone Linux desktop binary
- **`wallet-windows-amd64.zip`**: Standalone Windows desktop binary

To trigger a release:
```bash
git tag v0.1.0
git push origin v0.1.0
```

### Local Release Build
Build release desktop archives on your local machine:
```bash
./scripts/build-release.sh v0.1.0
```
Artifacts are generated in the `dist/` directory.

---

## 🧪 Testing

Run all unit tests for the SQLite engine, models, and P2P sync server:
```bash
go test -v ./...
```
