# Dockerfile для debug версии server приложения с Delve
# Используется для разработки с возможностью подключения дебаггера

FROM golang:1.24.3-alpine3.22

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata wget

# Устанавливаем Delve дебаггер
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Создаем непривилегированного пользователя
RUN adduser -D -s /bin/sh -u 1001 appuser

# Создаем необходимые директории
RUN mkdir -p /app && chown -R appuser:appuser /app

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей и загружаем их
COPY --chown=appuser:appuser go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY --chown=appuser:appuser . .

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт для приложения
EXPOSE 8080

# Открываем порт для дебаггера
EXPOSE 2345

# Устанавливаем переменные среды по умолчанию
ENV SERVER_PORT=8080
ENV SERVER_HOST=0.0.0.0

# Компилируем с флагами для дебаггинга и запускаем через Delve
CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "debug", "./cmd/server/main.go"] 