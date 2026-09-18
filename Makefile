# Makefile for Terminal Odyssey
# Modern CLI Settlement Management and Dungeon Crawler

BINARY_NAME=odyssey
BUILD_DIR=bin
CMD_DIR=./cmd/odyssey
GO=go
LDFLAGS=-s -w

.PHONY: all build run test test-coverage vet lint clean help build-all

all: vet test build

build:
	@echo "==> Membangun binary $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "==> Selesai: $(BUILD_DIR)/$(BINARY_NAME)"

run:
	@$(GO) run $(CMD_DIR)

test:
	@echo "==> Menjalankan seluruh pengujian unit..."
	$(GO) test -v -count=1 ./...

test-coverage:
	@echo "==> Menjalankan analisis cakupan kode..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out
	@rm -f coverage.out

vet:
	@echo "==> Menjalankan pemeriksaan kode statis (go vet)..."
	$(GO) vet ./...

lint: vet

clean:
	@echo "==> Membersihkan direktori build dan file sementara..."
	@rm -rf $(BUILD_DIR) dist coverage.out
	@echo "==> Direktori build telah dibersihkan"

build-all:
	@echo "==> Melakukan cross-compilation untuk berbagai platform..."
	@mkdir -p $(BUILD_DIR)
	@echo "--> Linux (amd64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@echo "--> Linux (arm64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	@echo "--> Windows (amd64)..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "--> macOS (amd64)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	@echo "--> macOS (arm64)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	@echo "==> Seluruh binary rilis berhasil dibuat di $(BUILD_DIR)/"

help:
	@echo "Terminal Odyssey - Makefile Commands"
	@echo "====================================="
	@echo "  make build          : Kompilasi binary ke $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  make run            : Jalankan game secara langsung dengan 'go run'"
	@echo "  make test           : Jalankan seluruh unit test (count=1)"
	@echo "  make test-coverage  : Tampilkan ringkasan cakupan kode test"
	@echo "  make vet            : Analisis kode statis dengan 'go vet'"
	@echo "  make clean          : Hapus folder build dan file sementara"
	@echo "  make build-all      : Cross-compile binary untuk Linux, Windows, dan macOS"
	@echo "  make all            : Jalankan vet, test, dan build secara berurutan"
	@echo "  make help           : Tampilkan panduan perintah Makefile ini"
