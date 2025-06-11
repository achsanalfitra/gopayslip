package empl

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/achsanalfitra/gopayslip/hlp"
	"github.com/achsanalfitra/gopayslip/internal/app"
	"github.com/achsanalfitra/gopayslip/internal/model"
	"github.com/google/uuid"
)

type User interface {
	CheckIn(userID int64, requestID uuid.UUID, ctx context.Context) error
	ProposeOvertime()
}

func CheckIn(userID int64, requestID uuid.UUID, ctx context.Context) error {
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	// scan for user role from user id
	var userRole string
	query := `SELECT role FROM users WHERE id = $1`
	err = db.QueryRowContext(ctx, query, userID).Scan(&userRole)
	if err == sql.ErrNoRows {
		return errors.New("user not found for check-in")
	}
	if err != nil {
		return errors.New("user query error while check-in")
	}

	var assertedRole model.Role
	if userRole != "employee" {
		assertedRole = model.ADMIN
	} else {
		assertedRole = model.EMPLOYEE
	}

	attendanceRecord := model.Attendance{
		UserID:    userID,
		CreatedBy: userID,
		UpdatedBy: userID,
		RequestId: requestID,
		UserRole:  assertedRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertQuery := `INSERT INTO attendance (user_id, created_by, updated_by, request_id, user_role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.ExecContext(ctx, insertQuery,
		attendanceRecord.UserID,
		attendanceRecord.CreatedBy,
		attendanceRecord.UpdatedBy,
		attendanceRecord.RequestId,
		attendanceRecord.UserRole,
		attendanceRecord.CreatedAt,
		attendanceRecord.UpdatedAt,
	)
	if err != nil {
		return errors.New("failed to insert attendance record")
	}

	return nil
}

func ProposeOvertime(userID int64, requestID uuid.UUID, overtimeDuration time.Duration, overtimeDate time.Time, ctx context.Context) error {
	// early exit when overtme duration > 3 hours
	if overtimeDuration > 3*time.Hour {
		return errors.New("maximum overtime is 3 hours")
	}

	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	// validate after work
	var latestAttendanceDate time.Time
	query := `SELECT overtime_date FROM attendance WHERE user_id = $1 ORDER BY created_date DESC LIMIT 1`
	err = db.QueryRowContext(ctx, query, userID).Scan(&latestAttendanceDate)
	if err == sql.ErrNoRows {
		return errors.New("no attendance record found for user, cannot propose overtime")
	}
	if err != nil {
		return errors.New("failed to query attendance during overtime proposal")
	}

	// check the latest attendance, if proposed date is lesser than the day the latest attendance happen on 5 PM, invalidate
	workEndTime := time.Date(latestAttendanceDate.Year(), latestAttendanceDate.Month(), latestAttendanceDate.Day(), 17, 0, 0, 0, latestAttendanceDate.Location()) // 5 PM on attendance day
	if overtimeDate.Before(workEndTime) {
		return errors.New("overtime can only be proposed after 5 PM on the day on working day")
	}

	// if overtime for that day already exists, prevent another overtime
	var existingOvertimeID int64
	queryExistingOvertime := `SELECT id FROM overtimes WHERE employee_id = $1 AND date = $2`
	err = db.QueryRowContext(ctx, queryExistingOvertime, userID, time.Now().Truncate(24*time.Hour)).Scan(&existingOvertimeID)
	if err != nil && err != sql.ErrNoRows {
		return errors.New("failed to check existing overtime")
	}
	if err == nil {
		return errors.New("overtime for this day already exists for this user")
	}

	// post the overtime payload
	overtimePayload := model.Overtime{
		UserID:    userID,
		CreatedBy: userID,
		UpdatedBy: userID,
		RequestId: requestID,
		Interval:  overtimeDuration,
		Date:      overtimeDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	insertQuery := `INSERT INTO overtimes (user_id, created_by, updated_by, request_id, overtime_duration, overtime_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.ExecContext(ctx, insertQuery,
		overtimePayload.UserID,
		overtimePayload.CreatedBy,
		overtimePayload.UpdatedBy,
		overtimePayload.RequestId,
		overtimePayload.Interval,
		overtimePayload.Date,
		overtimePayload.CreatedAt,
		overtimePayload.UpdatedAt,
	)

	if err != nil {
		return errors.New("failed to insert overtime record")
	}

	return nil
}

func ProposeReimbursement(userID int64, requestID uuid.UUID, amount float64, desc string, ctx context.Context) error {
	// invalidate minus amount, fail fast
	if amount <= 0 {
		return errors.New("reimbursement can't be smaller than 0")
	}

	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	// post the whole payload immediately
	reimbursementPayload := model.Reimbursement{
		UserID:              userID,
		CreatedBy:           userID,
		UpdatedBy:           userID,
		ReimbursementAmount: amount,
		RequestId:           requestID,
		Description:         desc,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	insertQuery := `INSERT INTO reimbursements (user_id, created_by, updated_by, amount, request_id, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.ExecContext(ctx, insertQuery,
		reimbursementPayload.UserID,
		reimbursementPayload.CreatedBy,
		reimbursementPayload.UpdatedBy,
		reimbursementPayload.ReimbursementAmount,
		reimbursementPayload.RequestId,
		reimbursementPayload.Description,
		reimbursementPayload.CreatedAt,
		reimbursementPayload.UpdatedAt,
	)
	if err != nil {
		return errors.New("failed to insert reimbursement record")
	}

	return nil

}
