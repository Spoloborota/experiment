# Dockerfile для migrations приложения
# Многоэтапная сборка для безопасности и уменьшения размера образа

# Этап 1: Сборка
FROM golang:1.24.3-alpine3.22 AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Создаем непривилегированного пользователя
RUN adduser -D -s /bin/sh -u 1001 appuser

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей и загружаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bin/migrations \
    cmd/migrations/main.go

# Этап 2: Production образ
FROM alpine:3.22.1

# Устанавливаем необходимые для работы пакеты
RUN apk --no-cache add ca-certificates tzdata

# Создаем непривилегированного пользователя
RUN adduser -D -s /bin/sh -u 1001 appuser

# Создаем необходимые директории
RUN mkdir -p /app/migrations && chown -R appuser:appuser /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение
COPY --from=builder --chown=appuser:appuser /app/bin/migrations ./migrate

# Копируем файлы миграций
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations

# Устанавливаем переменные среды по умолчанию
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=social_network
ENV DB_SSLMODE=disable

# По умолчанию запускаем миграции up
CMD ["./migrate", "up"] 