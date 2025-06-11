package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/achsanalfitra/gopayslip/hlp"
	"github.com/achsanalfitra/gopayslip/internal/app"
	"github.com/achsanalfitra/gopayslip/internal/model"
	"github.com/achsanalfitra/gopayslip/internal/services/empl"
)

type Admin interface {
	DefinePayroll(userID int64, start, end time.Time, ctx context.Context) error
	RunPayroll(ctx context.Context) (end time.Time, err error)
}

type adminSvcImpl struct{}

func NewAdminServices() Admin {
	return &adminSvcImpl{}
}

// get the service
var emplServices = empl.NewEmplServices()

func (a *adminSvcImpl) DefinePayroll(userID int64, start, end time.Time, ctx context.Context) error {
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
	query := `SELECT id, start_period, end_period, is_run FROM payroll ORDER BY end_period DESC LIMIT 1`
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
	insertQuery := `INSERT INTO payroll (created_by, updated_by, start_period, end_period, created_at, updated_at, is_run) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
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
func (a *adminSvcImpl) RunPayroll(ctx context.Context) (end time.Time, err error) {
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return time.Time{}, err
	}

	// find the latest payroll where statis is_run is false, if not found return error
	var latestUnrunPayroll model.Payroll
	query := `SELECT id, end_period FROM payroll WHERE is_run = FALSE ORDER BY end_period ASC LIMIT 1`
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
	updateQuery := `UPDATE payroll SET is_run = TRUE, updated_at = $1 WHERE id = $2`
	_, err = db.ExecContext(ctx, updateQuery, time.Now(), latestUnrunPayroll.ID)
	if err != nil {
		return time.Time{}, errors.New("failed to update payroll status")
	}

	return latestUnrunPayroll.EndPeriod, nil
}

func (a *adminSvcImpl) GeneratePayrollSummary(ctx context.Context, start, end time.Time) (PayslipList map[string]float64, Total float64, err error) {
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return make(map[string]float64), 0, err
	}

	// instantiate
	PayslipList = make(map[string]float64)
	Total = 0.0

	// get all the user id
	rows, err := db.QueryContext(ctx, `SELECT id FROM users`)
	if err != nil {
		return make(map[string]float64), 0, errors.New("failed to query users")
	}
	defer rows.Close()

	// iterate for all user id
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return make(map[string]float64), 0, fmt.Errorf("failed to scan user ID: %w", err)
		}

		payslip, err := emplServices.GeneratePayslip(userID, ctx, start, end)
		if err != nil {
			log.Printf("failed to generate payslip for user %d: %v", userID, err)
			continue
		}

		PayslipList[strconv.FormatInt(userID, 10)] = payslip.TakeHomePay
		Total += payslip.TakeHomePay
	}

	if err := rows.Err(); err != nil {
		return make(map[string]float64), 0, fmt.Errorf("error during user iteration")
	}

	return PayslipList, Total, nil
}
