package app

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations(connString string) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close DB: %v", err)
		}
	}(db)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}
	migrationPath := filepath.Join(wd, "services/group/internal/migrations")
	log.Printf("Running migrations: %v", migrationPath)
	if err := goose.Up(db, migrationPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}
