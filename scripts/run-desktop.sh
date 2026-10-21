#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."

echo "==> Launching Wallet on Desktop..."
export CGO_CFLAGS="-I$(pwd)/include"
go run ./cmd/wallet
