#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ARTIFACTS_DIR="$ROOT_DIR/artifacts"

mkdir -p "$ARTIFACTS_DIR"
find "$ARTIFACTS_DIR" -mindepth 1 -maxdepth 1 -type f -delete

go test ./...

targets=(
  "darwin amd64 pdfzus-merge-darwin-amd64"
  "linux amd64 pdfzus-merge-linux-amd64"
  "windows amd64 pdfzus-merge-windows-amd64.exe"
)

for target in "${targets[@]}"; do
  read -r goos goarch filename <<<"$target"
  GOOS="$goos" GOARCH="$goarch" go build -o "$ARTIFACTS_DIR/$filename" ./cmd/pdfzus-merge
done

go run ./tool/make-sample-artifacts --output-dir "$ARTIFACTS_DIR"

(
  cd "$ARTIFACTS_DIR"
  rm -f SHA256SUMS
  shasum -a 256 pdfzus-merge-* > SHA256SUMS
)

echo "Artefakte wurden in $ARTIFACTS_DIR erstellt."

