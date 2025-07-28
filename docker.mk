# ==========================================
# Docker-related Makefile
# ==========================================
# Этот файл содержит все команды для работы с Docker
# Включается в основной Makefile через: include Makefile.docker

.PHONY: docker-build docker-up docker-down docker-logs docker-debug docker-migrate docker-reset
.PHONY: docker-logs-server docker-logs-migrations docker-logs-postgres docker-migrate-status
.PHONY: docker-migrate-create docker-migrate-down docker-db-shell docker-dev docker-dev-debug
.PHONY: docker-clean docker-prune docker-secure-build

# === ОСНОВНЫЕ DOCKER КОМАНДЫ ===

# Сборка всех Docker образов
docker-build:
	@echo "🔨 Сборка Docker образов..."
	docker-compose build

# Сборка безопасной версии
docker-secure-build:
	@echo "🔒 Сборка secure версии сервера..."
	docker build -f docker/server-secure.dockerfile -t experiment-server-secure .

# Запуск всех сервисов в Docker
docker-up:
	@echo "🚀 Запуск сервисов в Docker..."
	docker-compose up -d

# Запуск с debug версией сервера (с дебаггером)
docker-debug:
	@echo "🐛 Запуск в debug режиме с Delve дебаггером..."
	@echo "Дебаггер будет доступен на порту 2345"
	@echo "Сервер будет доступен на порту 8081"
	docker-compose --profile debug up -d

# Остановка всех сервисов
docker-down:
	@echo "⏹️  Остановка Docker сервисов..."
	docker-compose down

# Остановка и удаление всех данных
docker-reset:
	@echo "🧹 Полная очистка Docker окружения..."
	docker-compose down -v
	docker-compose build --no-cache

# === ЛОГИ И МОНИТОРИНГ ===

# Просмотр логов всех сервисов
docker-logs:
	docker-compose logs -f

# Просмотр логов конкретного сервиса
docker-logs-server:
	docker-compose logs -f server

docker-logs-migrations:
	docker-compose logs -f migrations

docker-logs-postgres:
	docker-compose logs -f postgres

# Статус контейнеров
docker-status:
	@echo "📊 Статус контейнеров:"
	docker-compose ps

# === РАБОТА С МИГРАЦИЯМИ ===

# Выполнение миграций в Docker окружении
docker-migrate:
	@echo "📦 Выполнение миграций в Docker..."
	docker-compose run --rm migrations up

# Статус миграций в Docker
docker-migrate-status:
	docker-compose --profile tools run --rm migrations-cli ./migrate status

# Создание новой миграции в Docker
docker-migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make docker-migrate-create name=migration_name"; \
		exit 1; \
	fi
	docker-compose --profile tools run --rm migrations-cli ./migrate create $(name)

# Откат миграций в Docker
docker-migrate-down:
	docker-compose --profile tools run --rm migrations-cli ./migrate down

# === УТИЛИТЫ ===

# Подключение к PostgreSQL в Docker
docker-db-shell:
	docker-compose exec postgres psql -U postgres -d social_network

# Очистка неиспользуемых образов и контейнеров
docker-clean:
	@echo "🧹 Очистка неиспользуемых Docker ресурсов..."
	docker system prune -f

# Агрессивная очистка (включая образы)
docker-prune:
	@echo "🗑️  Агрессивная очистка Docker (включая образы)..."
	docker system prune -a -f

# Анализ размера образов
docker-analyze:
	@echo "📏 Анализ размера Docker образов:"
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | grep experiment

# === КОМПЛЕКСНЫЕ КОМАНДЫ ===

# Полный перезапуск для разработки
docker-dev: docker-down docker-build docker-up
	@echo "✅ Docker окружение готово!"
	@echo "🌐 Сервер доступен на: http://localhost:8080"
	@echo "📚 Swagger UI: http://localhost:8080/swagger/"
	@echo "🐘 PostgreSQL: localhost:6632"

# Полный перезапуск в debug режиме
docker-dev-debug: docker-down docker-build docker-debug
	@echo "✅ Docker debug окружение готово!"
	@echo "🌐 Сервер доступен на: http://localhost:8081"
	@echo "🐛 Delve дебаггер: localhost:2345"
	@echo "🐘 PostgreSQL: localhost:6632"

# Быстрый рестарт сервера (без пересборки)
docker-restart-server:
	@echo "🔄 Перезапуск только сервера..."
	docker-compose restart server

# Проверка health check
docker-health:
	@echo "🏥 Проверка health check эндпоинтов:"
	@curl -f http://localhost:8080/health || echo "❌ Сервер недоступен"
	@curl -f http://localhost:8080/api/v1/health || echo "❌ API недоступен"

# === БЕЗОПАСНОСТЬ ===

# Сканирование уязвимостей (если установлен trivy)
docker-scan:
	@echo "🔍 Сканирование уязвимостей образов..."
	@if command -v trivy >/dev/null 2>&1; then \
		trivy image experiment-server:latest; \
	else \
		echo "⚠️  trivy не установлен. Установите: https://github.com/aquasecurity/trivy"; \
	fi

# Запуск с дополнительными ограничениями безопасности
docker-secure-run:
	@echo "🔒 Запуск с усиленной безопасностью..."
	docker run -d \
		--name secure-server \
		--read-only \
		--no-new-privileges \
		--cap-drop=ALL \
		--security-opt=no-new-privileges \
		-p 8080:8080 \
		experiment-server-secure:latest

# === ПОМОЩЬ ===

# Показать все доступные Docker команды
docker-help:
	@echo "🐳 Доступные Docker команды:"
	@echo ""
	@echo "📦 ОСНОВНЫЕ:"
	@echo "  docker-build           - Собрать все образы"
	@echo "  docker-up              - Запустить сервисы"
	@echo "  docker-down            - Остановить сервисы"
	@echo "  docker-dev             - Быстрый старт для разработки"
	@echo ""
	@echo "🐛 ОТЛАДКА:"
	@echo "  docker-debug           - Запустить с дебаггером"
	@echo "  docker-dev-debug       - Debug окружение"
	@echo "  docker-logs            - Логи всех сервисов"
	@echo "  docker-logs-server     - Логи сервера"
	@echo ""
	@echo "🗃️  МИГРАЦИИ:"
	@echo "  docker-migrate         - Применить миграции"
	@echo "  docker-migrate-status  - Статус миграций"
	@echo ""
	@echo "🔒 БЕЗОПАСНОСТЬ:"
	@echo "  docker-secure-build    - Собрать secure образ"
	@echo "  docker-scan            - Сканировать уязвимости"
	@echo ""
	@echo "🧹 ОЧИСТКА:"
	@echo "  docker-clean           - Очистить ресурсы"
	@echo "  docker-reset           - Полная очистка"
	@echo "" 