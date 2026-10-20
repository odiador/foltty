.DEFAULT_GOAL := help

.PHONY: help install dev build lint fix clean clean-all tui link list add

help: ## Muestra este menú de ayuda
	@echo ""
	@echo "👻 Foltty - Project Manager para Ghostty (Raycast & TUI)"
	@echo "========================================================="
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

install: ## Instala todas las dependencias usando pnpm
	pnpm install

dev: ## Inicia el modo de desarrollo de Raycast con Hot Reload
	pnpm run dev

build: ## Compila la extensión para producción (directorio dist/)
	pnpm run build

lint: ## Ejecuta las verificaciones de formato y linter de Raycast
	pnpm run lint

fix: ## Corrige automáticamente errores de formato y linter
	pnpm run fix-lint

tui: ## Abre la TUI interactiva en la terminal
	./bin/foltty.js

list: ## Lista todos los proyectos registrados
	./bin/foltty.js list

add: ## Agrega el directorio actual a Foltty
	./bin/foltty.js add

link: ## Enlaza el comando 'foltty' globalmente en tu sistema
	pnpm link --global

clean: ## Elimina los archivos de compilación generados (dist/)
	rm -rf dist

clean-all: clean ## Elimina dist/ y node_modules
	rm -rf node_modules
