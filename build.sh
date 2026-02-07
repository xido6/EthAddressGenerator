#!/bin/bash

OLD_GOOS=$(go env GOOS)
OLD_GOARCH=$(go env GOARCH)

SRC_DIR="."
PROGRAM_NAME="generator"

# macOS Intel
echo "🔹 Building macOS (Intel)..."
go env -w GOOS=darwin GOARCH=amd64
mkdir -p ./resource/macOS/amd64
go build -o ./resource/macOS/amd64/${PROGRAM_NAME} "$SRC_DIR"

# macOS ARM (M1/M2/M3)
echo "🔹 Building macOS (ARM)..."
go env -w GOOS=darwin GOARCH=arm64
mkdir -p ./resource/macOS/arm64
go build -o ./resource/macOS/arm64/${PROGRAM_NAME} "$SRC_DIR"

# Windows 64-bit
echo "🔹 Building Windows 64-bit..."
go env -w GOOS=windows GOARCH=amd64
mkdir -p ./resource/windows_amd64
go build -o ./resource/windows_amd64/${PROGRAM_NAME}.exe "$SRC_DIR"

echo "🔹 Restoring original Go environment..."
go env -w GOOS=$OLD_GOOS GOARCH=$OLD_GOARCH

echo "✅ All builds completed!"