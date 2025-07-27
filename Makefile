.PHONY: build run test clean sqlc swagger deps migrate migrate-create migrate-down migrate-status migrate-reset migrate-version

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