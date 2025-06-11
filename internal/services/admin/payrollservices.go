package admin

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/achsanalfitra/gopayslip/hlp"
	"github.com/achsanalfitra/gopayslip/internal/app"
	"github.com/achsanalfitra/gopayslip/internal/model"
)

type Admin interface {
}

func DefinePayroll(userID int64, start, end time.Time, ctx context.Context) error {
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	// start period can't be after end period
	if start.After(end) {
		return errors.New("start period cannot be after end period")
	}

	// check interval with the latest payroll
	// if start < than latest end exit where is run == true
	var latestPayroll model.Payroll
	var tempStatus bool
	query := `SELECT id, start_period, end_period, is_run FROM payrolls ORDER BY end_period DESC LIMIT 1`
	err = db.QueryRowContext(ctx, query).Scan(
		&latestPayroll.ID,
		&latestPayroll.StartPeriod,
		&latestPayroll.EndPeriod,
		&tempStatus,
	)

	if err == nil {
		if start.Before(latestPayroll.EndPeriod) && tempStatus {
			return errors.New("new payroll period overlaps with a previously run payroll")
		}

		// check if the latest payroll is not run
		if !tempStatus {
			return errors.New("previous payroll period has not been run yet")
		}
	}

	if err != sql.ErrNoRows {
		return errors.New("failed to query payroll row during validation")
	}

	// populate the model
	payroll := model.Payroll{
		CreatedBy:   userID,
		UpdatedBy:   userID,
		StartPeriod: start,
		EndPeriod:   end,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsRun:       false,
	}

	// insert the payload to db
	insertQuery := `INSERT INTO payrolls (created_by, updated_by, start_period, end_period, created_at, updated_at, is_run) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = db.QueryRowContext(ctx, insertQuery,
		payroll.CreatedBy,
		payroll.UpdatedBy,
		payroll.StartPeriod,
		payroll.EndPeriod,
		payroll.CreatedAt,
		payroll.UpdatedAt,
		payroll.IsRun,
	).Scan(&payroll.ID)

	if err != nil {
		return errors.New("failed to insert payroll")
	}

	return nil
}

// run payroll
func RunPayroll(ctx context.Context) (end time.Time, err error) {
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	// find the latest payroll where statis is_run is false, if not found return error
	var latestUnrunPayroll model.Payroll
	query := `SELECT id, end_period FROM payrolls WHERE is_run = FALSE ORDER BY end_period ASC LIMIT 1`
	err = db.QueryRowContext(ctx, query).Scan(
		&latestUnrunPayroll.ID,
		&latestUnrunPayroll.EndPeriod,
	)

	if err == sql.ErrNoRows {
		return time.Time{}, errors.New("no pending payroll to run")
	}
	if err != nil {
		return time.Time{}, errors.New("failed to run payroll")
	}

	// insert
	updateQuery := `UPDATE payrolls SET is_run = TRUE, updated_at = $1 WHERE id = $2`
	_, err = db.ExecContext(ctx, updateQuery, time.Now(), latestUnrunPayroll.ID)
	if err != nil {
		return time.Time{}, errors.New("failed to update payroll status")
	}

	return latestUnrunPayroll.EndPeriod, nil
}
