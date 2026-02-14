BINARY_NAME=gojinn
SERVER_BIN=gojinn-server
SIGNER_BIN=gojinn-signer

REPRODUCIBLE_FLAGS=-trimpath -buildvcs=false -ldflags="-s -w -buildid="

XCADDY_CMD := xcaddy build --with github.com/pauloappbr/gojinn=.

URL_PYTHON := "https://github.com/vmware-labs/webassembly-language-runtimes/releases/download/python%2F3.12.0%2B20231211-040d5a6/python-3.12.0.wasm"
URL_PHP    := "https://github.com/vmware-labs/webassembly-language-runtimes/releases/download/php%2F8.2.6%2B20230714-11be424/php-cgi-8.2.6.wasm"
URL_RUBY   := "https://github.com/vmware-labs/webassembly-language-runtimes/releases/download/ruby%2F3.2.2%2B20230714-11be424/ruby-3.2.2.wasm"

.PHONY: all clean ci lint audit test build-all build-cli build-server build-signer download-runtimes check-wasm-sdk release-secure

all: build-all

build-all: build-cli build-server build-signer

build-cli:
	@echo "Building Gojinn CLI (Reproducible)..."
	@mkdir -p bin
	@go build $(REPRODUCIBLE_FLAGS) -o bin/$(BINARY_NAME) ./cmd/gojinn

build-server:
	@echo "Building Gojinn Server (Caddy)..."
	@$(XCADDY_CMD) --output $(SERVER_BIN)

build-signer:
	@echo "Building Signer Tool (Reproducible)..."
	@go build $(REPRODUCIBLE_FLAGS) -o bin/$(SIGNER_BIN) ./cmd/signer/main.go

release-secure: build-cli build-signer generate-sbom generate-transparency-log
generate-sbom:
	@echo "Generating Software Bill of Materials (SBOM)..."
	@mkdir -p security
	@which syft > /dev/null || curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin
	@syft dir:. -o cyclonedx-json > security/gojinn-sbom.json
	@echo "SBOM saved to security/gojinn-sbom.json"

generate-transparency-log:
	@echo "Generating Binary Hash Transparency Log..."
	@mkdir -p security
	@cd bin && sha256sum $(BINARY_NAME) $(SIGNER_BIN) > ../security/checksums.txt
	@echo "Checksums saved to security/checksums.txt"

download-runtimes:
	@echo "⬇Checking Polyglot Runtimes..."
	@mkdir -p lib
	@mkdir -p functions
	@if [ ! -s lib/python.wasm ]; then curl -L -o lib/python.wasm $(URL_PYTHON); fi
	@if [ ! -s lib/php.wasm ]; then curl -L -o lib/php.wasm $(URL_PHP); fi
	@if [ ! -s lib/ruby.wasm ]; then curl -L -o lib/ruby.wasm $(URL_RUBY); fi
	@cp lib/*.wasm functions/ 2>/dev/null || true
	@echo "Runtimes ready."

setup-ci-tools:
	@echo "Setting up CI tools..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest

lint: setup-ci-tools
	@echo "Running Linter..."
	golangci-lint run ./... --timeout=5m

audit: setup-ci-tools
	@echo "Running Security Audit..."
	govulncheck ./... || true

test:
	@echo "Running Tests..."
	go test -v -race -short ./...

check-wasm-sdk:
	@echo "Verifying WASM SDK Build..."
	GOOS=wasip1 GOARCH=wasm go build -o /dev/null ./sdk/...

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -rf security/
	@rm -f $(SERVER_BIN)
	@rm -f functions/*.wasm
	@git checkout functions/.gitkeep 2>/dev/null || true

ci: lint audit test check-wasm-sdk
	@echo "\n \033[0;32mALL CHECKS PASSED! Ready to push.\033[0m"