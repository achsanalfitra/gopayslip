package app

import (
	"database/sql"

	"github.com/achsanalfitra/gopayslip/cmd/depsconfig"
)

// app-wide key consistency
type DBKey string

const (
	PQ DBKey = "PostgresSQL"
	// insert other databases here
)

type AppConfig struct {
	DB     *sql.DB
	Server *depsconfig.Server
}

// create App for dependency injection
type App struct {
	DB     *sql.DB
	Server *depsconfig.Server
	// declare other app-dependencies here
}

func NewApp(cfg AppConfig) *App {
	return &App{
		DB:     cfg.DB,
		Server: cfg.Server,
		// don't forget to instantiate them
	}
}

func (a *App) Run() {
	a.Server.Start()
}
