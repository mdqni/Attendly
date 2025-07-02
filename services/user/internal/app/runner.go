package app

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"os"
	"path/filepath"
)

func RunMigrations(connString string) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}
	migrationPath := filepath.Join(wd, "services/user/internal/migrations")
	if err := goose.Up(db, migrationPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}
