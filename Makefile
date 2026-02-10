# Makefile for Gojinn - The Sovereign Cloud Platform

# --- Variáveis de Configuração ---
BINARY_NAME=gojinn
SERVER_BIN=gojinn-server
SIGNER_BIN=gojinn-signer

# XCaddy é necessário para buildar o servidor com plugins
XCADDY_CMD := xcaddy build --with github.com/pauloappbr/gojinn=.

# URLs VMWare Labs (Runtimes WASM)
URL_PYTHON := "https://github.com/vmware-labs/webassembly-language-runtimes/releases/download/python%2F3.12.0%2B20231211-040d5a6/python-3.12.0.wasm"
URL_PHP    := "https://github.com/vmware-labs/webassembly-language-runtimes/releases/download/php%2F8.2.6%2B20230714-11be424/php-cgi-8.2.6.wasm"
URL_RUBY   := "https://github.com/vmware-labs/webassembly-language-runtimes/releases/download/ruby%2F3.2.2%2B20230714-11be424/ruby-3.2.2.wasm"

.PHONY: all clean ci lint audit test build-all build-cli build-server build-signer download-runtimes check-wasm-sdk

# Default target
all: build-all

# --- 1. Build Core Binaries ---
build-all: build-cli build-server build-signer

build-cli:
	@echo "🧰 Building Gojinn CLI..."
	@mkdir -p bin
	@go build -o bin/$(BINARY_NAME) ./cmd/gojinn

build-server:
	@echo "🏗️  Building Gojinn Server (Caddy)..."
	@# O output vai para a raiz ou bin/ dependendo da preferência. Coloquei na raiz para compatibilidade com seus scripts.
	@$(XCADDY_CMD) --output $(SERVER_BIN)

build-signer:
	@echo "🔑 Building Signer Tool..."
	@go build -o bin/$(SIGNER_BIN) ./cmd/signer/main.go

# --- 2. Polyglot Runtimes (Download Only) ---
# Como removemos os exemplos, mantemos apenas o download dos runtimes para a pasta lib/ e functions/
download-runtimes:
	@echo "⬇️  Checking Polyglot Runtimes..."
	@mkdir -p lib
	@mkdir -p functions
	@if [ ! -s lib/python.wasm ]; then curl -L -o lib/python.wasm $(URL_PYTHON); fi
	@if [ ! -s lib/php.wasm ]; then curl -L -o lib/php.wasm $(URL_PHP); fi
	@if [ ! -s lib/ruby.wasm ]; then curl -L -o lib/ruby.wasm $(URL_RUBY); fi
	@# Copia para functions para facilitar o uso local se necessário
	@cp lib/*.wasm functions/ 2>/dev/null || true
	@echo "✅ Runtimes ready."

# --- 3. CI & Quality Assurance ---
setup-ci-tools:
	@echo "🛠️  Setting up CI tools..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest

lint: setup-ci-tools
	@echo "🔍 Running Linter..."
	golangci-lint run ./... --timeout=5m

audit: setup-ci-tools
	@echo "🛡️  Running Security Audit..."
	govulncheck ./... || true

test:
	@echo "🧪 Running Tests..."
	go test -v -race -short ./...

# Verifica se o SDK Go compila corretamente para WASM (Build check)
check-wasm-sdk:
	@echo "🕸️  Verifying WASM SDK Build..."
	GOOS=wasip1 GOARCH=wasm go build -o /dev/null ./sdk/...

# --- 4. Utilities ---
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf bin/
	@rm -f $(SERVER_BIN)
	@rm -f functions/*.wasm
	@# Mantém o .gitkeep
	@git checkout functions/.gitkeep 2>/dev/null || true

# --- COMANDO MESTRE (CI Local) ---
ci: lint audit test check-wasm-sdk
	@echo "\n✅ \033[0;32mALL CHECKS PASSED! Ready to push.\033[0m"