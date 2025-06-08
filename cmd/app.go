package cmd

import "database/sql"

// create App for dependency injection
type App struct {
	DB *sql.DB
}
