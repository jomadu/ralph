# Ralph — build, test, and lint
# Run from repository root. See AGENTS.md for quality gates.
# VERSION is set by CI (semantic-release) or from git: make build VERSION=$(git describe --tags --always)

BINARY ?= bin/ralph
VERSION ?= dev

LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

.PHONY: all build build-multi test lint lint-vet lint-fmt fmt clean

all: build

build:
	@mkdir -p bin
	@mkdir -p cmd/ralph/embed && cp -f docs/writing-ralph-prompts.md cmd/ralph/embed/writing-ralph-prompts.md
	go build $(LDFLAGS) -o $(BINARY) ./cmd/ralph

# Cross-build for release (GOOS/GOARCH set by CI). Output: dist/ralph-$(VERSION)-$(GOOS)-$(GOARCH)[.exe]
build-multi:
	@mkdir -p dist
	@mkdir -p cmd/ralph/embed && cp -f docs/writing-ralph-prompts.md cmd/ralph/embed/writing-ralph-prompts.md
	@case "$(GOOS)" in windows) SUF=".exe";; *) SUF="";; esac; \
	OUT="dist/ralph-$(VERSION)-$(GOOS)-$(GOARCH)$$SUF"; \
	go build $(LDFLAGS) -o "$$OUT" ./cmd/ralph; \
	echo "$$OUT"

test:
	go test ./...

lint: lint-vet lint-fmt

lint-vet:
	go vet ./...

lint-fmt:
	@test -z "$$(gofmt -s -l .)" || (echo "gofmt: the following files need 'gofmt -s -w' or 'make fmt':" && gofmt -s -l . && exit 1)

fmt:
	gofmt -s -w .

clean:
	rm -f $(BINARY)
	rm -rf dist
