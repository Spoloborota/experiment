# Максимально безопасный Dockerfile для production
# Follows OWASP Container Security best practices

# Этап 1: Сборка в безопасной среде
FROM golang:1.24.3-alpine3.22@sha256:b4f875e650466fa0fe62c6fd3f02517a392123eea85f1d7e69d85f780e4db1c1 AS builder

# Создаем непривилегированного пользователя для сборки
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Устанавливаем только необходимые пакеты + обновления безопасности
RUN apk add --no-cache \
    git=2.47.1-r0 \
    ca-certificates=20241010-r0 \
    tzdata=2025a-r0 && \
    apk upgrade --no-cache

# Переключаемся на непривилегированного пользователя для сборки
USER appuser

# Создаем рабочую директорию
WORKDIR /app

# Копируем только файлы зависимостей (для кеширования)
COPY --chown=appuser:appgroup go.mod go.sum ./

# Загружаем зависимости с проверкой
RUN go mod download && go mod verify

# Копируем исходный код
COPY --chown=appuser:appgroup . .

# Собираем с максимальными оптимизациями безопасности
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -buildmode=exe \
    -ldflags='-w -s -extldflags "-static" -X main.version=1.0.0' \
    -a \
    -installsuffix cgo \
    -tags netgo,osusergo \
    -trimpath \
    -o bin/server \
    cmd/server/main.go

# Этап 2: Distroless production образ (самый безопасный)
FROM gcr.io/distroless/static-debian12:nonroot@sha256:8dd8d3ca2cf283383304fd45a5c9c74d5f2cd9da8d3b077d720e264880077c65

# Копируем пользователя из builder (nonroot = UID 65532)
# Уже есть в distroless, дополнительно ничего не нужно

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только исполняемый файл и документацию
COPY --from=builder --chown=65532:65532 /app/bin/server ./server
COPY --from=builder --chown=65532:65532 /app/docs ./docs

# Открываем только необходимый порт
EXPOSE 8080

# Устанавливаем переменные среды
ENV SERVER_PORT=8080 \
    SERVER_HOST=0.0.0.0 \
    GIN_MODE=release \
    CGO_ENABLED=0

# Добавляем лейблы для безопасности и мониторинга
LABEL \
    org.opencontainers.image.title="Social Network Server" \
    org.opencontainers.image.description="Production-ready social network API" \
    org.opencontainers.image.vendor="Your Company" \
    org.opencontainers.image.version="1.0.0" \
    org.opencontainers.image.created="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    org.opencontainers.image.source="https://github.com/your-org/experiment" \
    maintainer="your-team@company.com" \
    security.scan.enabled="true"

# Healthcheck с улучшенной безопасностью
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD ["/app/server", "--health-check"] || exit 1

# Запускаем приложение как исполняемый файл
ENTRYPOINT ["/app/server"]
CMD [] 