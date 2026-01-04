.PHONY: help demo up down logs restart clean seed generator build test

help: ## Показать помощь
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

demo: ## Полное развертывание демо (Docker + seed + generator)
	@echo "🚀 Запуск полного демо развертывания..."
	@echo ""
	@echo "📦 1/4 Поднимаем Docker Compose..."
	@docker-compose up -d
	@echo ""
	@echo "⏳ 2/4 Ждем готовности PostgreSQL (15 сек)..."
	@sleep 15
	@echo ""
	@echo "🌱 3/4 Заполняем базу тестовыми данными..."
	@cd backend && $(MAKE) seed
	@echo ""
	@echo "⚙️  4/4 Генерируем конфиги Asterisk..."
	@cd backend && $(MAKE) generator
	@echo ""
	@echo "✅ Демо развернуто успешно!"
	@echo ""
	@echo "📍 API доступен по адресу: http://localhost:8080"
	@echo "📍 Проверка: curl http://localhost:8080/api/profiles"
	@echo "📍 Конфиги: ./backend/results/"
	@echo ""

up: ## Поднять Docker Compose
	@echo "🐳 Поднимаем Docker Compose..."
	@docker-compose up -d
	@echo "✅ Сервисы запущены"

down: ## Остановить Docker Compose
	@echo "🛑 Останавливаем Docker Compose..."
	@docker-compose down
	@echo "✅ Сервисы остановлены"

restart: ## Перезапустить Docker Compose
	@echo "🔄 Перезапускаем сервисы..."
	@docker-compose restart
	@echo "✅ Сервисы перезапущены"

logs: ## Показать логи всех сервисов
	@docker-compose logs -f

logs-backend: ## Показать логи backend
	@docker-compose logs -f backend

logs-postgres: ## Показать логи postgres
	@docker-compose logs -f postgres

seed: up ## Заполнить базу тестовыми данными
	@echo "🌱 Заполняем базу данных..."
	@sleep 5
	@cd backend && $(MAKE) seed

generator: ## Запустить генератор конфигов
	@echo "⚙️  Генерируем конфиги..."
	@cd backend && $(MAKE) generator

build: ## Пересобрать Docker образы
	@echo "🔨 Пересобираем Docker образы..."
	@docker-compose build
	@echo "✅ Образы пересобраны"

clean: down ## Полная очистка (контейнеры + volumes + конфиги)
	@echo "🗑️  Очищаем все данные..."
	@docker-compose down -v
	@rm -rf backend/results/
	@rm -rf backend/bin/
	@echo "✅ Все очищено"

test: ## Запустить тесты
	@echo "🧪 Запускаем тесты..."
	@cd backend && go test ./...

status: ## Показать статус сервисов
	@docker-compose ps

shell-backend: ## Войти в shell backend контейнера
	@docker-compose exec backend sh

shell-postgres: ## Войти в psql консоль
	@docker-compose exec postgres psql -U postgres -d asterisk_manager

api-test: ## Протестировать API endpoints
	@echo "📡 Тестируем API endpoints..."
	@echo ""
	@echo "GET /:"
	@curl -s http://localhost:8080/ && echo ""
	@echo ""
	@echo "GET /api/locations:"
	@curl -s http://localhost:8080/api/locations | jq -r '.[0:2] | length' | xargs -I {} echo "  ✓ {} локации получены"
	@echo ""
	@echo "GET /api/devices:"
	@curl -s http://localhost:8080/api/devices | jq -r '.[0:2] | length' | xargs -I {} echo "  ✓ {} устройства получены"
	@echo ""
	@echo "GET /api/profiles:"
	@curl -s http://localhost:8080/api/profiles | jq -r '.[0:2] | length' | xargs -I {} echo "  ✓ {} профиля получены"
	@echo ""
	@echo "✅ API работает корректно!"

.DEFAULT_GOAL := help
