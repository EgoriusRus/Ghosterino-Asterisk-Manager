.PHONY: help build run seed generator clean docker-up docker-down

help: ## Показать помощь
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Собрать API сервер
	@echo "🔨 Building API server..."
	@go build -o bin/api main.go routes.go

run: build ## Запустить API сервер
	@echo "🚀 Starting API server..."
	@./bin/api

seed: ## Заполнить базу тестовыми данными
	@echo "🌱 Seeding database..."
	@go run cmd/seed/main.go

generator: ## Запустить генератор конфигов
	@echo "⚙️  Running generator..."
	@go run cmd/generator/main.go

clean: ## Очистить собранные файлы
	@echo "🗑️  Cleaning..."
	@rm -rf bin/
	@rm -rf results/

docker-up: ## Поднять Docker Compose
	@echo "🐳 Starting Docker Compose..."
	@cd .. && docker-compose up -d

docker-down: ## Остановить Docker Compose
	@echo "🐳 Stopping Docker Compose..."
	@cd .. && docker-compose down

docker-logs: ## Показать логи Docker
	@cd .. && docker-compose logs -f

test: ## Запустить тесты
	@echo "🧪 Running tests..."
	@go test -v ./...

.DEFAULT_GOAL := help
