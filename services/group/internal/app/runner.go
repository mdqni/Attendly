package app

import (
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(connString string, migrationsPath string) {
	var m *migrate.Migrate
	var err error

	for i := 0; i < 10; i++ {
		m, err = migrate.New("file://"+migrationsPath, connString)
		if err == nil {
			break
		}
		log.Printf("Migration tool not ready yet (%d/10), retrying in 3s...: %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
