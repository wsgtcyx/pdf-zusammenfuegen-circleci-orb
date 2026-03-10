#!/usr/bin/env bash
set -euo pipefail

workdir="${ORB_PARAM_WORKDIR:-.}"
output="${ORB_PARAM_OUTPUT:-}"
inputs_raw="${ORB_PARAM_INPUTS:-}"
bookmarks="${ORB_PARAM_BOOKMARKS:-false}"
divider="${ORB_PARAM_DIVIDER:-false}"
optimize="${ORB_PARAM_OPTIMIZE:-true}"

if [[ -z "$output" ]]; then
  echo "Fehler: Der Ausgabe-Pfad darf nicht leer sein." >&2
  exit 1
fi

if [[ -n "$workdir" ]]; then
  cd "$workdir"
fi

declare -a inputs
tmp_inputs_file="$(mktemp)"
trap 'rm -f "$tmp_inputs_file"' EXIT
printf '%s' "$inputs_raw" > "$tmp_inputs_file"

while IFS= read -r line || [[ -n "$line" ]]; do
  [[ -z "$line" ]] && continue
  inputs+=("$line")
done < "$tmp_inputs_file"

if [[ "${#inputs[@]}" -lt 2 ]]; then
  echo "Fehler: Bitte mindestens zwei PDF-Dateien ueber den Parameter inputs angeben." >&2
  exit 1
fi

declare -a flags
if [[ "$bookmarks" == "true" ]]; then
  flags+=("--bookmarks")
fi
if [[ "$divider" == "true" ]]; then
  flags+=("--divider")
fi
if [[ "$optimize" != "true" ]]; then
  flags+=("--no-optimize")
fi

mkdir -p "$(dirname "$output")"

echo "Fuehre Merge fuer ${#inputs[@]} PDF-Dateien aus."
pdfzus-merge "${flags[@]}" -o "$output" "${inputs[@]}"
