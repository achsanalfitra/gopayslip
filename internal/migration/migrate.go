package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"time"
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
	stmt, err := os.ReadFile(filepath.Join(m.src, m.initSchema))
	if err != nil {
		return fmt.Errorf("fail reading schema migration: can't find file named %s", m.initSchema)
	}

	if _, err := m.db.Exec(string(stmt)); err != nil {
		return fmt.Errorf("can't migrate schema_migration from %s: ensure the SQL format is correct", m.initSchema)
	}

	return nil
}

func (m *Migrate) Up() error {
	// get the sorted files
	versQ, err := m.hlpUp()
	if err != nil {
		return err
	}

	// IMPORTANT: this assumes ordered migration, e.g., 001, 002, 003 for simplicity. Not skipping 001, 003 et
	// cleanup versQ against existing data
	var pendQ []string
	for _, name := range versQ {
		var foundMatch string
		err := m.db.QueryRow("SELECT schema FROM schema_migration WHERE schema=$1", name).Scan(&foundMatch)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				pendQ = append(pendQ, name)
			} else {
				return err
			}
		}
	}

	// start the transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed starting the transaction")
	}
	defer tx.Rollback()

	// get statement
	for _, name := range pendQ {
		stmt, err := os.ReadFile(filepath.Join(m.src, name))
		if err != nil {
			return fmt.Errorf("file can't be read")
		}

		if _, err := tx.Exec(string(stmt)); err != nil {
			return fmt.Errorf("can't execute the statement in %s", name)
		}

		if _, err := tx.Exec("INSERT INTO schema_migration (schema, created_at) VALUES ($1, $2)", name, time.Now()); err != nil {
			return fmt.Errorf("can't insert %s into schema migration", name)
		}
	}

	// finally, commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed")
	}

	return nil
}

func (m *Migrate) hlpUp() ([]string, error) {
	// parse all files that ends with *_up
	pattern := fmt.Sprintf("%s/*_up.sql", m.src)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse files matching the pattern *_up.sql in %s due to permission issues", m.src)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("can't find files matching the pattern *_up.sql in %s", m.src)
	}

	// instantiate regex check instance
	regex := regexp.MustCompile(`(\d+)_.*_up\.sql$`)

	// order the file
	var versQ []string
	uniqueSet := make(map[string]struct{})
	for _, file := range files {
		name := filepath.Base(file)

		// catch filename validity with regex
		if match := regex.MatchString(name); !match {
			return nil, fmt.Errorf("invalid sql file name for %s", name)
		}

		// find matching substring, in this case we want the first one
		vers := regex.FindStringSubmatch(name)

		versQ = append(versQ, name)
		if _, exists := uniqueSet[vers[1]]; exists {
			return nil, fmt.Errorf("can't have the same sequence %s which is found in %s", vers[1], name)
		}

		uniqueSet[vers[1]] = struct{}{}
	}

	slices.Sort(versQ)
	return versQ, nil
}

// TO BE IMPLEMENTED: which might be unnecessary for a prototype
func (m *Migrate) Down() {
	// parse all files that ends with *_down
	// check the "up version" in the database, if exists delete it and run the respective down.sql
}
