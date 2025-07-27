.PHONY: build run test clean sqlc swagger deps migrate migrate-create migrate-down migrate-status migrate-reset migrate-version docker-build docker-up docker-down docker-logs docker-debug docker-migrate docker-reset

# Сборка проекта
build:
	go build -o bin/server cmd/server/main.go

# Запуск приложения
run:
	go run cmd/server/main.go

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Генерация sqlc кода
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

# Генерация swagger документации
swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go --parseDependency --parseInternal

# Очистка сгенерированных файлов
clean:
	rm -rf bin/
	rm -rf docs/swagger.json docs/swagger.yaml

# Запуск тестов
test:
	go test -v ./...

# Применение миграций с goose
migrate:
	@echo "Применение миграций с goose..."
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

# Полная перегенерация
regen: clean sqlc swagger

# Быстрый старт для разработки
dev: deps regen build
	@echo "Проект готов к запуску!"
	@echo "1. Настройте .env файл"
	@echo "2. Запустите PostgreSQL (docker-compose up -d)"
	@echo "3. Примените миграции: make migrate"
	@echo "4. Запустите сервер: make run"

# === DOCKER КОМАНДЫ ===

# Сборка всех Docker образов
docker-build:
	@echo "Сборка Docker образов..."
	docker-compose build

# Запуск всех сервисов в Docker
docker-up:
	@echo "Запуск сервисов в Docker..."
	docker-compose up -d

# Запуск с debug версией сервера (с дебаггером)
docker-debug:
	@echo "Запуск в debug режиме с Delve дебаггером..."
	@echo "Дебаггер будет доступен на порту 2345"
	@echo "Сервер будет доступен на порту 8081"
	docker-compose --profile debug up -d

# Остановка всех сервисов
docker-down:
	@echo "Остановка Docker сервисов..."
	docker-compose down

# Остановка и удаление всех данных
docker-reset:
	@echo "Полная очистка Docker окружения..."
	docker-compose down -v
	docker-compose build --no-cache

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

# Выполнение миграций в Docker окружении
docker-migrate:
	@echo "Выполнение миграций в Docker..."
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

# Подключение к PostgreSQL в Docker
docker-db-shell:
	docker-compose exec postgres psql -U postgres -d social_network

# Полный перезапуск для разработки
docker-dev: docker-down docker-build docker-up
	@echo "Docker окружение готово!"
	@echo "Сервер доступен на: http://localhost:8080"
	@echo "Swagger UI: http://localhost:8080/swagger/"
	@echo "PostgreSQL: localhost:6632"

# Полный перезапуск в debug режиме
docker-dev-debug: docker-down docker-build docker-debug
	@echo "Docker debug окружение готово!"
	@echo "Сервер доступен на: http://localhost:8081"
	@echo "Delve дебаггер: localhost:2345"
	@echo "PostgreSQL: localhost:6632" 