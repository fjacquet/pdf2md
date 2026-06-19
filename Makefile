# Canonical Go Makefile — fjacquet/ci standard interface (do not rename targets)
.DEFAULT_GOAL := all
DIST  ?= dist
COVER ?= coverage.out
GOLANGCI_VERSION ?= v2.12.2
GORELEASER_VERSION ?= v2.16.0

BINARY_NAME ?= pdf2md

.PHONY: all clean install tools lint format test build vuln sbom security docs coverage-upload release ci \
        benchmark fmt run-test

all: clean lint test build

clean:
	rm -rf $(DIST) site $(COVER) *.sarif
	rm -f $(BINARY_NAME) test.md

install:
	go mod download

tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/goreleaser/goreleaser/v2@$(GORELEASER_VERSION)

lint:
	golangci-lint run --timeout=5m

format:
	golangci-lint fmt

test:
	go test -race -coverprofile=$(COVER) -covermode=atomic ./...

build:
	go build -v ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

sbom:
	mkdir -p $(DIST)
	go run github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest mod -json -output $(DIST)/sbom.cdx.json

security:  # advisory: reports findings but never blocks the build (CodeQL/osv are the blocking gates)
	uvx semgrep scan --config auto --skip-unknown-extensions || true

docs:
	uvx --with mkdocs-material --with pymdown-extensions mkdocs build --strict --site-dir site

coverage-upload:
	uvx --from codecov-cli codecov upload-process --file $(COVER) || true

release:
	goreleaser release --clean

ci: lint test build vuln

# ---------------------------------------------------------------------------
# Convenience targets (preserved from the original Makefile)
# ---------------------------------------------------------------------------

## benchmark: Run the benchmark harness
benchmark:
	go run cmd/benchmark/main.go

## fmt: Format code
fmt:
	go fmt ./...

## run-test: Build and run the tool on the bundled test fixture
run-test: build
	go build -o $(BINARY_NAME) cmd/pdf2md/main.go
	./$(BINARY_NAME) testdata/test.pdf test.md
