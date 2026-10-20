.DEFAULT_GOAL := help

.PHONY: help install dev build lint fix clean clean-all tui list add link build-go install-go tui-go

help: ## Muestra este menú de ayuda
	@echo ""
	@echo "👻 Foltty - Project Manager para Ghostty (Raycast & Go TUI)"
	@echo "============================================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

install: ## Instala dependencias de Raycast (pnpm) y Go
	pnpm install
	go mod download

dev: ## Inicia el modo de desarrollo de Raycast con Hot Reload
	pnpm run dev

build: ## Compila la extensión de Raycast para producción (dist/)
	pnpm run build

build-go: ## Compila el binario nativo en Go (bin/foltty-go)
	go build -o bin/foltty-go ./cmd/foltty

install-go: ## Instala el binario 'foltty' globalmente vía Go (~/go/bin)
	go install ./cmd/foltty

tui-go: build-go ## Ejecuta la TUI nativa en Go
	./bin/foltty-go

tui: ## Abre la TUI interactiva en la terminal (Node/TSX)
	./bin/foltty.js

list: ## Lista todos los proyectos registrados
	./bin/foltty.js list

add: ## Agrega el directorio actual a Foltty
	./bin/foltty.js add

link: ## Enlaza el comando 'foltty' globalmente en tu sistema vía pnpm
	pnpm link --global

lint: ## Ejecuta las verificaciones de formato y linter de Raycast
	pnpm run lint

fix: ## Corrige automáticamente errores de formato y linter
	pnpm run fix-lint

clean: ## Elimina los archivos de compilación generados (dist/ y bin/foltty-go)
	rm -rf dist bin/foltty-go

clean-all: clean ## Elimina dist/, bin/foltty-go y node_modules
	rm -rf node_modules
