package main

import (
	"log"

	"github.com/achsanalfitra/gopayslip/internal/config"
	"github.com/achsanalfitra/gopayslip/internal/migration"
)

func main() {
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.DB.Ping(); err != nil {
		log.Fatalf("can't connect to database: %s", err)
	}

	// for demo project, we will only demosntrate the UP function
	m := migration.NewMigration(
		string(migration.UP),
		"internal/migration/migrations/",
		"schema_migration.sql",
		db.DB,
	)

	if err := m.InitMigrationSchema(); err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	db.DB.Close()
}
