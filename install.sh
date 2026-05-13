#!/usr/bin/env bash
# kylang - Created and Developed by Farhan Kertadiwangsa
# Copyright (c) 2026 Farhan Kertadiwangsa. All rights reserved.

set -e

REPO_URL="https://github.com/Farhankerdiwangsa/ky-lang" 
BINARY_NAME="ky"
INSTALL_DIR="/usr/local/bin"
TMP_DIR="$(mktemp -d)"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${BOLD}Bahasa Pemrograman kylang${NC} - Created by Farhan Kertadiwangsa"

if ! command -v go &>/dev/null; then
    echo -e "${RED}[Error]${NC} Go tidak ditemukan. Silakan install Go terlebih dahulu: https://go.dev/dl/"
    exit 1
fi

echo -e "${GREEN}[1/3]${NC} Mengunduh source code..."
git clone --depth 1 "$REPO_URL" "$TMP_DIR/ky-lang-src" &>/dev/null || {
    echo -e "${RED}[Error]${NC} Gagal mengunduh source code."
    exit 1
}

echo -e "${GREEN}[2/3]${NC} Mengompilasi interpreter..."
cd "$TMP_DIR/ky-lang-src"
go build -ldflags="-s -w" -o "$BINARY_NAME" . || {
    echo -e "${RED}[Error]${NC} Gagal melakukan kompilasi."
    exit 1
}

echo -e "${GREEN}[3/3]${NC} Memasang ke sistem..."
if [ -w "$INSTALL_DIR" ]; then
    cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo -e "\n${BOLD}${GREEN}Berhasil!${NC} Bahasa kylang telah terpasang."
echo -e "Ketik '${BOLD}ky${NC}' untuk memulai."

rm -rf "$TMP_DIR"