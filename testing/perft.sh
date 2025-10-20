#!/bin/bash

# Intended to be used with https://github.com/agausmann/perftree (Incredible tool)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ "$#" -ne 2 ] && [ "$#" -ne 3 ]; then
  echo "Usage: $0 <depth> <fen> [moves]"
  exit 1
fi

if [ "$#" -eq 2 ]; then
  go run "$SCRIPT_DIR/../cmd/chess.go" perft "$1" "$2"
else
  go run "$SCRIPT_DIR/../cmd/chess.go" perft "$1" "$2" "$3"
fi
