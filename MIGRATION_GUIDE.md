# Руководство по миграциям

Проект использует [goose](https://github.com/pressly/goose) для управления миграциями базы данных.

## Быстрый старт

1. **Запустите PostgreSQL:**
   ```bash
   docker-compose up -d
   ```

2. **Примените существующие миграции:**
   ```bash
   make migrate
   ```

3. **Проверьте статус:**
   ```bash
   make migrate-status
   ```

## Создание новых миграций

```bash
# Создать новую миграцию
make migrate-create name=add_user_avatar

# Это создаст файл вида: migrations/20240127021234_add_user_avatar.sql
```

Отредактируйте созданный файл:
```sql
-- +goose Up
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);

-- +goose Down
ALTER TABLE users DROP COLUMN avatar_url;
```

Примените миграцию:
```bash
make migrate
```

## Управление миграциями

```bash
# Применить все ожидающие миграции
make migrate

# Откатить последнюю миграцию
make migrate-down

# Показать статус всех миграций
make migrate-status

# Показать текущую версию БД
make migrate-version

# Сбросить все миграции (ОСТОРОЖНО!)
make migrate-reset

# Откатиться до конкретной версии
go run cmd/migrations/main.go down-to 20240127000000
```

## Прямые команды

Можно также использовать команды напрямую:

```bash
# Статус миграций
go run cmd/migrations/main.go status

# Применить миграции
go run cmd/migrations/main.go up

# Создать миграцию
go run cmd/migrations/main.go create add_new_table

# Помощь
go run cmd/migrations/main.go --help
```

## Структура миграции

Все миграции должны содержать блоки `-- +goose Up` и `-- +goose Down`:

```sql
-- +goose Up
-- SQL для применения миграции
CREATE TABLE example (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- +goose Down  
-- SQL для отката миграции
DROP TABLE example;
```

## Важные замечания

- Всегда проверяйте миграции перед применением в продакшене
- Создавайте резервные копии базы данных перед важными изменениями
- Используйте `make migrate-status` для проверки состояния миграций
- Миграции применяются в порядке их timestamp'ов в именах файлов 