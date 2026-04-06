package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	m, err := migrate.New(
		"file:///app/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("migration init error:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration up error:", err)
	}

	log.Println("migrations applied")
}