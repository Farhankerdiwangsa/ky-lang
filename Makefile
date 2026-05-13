#  Created and Developed by Farhan Kertadiwangsa
# Makefile untuk build, instalasi, dan penghapusan interpreter

BINARY    := ky
BUILD_DIR := .
INSTALL   := /usr/local/bin
GO        := go
LDFLAGS   := -ldflags="-s -w"

.PHONY: all build install uninstall clean deb help

## Default: build interpreter
all: build

## Build binary ke direktori saat ini
build:
        @echo "[ky] Mengompilasi interpreter..."
        @$(GO) build $(LDFLAGS) -o $(BINARY) $(BUILD_DIR)
        @echo "[ky] Selesai → ./$(BINARY)"

## Install binary ke $(INSTALL), atau ~/.local/bin jika tidak ada akses
install: build
        @TARGET="$(INSTALL)"; \
        if [ ! -w "$$TARGET" ]; then \
                if command -v sudo >/dev/null 2>&1; then \
                        echo "[ky] Memasang ke $$TARGET (butuh sudo)..."; \
                        sudo cp $(BINARY) $$TARGET/$(BINARY) && sudo chmod +x $$TARGET/$(BINARY); \
                else \
                        TARGET="$$HOME/.local/bin"; \
                        mkdir -p "$$TARGET"; \
                        echo "[ky] Memasang ke $$TARGET (tanpa sudo)..."; \
                        cp $(BINARY) $$TARGET/$(BINARY); \
                        chmod +x $$TARGET/$(BINARY); \
                        echo "[ky] Tambahkan ke PATH jika belum ada:"; \
                        echo "       export PATH=\"\$$HOME/.local/bin:\$$PATH\""; \
                fi; \
        else \
                echo "[ky] Memasang ke $$TARGET..."; \
                cp $(BINARY) $$TARGET/$(BINARY); \
                chmod +x $$TARGET/$(BINARY); \
        fi
        @echo "[ky] Berhasil! Jalankan: ky"

## Hapus binary dari sistem
uninstall:
        @if [ -f "$(INSTALL)/$(BINARY)" ]; then \
                echo "[ky] Menghapus $(INSTALL)/$(BINARY)..."; \
                if [ -w "$(INSTALL)" ]; then rm -f $(INSTALL)/$(BINARY); \
                else sudo rm -f $(INSTALL)/$(BINARY); fi; \
        elif [ -f "$$HOME/.local/bin/$(BINARY)" ]; then \
                echo "[ky] Menghapus $$HOME/.local/bin/$(BINARY)..."; \
                rm -f "$$HOME/.local/bin/$(BINARY)"; \
        else \
                echo "[ky] ky tidak ditemukan di sistem."; \
        fi
        @echo "[ky] Terhapus."

## Build lalu buat paket .deb (membutuhkan script build-deb.sh di root)
deb: build
        @cd .. && bash scripts/build-deb.sh

## Hapus binary hasil build lokal
clean:
        @rm -f $(BINARY)
        @echo "[ky] Dibersihkan."

## Tampilkan bantuan ini
help:
        @echo ""
        @echo "  Bahasa Pemrograman kylang"
        @echo "  Diciptakan oleh Farhan Kertadiwangsa"
        @echo ""
        @echo "  Perintah yang tersedia:"
        @echo "    make            — build interpreter"
        @echo "    make install    — build + pasang ke /usr/local/bin"
        @echo "    make uninstall  — hapus dari sistem"
        @echo "    make deb        — buat paket .deb"
        @echo "    make clean      — hapus file build lokal"
        @echo ""
