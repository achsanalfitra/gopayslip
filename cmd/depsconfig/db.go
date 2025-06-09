package depsconfig

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	// practice ordering from small -> large
	port     string
	host     string
	sslmode  string
	user     string
	password string
	db       string
	DB       *sql.DB
}

func InitDatabase() (*Database, error) {
	database := Database{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		db:       os.Getenv("POSTGRES_DB"),
		port:     os.Getenv("POSTGRES_PORT"),
		host:     os.Getenv("POSTGRES_HOST"),
		sslmode:  os.Getenv("SSLMODE"),
	}

	// check .env completeness
	if database.user == "" || database.password == "" || database.db == "" || database.port == "" || database.host == "" {
		return nil, fmt.Errorf("missing one or more required PostgreSQL environment variables (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT, POSTGRES_HOST)")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=%s",
		database.user, database.password, database.db, database.port, database.host, database.sslmode)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("cannot establish connection to database: %w", err)
	}

	database.DB = dbConn

	return &database, nil
}
