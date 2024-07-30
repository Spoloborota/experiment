package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Migrate(db *pgxpool.Pool) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            first_name VARCHAR(50),
            last_name VARCHAR(50),
            age INT,
            gender VARCHAR(10),
            interests TEXT,
            city VARCHAR(50),
            password_hash TEXT
        )`,
	}

	for _, query := range queries {
		_, err := db.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Failed to execute migration: %v\n", err)
		}
	}
}
