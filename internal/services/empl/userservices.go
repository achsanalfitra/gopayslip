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
