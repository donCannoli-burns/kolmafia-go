#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p web bin
# The dashboard and native server use the same go-app component source.
GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/workbench
go build -o bin/kol-workbench ./cmd/workbench
echo "Built bin/kol-workbench + web/app.wasm"
