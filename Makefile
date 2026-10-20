.DEFAULT_GOAL := help

.PHONY: help install install-tui install-raycast build dev dev-raycast build-raycast lint-raycast fix-raycast build-tui tui list add preset term clean clean-all

help: ## Show this help menu
	@echo ""
	@echo "👻 Foltty Monorepo (Raycast + Go TUI)"
	@echo "======================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

install: install-tui install-raycast ## Install both Go binary and Raycast extension

install-tui: build-tui ## Install 'foltty' binary into your system PATH (~/go/bin or ~/.local/bin)
	mkdir -p $(HOME)/.local/bin
	cp bin/foltty $(HOME)/.local/bin/foltty

install-raycast: ## Install and bundle Raycast extension
	cd raycast && npm ci
	cd raycast && npm run build

build: build-tui build-raycast ## Build both Go CLI binary and Raycast extension

# --- Go TUI / CLI ---
build-tui: ## Build Go binary into bin/foltty
	mkdir -p bin
	cd tui && go build -ldflags="-s -w" -o ../bin/foltty ./cmd/foltty

tui: build-tui ## Launch interactive terminal UI
	./bin/foltty

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
	cd raycast && npm run dev

build-raycast: ## Build production Raycast extension bundle (raycast/dist/)
	cd raycast && npm run build

lint-raycast: ## Run official Raycast linter
	cd raycast && npm run lint

fix-raycast: ## Automatically fix style and linter issues in Raycast
	cd raycast && npm run fix-lint

service-status: ## Check background LaunchAgent daemon status
	@launchctl list | grep foltty || echo "Daemon not running"

service-start: ## Start background daemon so extension works without open terminal
	@launchctl load $(HOME)/Library/LaunchAgents/com.foltty.raycast.plist 2>/dev/null || true
	@echo "✅ Foltty Raycast background service active"

service-stop: ## Stop background daemon
	@launchctl unload $(HOME)/Library/LaunchAgents/com.foltty.raycast.plist 2>/dev/null || true
	@echo "🛑 Foltty Raycast background service stopped"

# --- Cleanup ---
clean: ## Remove generated build artifacts (bin/ and raycast/dist)
	rm -rf bin raycast/dist

clean-all: clean ## Remove bin/, raycast/dist and node_modules
	rm -rf raycast/node_modules
