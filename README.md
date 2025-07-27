# Social Network Backend

Бэкенд для социальной сети на Go 1.24 с созданием и просмотром анкет пользователей.

## Особенности

- 🔐 JWT авторизация с bcrypt хешированием паролей
- 👤 Регистрация и управление профилями пользователей  
- 🔍 Поиск и фильтрация анкет
- 📚 Swagger документация API
- 🗄️ PostgreSQL с sqlc для генерации запросов
- 🏗️ Clean Architecture + DDD
- 📝 Структурированные логи с zap
- ⚙️ Конфигурация через .env файлы

## API Эндпойнты

### Публичные
- `POST /api/v1/register` - Регистрация пользователя
- `POST /api/v1/login` - Авторизация
- `GET /api/v1/profile/{id}` - Просмотр анкеты по ID
- `GET /api/v1/profiles` - Поиск анкет с фильтрацией

### Защищенные (требуют JWT токен)
- `GET /api/v1/profile/me` - Просмотр собственной анкеты
- `POST /api/v1/profile` - Создание анкеты
- `PUT /api/v1/profile/me` - Редактирование анкеты

## Быстрый старт

### 1. Клонирование и установка зависимостей

```bash
git clone <repository>
cd experiment
go mod download
```

### 2. Настройка базы данных

Запустите PostgreSQL с помощью Docker:
```bash
docker-compose up -d
```

Или создайте базу данных вручную:
```sql
CREATE DATABASE social_network;
```

Выполните миграции с помощью goose:
```bash
make migrate
# или напрямую:
go run cmd/migrations/main.go up
```

### 3. Конфигурация

Скопируйте файл с примером переменных окружения:
```bash
cp .env.example .env
```

Отредактируйте `.env` файл с вашими настройками:
```env
# Конфигурация сервера
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Конфигурация базы данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=social_network
DB_SSLMODE=disable

# JWT конфигурация
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY_HOURS=24
```

### 4. Запуск приложения

```bash
go run cmd/server/main.go
```

Приложение будет доступно по адресу: http://localhost:8080

## Документация API

После запуска сервера, Swagger документация доступна по адресу:
http://localhost:8080/swagger/

## Примеры использования

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Создание профиля
```bash
curl -X POST http://localhost:8080/api/v1/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "first_name": "Иван",
    "last_name": "Иванов", 
    "age": 25,
    "gender": "male",
    "city": "Москва",
    "interests": ["программирование", "музыка", "спорт"]
  }'
```

### Поиск профилей
```bash
curl "http://localhost:8080/api/v1/profiles?gender=male&city=Москва&interests=программирование&limit=10&offset=0"
```

## Структура проекта

```
.
├── cmd/
│   ├── server/          # Точка входа приложения
│   └── migrations/      # Команда для управления миграциями
├── internal/
│   ├── config/          # Конфигурация
│   ├── domain/          # Доменный слой (DDD)
│   │   ├── entities/    # Доменные сущности
│   │   ├── repositories/# Интерфейсы репозиториев
│   │   └── services/    # Доменные сервисы
│   ├── infrastructure/  # Инфраструктурный слой
│   │   ├── database/    # Подключение к БД и sqlc
│   │   └── repository/  # Реализация репозиториев
│   └── interfaces/      # Слой интерфейсов
│       └── http/        # HTTP handlers, middleware, routes
├── migrations/          # SQL миграции (goose)
├── docs/               # Swagger документация
└── README.md
```

## Технологический стек

- **Go 1.24** - Язык программирования
- **PostgreSQL** - База данных
- **sqlc** - Генератор SQL кода
- **goose** - Миграции базы данных
- **chi** - HTTP роутер
- **JWT** - Авторизация
- **bcrypt** - Хеширование паролей
- **zap** - Логирование
- **Swagger** - Документация API

## Архитектура

Проект построен на принципах **Clean Architecture** и **Domain Driven Design (DDD)**:

- **Domain Layer** (`internal/domain/`) - Бизнес-логика, сущности и интерфейсы репозиториев
- **Infrastructure Layer** (`internal/infrastructure/`) - Внешние зависимости (БД, реализация репозиториев)  
- **Interface Layer** (`internal/interfaces/`) - HTTP handlers, middleware и роуты
- **Configuration** (`internal/config/`) - Настройки приложения с поддержкой env переменных

## Управление миграциями

Проект использует [goose](https://github.com/pressly/goose) для управления миграциями базы данных.

📖 **Подробное руководство:** [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)

### Основные команды:
```bash
# Применить все миграции
make migrate

# Создать новую миграцию
make migrate-create name=add_new_table

# Откатить последнюю миграцию
make migrate-down

# Показать статус миграций
make migrate-status

# Показать текущую версию БД
make migrate-version
```

## Разработка

### Генерация sqlc кода
```bash
go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate
```

### Генерация Swagger документации
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go
```

## Лицензия

MIT 