package app

import (
	"context"
	"database/sql"
)

// config
type ctxKey string

const dbKey ctxKey = "db"

// create App for dependency injection
type App struct {
	DB *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{
		DB: db,
	}
}

func InjectDB(ctx context.Context, a *App) context.Context {
	return context.WithValue(ctx, dbKey, a.DB)
}

func GetDB(ctx context.Context) *sql.DB {
	if db, ok := ctx.Value(dbKey).(*sql.DB); ok {
		return db
	}
	panic("database not found") // development error
}
