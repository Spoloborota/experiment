# ==========================================
# Main Makefile for Social Network Project
# ==========================================

# Включаем модульные Makefile файлы
include docker.mk

.PHONY: build run test clean sqlc swagger deps migrate migrate-create migrate-down migrate-status migrate-reset migrate-version
.PHONY: regen dev help

# === ОСНОВНЫЕ КОМАНДЫ РАЗРАБОТКИ ===

# Сборка проекта
build:
	go build -o bin/server cmd/server/main.go

# Запуск приложения локально
run:
	go run cmd/server/main.go

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Запуск тестов
test:
	go test -v ./...

# Очистка сгенерированных файлов
clean:
	rm -rf bin/
	rm -rf docs/swagger.json docs/swagger.yaml

# === ГЕНЕРАЦИЯ КОДА ===

# Генерация sqlc кода
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

# Генерация swagger документации
swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go --parseDependency --parseInternal

# Полная перегенерация
regen: clean sqlc swagger
	@echo "✅ Код перегенерирован!"

# === РАБОТА С МИГРАЦИЯМИ (локально) ===

# Применение миграций с goose
migrate:
	@echo "📦 Применение миграций с goose..."
	go run cmd/migrations/main.go up

# Создание новой миграции
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	go run cmd/migrations/main.go create $(name)

# Откат последней миграции
migrate-down:
	go run cmd/migrations/main.go down

# Статус миграций
migrate-status:
	go run cmd/migrations/main.go status

# Сброс всех миграций
migrate-reset:
	go run cmd/migrations/main.go reset

# Версия базы данных
migrate-version:
	go run cmd/migrations/main.go version

# === КОМПЛЕКСНЫЕ КОМАНДЫ ===

# Быстрый старт для разработки
dev: deps regen build
	@echo "🚀 Проект готов к запуску!"
	@echo ""
	@echo "📋 Следующие шаги:"
	@echo "1. Настройте .env файл"
	@echo "2. Запустите PostgreSQL: make docker-up (только postgres)"
	@echo "3. Примените миграции: make migrate"
	@echo "4. Запустите сервер: make run"
	@echo ""
	@echo "🐳 Или используйте Docker: make docker-dev"

# === ПОМОЩЬ ===

# Показать все доступные команды
help:
	@echo "🛠️  Доступные команды для Social Network проекта:"
	@echo ""
	@echo "🔧 РАЗРАБОТКА:"
	@echo "  build          - Собрать приложение"
	@echo "  run            - Запустить локально"
	@echo "  test           - Запустить тесты"
	@echo "  dev            - Быстрый старт разработки"
	@echo ""
	@echo "📝 ГЕНЕРАЦИЯ КОДА:"
	@echo "  sqlc           - Генерация SQLC кода"
	@echo "  swagger        - Генерация Swagger документации"
	@echo "  regen          - Полная перегенерация"
	@echo ""
	@echo "🗃️  МИГРАЦИИ (локально):"
	@echo "  migrate        - Применить миграции"
	@echo "  migrate-create - Создать миграцию"
	@echo "  migrate-status - Статус миграций"
	@echo ""
	@echo "🐳 DOCKER:"
	@echo "  docker-help    - Показать Docker команды"
	@echo "  docker-dev     - Быстрый старт в Docker"
	@echo ""
	@echo "🧹 ОЧИСТКА:"
	@echo "  clean          - Очистить сгенерированные файлы"
	@echo ""

# Команда по умолчанию
.DEFAULT_GOAL := help 