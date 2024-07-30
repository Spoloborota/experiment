package config

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DB *pgxpool.Pool
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	dbUser := viper.GetString("db.user")
	dbPassword := viper.GetString("db.password")
	dbName := viper.GetString("db.name")
	dbHost := viper.GetString("db.host")

	dbUrl := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":5432/" + dbName + "?sslmode=disable"

	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &Config{
		DB: db,
	}
}
