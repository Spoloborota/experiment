package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/Spoloborota/experiment/internal/config"
)

const (
	dialect       = "postgres"
	migrationsDir = "migrations"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", migrationsDir, "directory with migration files")
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к базе данных
	db, err := sql.Open(dialect, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Выполняем команду goose
	switch command {
	case "create":
		if len(args) < 2 {
			fmt.Println("Usage: go run cmd/migrations/main.go create MIGRATION_NAME")
			return
		}
		if err := goose.Create(db, *dir, args[1], "sql"); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		return
	}

	if err := goose.SetDialect(dialect); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	switch command {
	case "up":
		if err := goose.Up(db, *dir); err != nil {
			log.Fatalf("Failed to run migrations up: %v", err)
		}
	case "up-by-one":
		if err := goose.UpByOne(db, *dir); err != nil {
			log.Fatalf("Failed to run migration up by one: %v", err)
		}
	case "down":
		if err := goose.Down(db, *dir); err != nil {
			log.Fatalf("Failed to run migration down: %v", err)
		}
	case "down-to":
		if len(args) < 2 {
			fmt.Println("Usage: go run cmd/migrations/main.go down-to VERSION")
			return
		}
		version, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := goose.DownTo(db, *dir, version); err != nil {
			log.Fatalf("Failed to run migration down-to %d: %v", version, err)
		}
	case "redo":
		if err := goose.Redo(db, *dir); err != nil {
			log.Fatalf("Failed to redo migration: %v", err)
		}
	case "reset":
		if err := goose.Reset(db, *dir); err != nil {
			log.Fatalf("Failed to reset migrations: %v", err)
		}
	case "status":
		if err := goose.Status(db, *dir); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	case "version":
		version, err := goose.GetDBVersion(db)
		if err != nil {
			log.Fatalf("Failed to get database version: %v", err)
		}
		fmt.Printf("Current database version: %d\n", version)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		flags.Usage()
	}
}

func usage() {
	fmt.Print(`Usage: go run cmd/migrations/main.go COMMAND [ARGS...]

Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql]    Creates new migration file with the current timestamp

Examples:
    go run cmd/migrations/main.go up
    go run cmd/migrations/main.go create add_users_table
    go run cmd/migrations/main.go status
    go run cmd/migrations/main.go down-to 20200101000000

Options:
`)
	flags.PrintDefaults()
}
