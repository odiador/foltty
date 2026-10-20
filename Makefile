.DEFAULT_GOAL := help

.PHONY: help install build dev dev-raycast build-raycast lint-raycast fix-raycast build-tui tui install-tui list add preset term clean clean-all

help: ## Show this help menu
	@echo ""
	@echo "👻 Foltty Monorepo (Raycast + Go TUI)"
	@echo "======================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

install: ## Install dependencies for both Raycast (pnpm) and Go
	cd raycast && pnpm install
	cd tui && go mod download

build: build-tui build-raycast ## Build both Go CLI binary and Raycast extension

# --- Go TUI / CLI ---
build-tui: ## Build Go binary into bin/foltty
	mkdir -p bin
	cd tui && go build -ldflags="-s -w" -o ../bin/foltty ./cmd/foltty

tui: build-tui ## Launch interactive terminal UI
	./bin/foltty

install-tui: ## Install 'foltty' binary into your system PATH (~/go/bin)
	cd tui && go install ./cmd/foltty

list: build-tui ## List registered projects via Go CLI
	./bin/foltty list

add: build-tui ## Register current directory in Foltty
	./bin/foltty add

preset: build-tui ## List available folder presets (Home, Docs, Dev, Customs)
	./bin/foltty preset

term: build-tui ## Inspect and list detected terminal emulators
	./bin/foltty terminal

# --- Raycast Extension ---
dev: dev-raycast ## Alias for dev-raycast

dev-raycast: ## Start Raycast extension development server with Hot Reload
	cd raycast && pnpm run dev

build-raycast: ## Build production Raycast extension bundle (raycast/dist/)
	cd raycast && pnpm run build

lint-raycast: ## Run official Raycast linter
	cd raycast && pnpm run lint

fix-raycast: ## Automatically fix style and linter issues in Raycast
	cd raycast && pnpm run fix-lint

# --- Cleanup ---
clean: ## Remove generated build artifacts (bin/ and raycast/dist)
	rm -rf bin raycast/dist

clean-all: clean ## Remove bin/, raycast/dist and node_modules
	rm -rf raycast/node_modules
