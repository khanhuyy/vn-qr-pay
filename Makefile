# qrpay — Vietnam QR Pay (VietQR / VNPAYQR) library + CLI
#
# Common targets:
#   make build     -> compile the CLI binary into ./bin/qrpay
#   make install   -> go install the CLI under $GOPATH/bin
#   make test      -> run the full test suite
#   make cover     -> run tests with coverage profile -> coverage.out
#   make lint      -> go vet over the whole module
#   make tidy      -> go mod tidy
#   make demo      -> build, then run a small end-to-end demo
#   make clean     -> remove ./bin

GO       ?= go
BIN_DIR  ?= bin
BINARY   ?= qrpay
PKG_CLI  := ./cmd/qrpay
PKG_ALL  := ./...

.PHONY: all
all: build

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY) $(PKG_CLI)
	@echo "built -> $(BIN_DIR)/$(BINARY)"

.PHONY: install
install:
	$(GO) install $(PKG_CLI)

.PHONY: test
test:
	$(GO) test -count=1 $(PKG_ALL)

.PHONY: cover
cover:
	$(GO) test -count=1 -coverprofile=coverage.out -coverpkg=./serve/...,./constants/... $(PKG_ALL)
	$(GO) tool cover -func=coverage.out | tail -n 20

.PHONY: lint
lint:
	$(GO) vet $(PKG_ALL)

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: demo
demo: build
	@echo "== build VietQR (Vietcombank 1234567890, 50 000 VND, 'Cam on') =="
	./$(BIN_DIR)/$(BINARY) build --bin 970436 --account 1234567890 --amount 50000 --purpose "Cam on"
	@echo
	@echo "== parse ACB static QR =="
	./$(BIN_DIR)/$(BINARY) parse 00020101021138530010A0000007270123000697041601092576788590208QRIBFTTA53037045802VN6304AE9F
	@echo
	@echo "== filter banks containing 'techcom' =="
	./$(BIN_DIR)/$(BINARY) banks --filter techcom

.PHONY: clean
clean:
	rm -rf $(BIN_DIR) coverage.out
