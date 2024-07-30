# Используем официальный образ Golang как базовый образ
FROM golang:1.20-alpine

# Устанавливаем необходимые пакеты
RUN apk update && apk add --no-cache git

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/server

# Указываем команду запуска контейнера
CMD ["./main"]

# Указываем порт, который будет прослушивать контейнер
EXPOSE 8080
