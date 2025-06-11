package app

import (
	"database/sql"
	"os"

	"github.com/achsanalfitra/gopayslip/internal/config"
)

// app-wide key consistency
type DBKey string

// InitDB precedes this, so this is never empty string because the error is catch earlier
var postgresKey DBKey = DBKey(os.Getenv("POSTGRES_DB"))

var (
	PQ DBKey = postgresKey
	// insert other databases here
)

type AppConfig struct {
	DB         *sql.DB
	Server     *config.Server
	InitStates map[string]any
}

// create App for dependency injection
type App struct {
	DB         *sql.DB
	Server     *config.Server
	InitStates map[string]any // data init
	// declare other app-dependencies here
}

func NewApp(cfg AppConfig) *App {
	return &App{
		DB:         cfg.DB,
		Server:     cfg.Server,
		InitStates: make(map[string]any),
		// don't forget to instantiate them
	}
}

func (a *App) Run() {
	a.Server.Start()
}
