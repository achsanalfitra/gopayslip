package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/achsanalfitra/gopayslip/internal/app"
	"github.com/achsanalfitra/gopayslip/internal/config"
	"github.com/achsanalfitra/gopayslip/internal/router"
)

func main() {
	// init DB
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.DB.Ping(); err != nil {
		log.Fatalf("can't connect to database: %s", err)
	}

	serverConfig := config

	appConfig := app.AppConfig{
		DB:         db.DB,
		InitStates: make(map[string]interface{}),
		Server: config.CreateServer(
			addr,
			router.NewRouter().ServeHTTP(),
		),
	}

	// get the initial states
	var startPeriod, endPeriod time.Time
	query := `SELECT start_period, end_period FROM payroll ORDER BY end_period DESC LIMIT 1`

	err = appConfig.DB.QueryRowContext(context.Background(), query).Scan(&startPeriod, &endPeriod)

	if err == sql.ErrNoRows {
		log.Println("No existing payroll periods found in the database. InitStates for payroll dates will be default (zero time).")
	} else if err != nil {
		log.Fatalf("Failed to query latest payroll period for InitStates: %v", err)
	} else {
		appConfig.InitStates[string(router.CtxStartKey)] = startPeriod
		appConfig.InitStates[string(router.CtxEndKey)] = endPeriod
		log.Printf("InitStates populated with latest payroll period: Start=%v, End=%v", startPeriod, endPeriod)
	}

	// get InitStates

	appConfig := app.AppConfig{
		DB: db.DB,
	}
}
