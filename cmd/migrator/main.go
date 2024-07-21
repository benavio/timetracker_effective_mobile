package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbPath := os.Getenv("DATABASE_PATH")
	migrationsPath := os.Getenv("DATABASE_MIGRATIONS_PATH")

	if dbPath == "" || migrationsPath == "" {
		log.Fatalf("env not set")
	}

	log.Printf("Connecting to database at %s", dbPath)

	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		log.Fatalf("%s : %s", "can't open db", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Can't connect to database: %s", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("%s : %s", "can't init driver", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		log.Fatalf("%s : %s", "can't init migrator", err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no changes")
			return
		}
	}

	log.Println("migrations done")
}
