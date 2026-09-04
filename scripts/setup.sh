#!/bin/sh

set -eu

if ! command -v go >/dev/null 2>&1; then
	printf '%s\n' "Error: Go is required but was not found in PATH." >&2
	exit 1
fi

cd "$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"

echo "Downloading dependencies..."
go mod download

echo "Running tests..."
go test ./...

echo "Building purrpeek..."
mkdir -p bin
go build -o bin/purrpeek ./cmd/purrpeek

echo "Ready: bin/purrpeek"
