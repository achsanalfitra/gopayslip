package hlp

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/achsanalfitra/gopayslip/internal/app"
)

func GetDB(ctx context.Context, dbKey app.DBKey) (*sql.DB, error) {
	if db, ok := ctx.Value(dbKey).(*sql.DB); ok {
		return db, nil
	}
	log.Println("error reaching database: ensure proper DBKey is used and connection is valid")
	return nil, errors.New("database can't be reached")
}
