#!/usr/bin/env bash
set -euo pipefail
files=$(rg --files -g '*.php' || true)
count=$(printf "%s" "$files" | sed '/^$/d' | wc -l | tr -d ' ')
if [ "$count" != "0" ]; then
  echo "found $count PHP files"
  printf "%s\n" "$files"
  exit 1
fi

echo "OK: no PHP files found"
