package database

import (
	"database/sql"
	"fmt"

	"github.com/Spoloborota/experiment/internal/config"
	_ "github.com/lib/pq"
)

// Connect создает подключение к PostgreSQL
func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
