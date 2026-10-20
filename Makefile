.DEFAULT_GOAL := help

.PHONY: help install build dev dev-raycast build-raycast lint-raycast fix-raycast build-tui tui install-tui list add clean clean-all

help: ## Muestra este menú de ayuda
	@echo ""
	@echo "👻 Foltty Monorepo (Raycast + Go TUI)"
	@echo "======================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

install: ## Instala dependencias de ambos proyectos (Raycast y Go)
	cd raycast && pnpm install
	cd tui && go mod download

build: build-tui build-raycast ## Compila ambos proyectos (Go binario y extensión Raycast)

# --- Go TUI / CLI ---
build-tui: ## Compila el binario de Go en bin/foltty
	mkdir -p bin
	cd tui && go build -o ../bin/foltty ./cmd/foltty

tui: build-tui ## Ejecuta la TUI interactiva en Go
	./bin/foltty

install-tui: ## Instala 'foltty' globalmente en tu sistema (~/go/bin)
	cd tui && go install ./cmd/foltty

list: build-tui ## Lista los proyectos registrados usando el binario de Go
	./bin/foltty list

add: build-tui ## Agrega el directorio actual a la lista
	./bin/foltty add

# --- Raycast Extension ---
dev: dev-raycast ## Alias para dev-raycast

dev-raycast: ## Inicia el servidor de desarrollo de Raycast con Hot Reload
	cd raycast && pnpm run dev

build-raycast: ## Compila la extensión de Raycast para distribución (raycast/dist/)
	cd raycast && pnpm run build

lint-raycast: ## Valida el linter oficial de Raycast
	cd raycast && pnpm run lint

fix-raycast: ## Corrige automáticamente estilos y formato de Raycast
	cd raycast && pnpm run fix-lint

# --- Limpieza ---
clean: ## Limpia los artefactos generados (bin/ y raycast/dist)
	rm -rf bin raycast/dist

clean-all: clean ## Limpia bin/, raycast/dist y node_modules
	rm -rf raycast/node_modules
