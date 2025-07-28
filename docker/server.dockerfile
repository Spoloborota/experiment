# Dockerfile для server приложения
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
    -o bin/server \
    cmd/server/main.go

# Этап 2: Production образ
FROM alpine:3.22.1

# Устанавливаем необходимые для работы пакеты
RUN apk --no-cache add ca-certificates tzdata wget

# Создаем непривилегированного пользователя
RUN adduser -D -s /bin/sh -u 1001 appuser

# Создаем необходимые директории
RUN mkdir -p /app && chown -R appuser:appuser /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение
COPY --from=builder --chown=appuser:appuser /app/bin/server ./server

# Копируем swagger документацию
COPY --from=builder --chown=appuser:appuser /app/docs ./docs

# Открываем порт для приложения
EXPOSE 8080

# Устанавливаем переменные среды по умолчанию
ENV SERVER_PORT=8080
ENV SERVER_HOST=0.0.0.0

# Настраиваем healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Запускаем приложение
CMD ["./server"] 