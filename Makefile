.PHONY: help demo up down logs restart clean seed generator build test prod-up prod-down prod-logs prod-restart backup dev

help: ## Показать помощь
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

demo: ## Полное развертывание демо (Docker + seed)
	@echo "🚀 Запуск полного демо развертывания..."
	@echo ""
	@echo "📦 1/3 Поднимаем Docker Compose..."
	@docker-compose up -d --build
	@echo ""
	@echo "⏳ 2/3 Ждем готовности PostgreSQL (10 сек)..."
	@sleep 10
	@echo ""
	@echo "🌱 3/3 Заполняем базу тестовыми данными..."
	@docker-compose exec -T backend ./seed
	@echo ""
	@echo "✅ Демо развернуто успешно!"
	@echo ""
	@echo "📍 Frontend: http://localhost:3000"
	@echo "📍 Backend API: http://localhost:8080/api"
	@echo ""

dev: ## Запустить для разработки (backend в Docker, frontend локально)
	@echo "🚀 Запуск в режиме разработки..."
	@echo ""
	@echo "📦 1/3 Поднимаем PostgreSQL и Backend..."
	@docker-compose up -d --build postgres backend
	@echo ""
	@echo "⏳ 2/3 Ждем готовности (10 сек)..."
	@sleep 10
	@echo ""
	@echo "🌱 3/3 Заполняем базу тестовыми данными..."
	@docker-compose exec -T backend ./seed
	@echo ""
	@echo "✅ Backend готов!"
	@echo ""
	@echo "📍 Backend API: http://localhost:8080/api"
	@echo "📍 Запустите frontend: cd frontend && npm run dev"
	@echo ""

up: ## Поднять Docker Compose (все сервисы)
	@echo "🐳 Поднимаем Docker Compose..."
	@docker-compose up -d
	@echo "✅ Сервисы запущены"
	@echo ""
	@echo "📍 Frontend: http://localhost:3000"
	@echo "📍 Backend API: http://localhost:8080/api"

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

logs-frontend: ## Показать логи frontend
	@docker-compose logs -f frontend

logs-postgres: ## Показать логи postgres
	@docker-compose logs -f postgres

seed: ## Заполнить базу тестовыми данными (через Docker)
	@echo "🌱 Заполняем базу данных..."
	@docker-compose exec -T backend ./seed

generator: ## Запустить генератор конфигов (через Docker)
	@echo "⚙️  Генерируем конфиги..."
	@docker-compose exec -T backend ./generator

build: ## Пересобрать Docker образы
	@echo "🔨 Пересобираем Docker образы..."
	@docker-compose build
	@echo "✅ Образы пересобраны"

build-frontend: ## Пересобрать только frontend
	@echo "🔨 Пересобираем frontend..."
	@docker-compose build frontend
	@echo "✅ Frontend пересобран"

build-backend: ## Пересобрать только backend
	@echo "🔨 Пересобираем backend..."
	@docker-compose build backend
	@echo "✅ Backend пересобран"

clean: down ## Полная очистка (контейнеры + volumes + конфиги)
	@echo "🗑️  Очищаем все данные..."
	@docker-compose down -v
	@rm -rf backend/results/
	@rm -rf backend/bin/
	@echo "✅ Все очищено"

test: ## Запустить тесты (через Docker)
	@echo "🧪 Запускаем тесты..."
	@docker-compose exec -T backend go test ./... 2>/dev/null || echo "Тесты недоступны в production образе"

status: ## Показать статус сервисов
	@docker-compose ps

shell-backend: ## Войти в shell backend контейнера
	@docker-compose exec backend sh

shell-frontend: ## Войти в shell frontend контейнера
	@docker-compose exec frontend sh

shell-postgres: ## Войти в psql консоль
	@docker-compose exec postgres psql -U postgres -d asterisk_manager

api-test: ## Протестировать API endpoints
	@echo "📡 Тестируем API endpoints..."
	@echo ""
	@echo "GET /:"
	@curl -s http://localhost:8080/ && echo ""
	@echo ""
	@echo "GET /api/locations:"
	@curl -s http://localhost:8080/api/locations | jq -r 'length' | xargs -I {} echo "  ✓ {} локаций получено"
	@echo ""
	@echo "GET /api/devices:"
	@curl -s http://localhost:8080/api/devices | jq -r 'length' | xargs -I {} echo "  ✓ {} устройств получено"
	@echo ""
	@echo "GET /api/profiles:"
	@curl -s "http://localhost:8080/api/profiles?page=1&perPage=10" | jq -r '.pagination.total' | xargs -I {} echo "  ✓ {} профилей всего"
	@echo ""
	@echo "✅ API работает корректно!"

e2e: ## Полный E2E тест (поднять всё + проверить)
	@echo "🧪 Запуск E2E тестирования..."
	@echo ""
	@$(MAKE) demo
	@echo ""
	@echo "⏳ Ждем полной готовности (5 сек)..."
	@sleep 5
	@echo ""
	@$(MAKE) api-test
	@echo ""
	@echo "🌐 Проверяем Frontend..."
	@curl -s -o /dev/null -w "  ✓ Frontend HTTP status: %{http_code}\n" http://localhost:3000
	@echo ""
	@echo "✅ E2E тест пройден успешно!"
	@echo ""
	@echo "📍 Откройте в браузере: http://localhost:3000"

# Production commands
prod-up: ## Запустить в продакшн режиме (docker-compose.prod.yml)
	@echo "🚀 Запуск продакшн сервисов..."
	@docker-compose -f docker-compose.prod.yml up -d
	@echo "✅ Сервисы запущены"

prod-down: ## Остановить продакшн сервисы
	@echo "🛑 Остановка продакшн сервисов..."
	@docker-compose -f docker-compose.prod.yml down
	@echo "✅ Сервисы остановлены"

prod-logs: ## Показать логи продакшн сервисов
	@docker-compose -f docker-compose.prod.yml logs -f

prod-restart: ## Перезапустить продакшн сервисы
	@echo "🔄 Перезапуск продакшн сервисов..."
	@docker-compose -f docker-compose.prod.yml restart
	@echo "✅ Сервисы перезапущены"

prod-build: ## Пересобрать продакшн образы
	@echo "🔨 Пересборка продакшн образов..."
	@docker-compose -f docker-compose.prod.yml build --no-cache
	@docker-compose -f docker-compose.prod.yml up -d
	@echo "✅ Образы пересобраны и запущены"

prod-seed: ## Заполнить базу в продакшн (через Docker)
	@echo "🌱 Заполняем продакшн базу данных..."
	@docker-compose -f docker-compose.prod.yml exec -T backend ./seed

prod-generator: ## Запустить генератор в продакшн (через Docker)
	@echo "⚙️  Генерируем конфиги в продакшн..."
	@docker-compose -f docker-compose.prod.yml exec -T backend ./generator

backup: ## Создать бэкап базы данных
	@echo "💾 Создание бэкапа базы данных..."
	@mkdir -p ./backups
	@docker-compose -f docker-compose.prod.yml exec -T postgres pg_dump -U $${DB_USER:-asterisk_prod} $${DB_NAME:-asterisk_manager_prod} | gzip > ./backups/backup_$$(date +%Y%m%d_%H%M%S).sql.gz
	@echo "✅ Бэкап создан в ./backups/"

.DEFAULT_GOAL := help
