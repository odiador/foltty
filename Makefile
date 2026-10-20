.DEFAULT_GOAL := help

.PHONY: help install dev build lint fix clean clean-all test-pman

help: ## Muestra este menú de ayuda
	@echo ""
	@echo "👻 Ghostty Project Manager - Makefile"
	@echo "======================================"
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

clean: ## Elimina los archivos de compilación generados (dist/)
	rm -rf dist

clean-all: clean ## Elimina dist/ y node_modules
	rm -rf node_modules

test-pman: ## Muestra el contenido actual de ~/.pman_projects
	@echo "📁 Contenido de ~/.pman_projects:"
	@cat ~/.pman_projects 2>/dev/null || echo "(El archivo no existe o está vacío)"
