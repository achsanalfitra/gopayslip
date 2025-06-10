package migration

import (
	"database/sql"
	"fmt"
	"os"
)

type Flag string

const (
	UP   Flag = "up"
	DOWN Flag = "down"
)

type Migrate struct {
	flag       string
	src        string // migrations source files path
	initSchema string // your own schema_migration for tracking migration versioning
	db         *sql.DB
}

func NewMigration(flag, src, initSchema string, db *sql.DB) *Migrate {
	return &Migrate{
		flag:       flag,
		src:        src,
		initSchema: initSchema,
		db:         db,
	}
}

func (m *Migrate) InitMigrationSchema() error {
	// reads the migration_schema_init.sql
	stmt, err := os.ReadFile(m.initSchema)
	if err != nil {
		return fmt.Errorf("fail reading schema migration: can't find file named %s", m.initSchema)
	}

	if _, err := m.db.Exec(string(stmt)); err != nil {
		return fmt.Errorf("can't migrate schema_migration from %s: ensure the SQL format is correct", m.initSchema)
	}

	return nil
}

func (m *Migrate) Up() {
	// parse all files that ends with *_up
	// check them against the database
	// run it if there is an untracked migration in the migration_schema table
}

func (m *Migrate) Down() {
	// parse all files that ends with *_down
	// check the "up version" in the database, if exists delete it and run the respective down.sql
}
