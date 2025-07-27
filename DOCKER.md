# Docker Deployment Guide

Этот документ описывает как запускать приложение в Docker окружении с возможностью отладки.

## Требования

- Docker 20.10+
- Docker Compose 2.0+
- Make (для удобных команд)

## Структура Docker файлов

### Dockerfile.server
Production версия сервера с многоэтапной сборкой:
- Использует фиксированные версии образов
- Создает непривилегированного пользователя
- Оптимизирован для минимального размера
- Включает health check

### Dockerfile.debug
Debug версия с Delve дебаггером:
- Включает исходный код для live reload
- Настроенный Delve дебаггер на порту 2345
- Удобен для разработки

### Dockerfile.migrations
Контейнер для миграций базы данных:
- Изолированное выполнение миграций
- Возможность ручного управления

## Быстрый старт

### Production запуск
```bash
# Полная сборка и запуск
make docker-dev

# Или пошагово:
make docker-build
make docker-up
```

Приложение будет доступно:
- **API:** http://localhost:8080
- **Swagger UI:** http://localhost:8080/swagger/
- **PostgreSQL:** localhost:6632

### Debug запуск
```bash
# Запуск с дебаггером
make docker-dev-debug
```

Приложение будет доступно:
- **API:** http://localhost:8081  
- **Delve дебаггер:** localhost:2345
- **PostgreSQL:** localhost:6632

## Подключение дебаггера

### VS Code
Добавьте в `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Connect to Docker Delve",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "/app",
            "port": 2345,
            "host": "127.0.0.1",
            "program": "${workspaceFolder}",
            "env": {},
            "args": []
        }
    ]
}
```

### GoLand/IntelliJ
1. Run -> Edit Configurations
2. Add "Go Remote"
3. Host: localhost, Port: 2345

## Управление сервисами

### Основные команды
```bash
# Сборка образов
make docker-build

# Запуск всех сервисов
make docker-up

# Запуск с дебаггером
make docker-debug

# Остановка сервисов
make docker-down

# Полная очистка (включая данные)
make docker-reset
```

### Логи
```bash
# Все сервисы
make docker-logs

# Конкретный сервис
make docker-logs-server
make docker-logs-migrations
make docker-logs-postgres
```

## Работа с миграциями

### В Docker окружении
```bash
# Применить миграции
make docker-migrate

# Статус миграций
make docker-migrate-status

# Создать новую миграцию
make docker-migrate-create name=add_new_table

# Откатить миграцию
make docker-migrate-down
```

### Подключение к базе данных
```bash
# Прямое подключение к PostgreSQL
make docker-db-shell
```

## Профили Docker Compose

### Стандартный запуск
```bash
docker-compose up -d
```
Запускает: postgres, server, migrations

### Debug профиль
```bash
docker-compose --profile debug up -d
```
Запускает: postgres, server-debug, migrations

### Tools профиль
```bash
docker-compose --profile tools run migrations-cli ./migrations status
```
Для ручного управления миграциями

## Переменные среды

Основные переменные настроены в `docker-compose.yml`:

```yaml
environment:
  - SERVER_PORT=8080
  - SERVER_HOST=0.0.0.0
  - DB_HOST=postgres
  - DB_PORT=5432
  - DB_USER=postgres
  - DB_PASSWORD=postgres
  - DB_NAME=social_network
  - DB_SSLMODE=disable
  - JWT_SECRET=docker-development-secret-key-change-in-production
  - JWT_EXPIRY_HOURS=24
```

Для production необходимо изменить `JWT_SECRET` и `DB_PASSWORD`.

## Health Checks

Приложение включает health check эндпоинты:
- `/health` - основной health check
- `/api/v1/health` - API health check

Docker автоматически проверяет состояние сервиса через health check.

## Мониторинг

### Статус контейнеров
```bash
docker-compose ps
```

### Ресурсы
```bash
docker stats
```

### Логи в реальном времени
```bash
docker-compose logs -f server
```

## Troubleshooting

### Контейнер не запускается
1. Проверьте логи: `make docker-logs-server`
2. Убедитесь что PostgreSQL готов: `make docker-logs-postgres`
3. Проверьте что миграции применились: `make docker-logs-migrations`

### Проблемы с подключением к БД
1. Проверьте что PostgreSQL health check проходит
2. Убедитесь что миграции выполнились успешно
3. Проверьте переменные среды в docker-compose.yml

### Дебаггер не подключается
1. Убедитесь что используете debug профиль: `make docker-debug`
2. Проверьте что порт 2345 не занят: `lsof -i :2345`
3. Проверьте логи debug контейнера: `docker-compose logs server-debug`

### Медленная сборка
1. Очистите Docker кеш: `docker system prune -a`
2. Пересоберите без кеша: `make docker-reset`

## Безопасность

### Production рекомендации
1. Измените `JWT_SECRET` на случайную строку
2. Используйте сильный пароль для PostgreSQL
3. Не используйте debug версию в production
4. Настройте SSL/TLS для внешних подключений
5. Ограничьте доступ к портам через firewall

### Пользователи
Все контейнеры запускаются под непривилегированным пользователем (UID 1001) для повышения безопасности.

## Kubernetes Deployment

Для развертывания в Kubernetes создайте ConfigMap и Secrets:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: social-network-config
data:
  SERVER_PORT: "8080"
  SERVER_HOST: "0.0.0.0"
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_NAME: "social_network"
  DB_SSLMODE: "disable"
---
apiVersion: v1
kind: Secret
metadata:
  name: social-network-secrets
type: Opaque
stringData:
  DB_USER: "postgres"
  DB_PASSWORD: "your-secure-password"
  JWT_SECRET: "your-jwt-secret"
```

Используйте образы из вашего Docker registry для deployment в Kubernetes. 