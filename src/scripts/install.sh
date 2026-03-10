#!/usr/bin/env bash
set -euo pipefail

version="${ORB_PARAM_VERSION#v}"
destination="${ORB_PARAM_DESTINATION/#\~/$HOME}"

if [[ -z "$version" ]]; then
  echo "Fehler: Die Release-Version darf nicht leer sein." >&2
  exit 1
fi

if [[ -z "$destination" ]]; then
  echo "Fehler: Das Zielverzeichnis darf nicht leer sein." >&2
  exit 1
fi

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$os" in
  linux|darwin)
    ;;
  *)
    echo "Fehler: Nicht unterstuetztes Betriebssystem: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64)
    arch="amd64"
    ;;
  *)
    echo "Fehler: Nicht unterstuetzte Architektur: $arch" >&2
    exit 1
    ;;
esac

asset="pdfzus-merge-${os}-${arch}"
url="https://github.com/wsgtcyx/pdf-zusammenfuegen-circleci-orb/releases/download/v${version}/${asset}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

mkdir -p "$destination"

echo "Lade ${asset} aus ${url}"
curl -fsSL "$url" -o "$tmpdir/pdfzus-merge"
chmod +x "$tmpdir/pdfzus-merge"
mv "$tmpdir/pdfzus-merge" "$destination/pdfzus-merge"

echo "export PATH=\"$destination:\$PATH\"" >> "$BASH_ENV"
echo "pdfzus-merge wurde unter $destination installiert."

