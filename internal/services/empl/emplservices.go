package empl

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/achsanalfitra/gopayslip/hlp"
	"github.com/achsanalfitra/gopayslip/internal/app"
)

type Empl interface {
	GeneratePayslip(userID int64, ctx context.Context, start, end time.Time) (Payslip, error)
}

type emplImplementation struct{}

func NewEmplServices() Empl {
	return &emplImplementation{}
}

// define payslip for easier payload distribution
type Payslip struct {
	UserID           int64
	Attendance       int
	Overtime         float64
	TakeHomePay      float64
	Salary           float64
	AttendancePay    float64
	OvertimePay      float64
	ReimbursementPay float64
	CreatedAt        time.Time
	PayrollStart     time.Time
	PayrollEnd       time.Time
}

func (e *emplImplementation) GeneratePayslip(userID int64, ctx context.Context, start, end time.Time) (Payslip, error) {
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return Payslip{}, err
	}

	// implement the private functions
	attendance, totalWorkingDays, err := e.countAttendance(userID, db, ctx, start, end)
	if err != nil {
		return Payslip{}, errors.New("failed to count attendance")
	}

	totalReimb, err := e.totalReimbursement(userID, db, ctx, start, end)
	if err != nil {
		return Payslip{}, errors.New("failed to get total reimbursement")
	}

	overtimeHrs, err := e.overtimeDuration(userID, db, ctx, start, end)
	if err != nil {
		return Payslip{}, errors.New("failed to get overtime duration")
	}

	baseSalary, err := e.getUserSalary(userID, db, ctx)
	if err != nil {
		return Payslip{}, errors.New("failed to get user salary")
	}

	// business logic calculation
	hourlyRate := baseSalary / (float64(totalWorkingDays) * 8) // assume 9 to 5 is 8 hours workday

	attendancePay := baseSalary
	if totalWorkingDays > 0 {
		attendancePay = baseSalary * (float64(attendance) / float64(totalWorkingDays))
	}

	overtimePay := hourlyRate * 2 * overtimeHrs

	takeHomePay := attendancePay + overtimePay + totalReimb

	// populate payslip payload
	payslip := Payslip{
		UserID:           userID,
		Attendance:       attendance,
		Overtime:         overtimeHrs,
		TakeHomePay:      takeHomePay,
		Salary:           baseSalary,
		AttendancePay:    attendancePay,
		OvertimePay:      overtimePay,
		ReimbursementPay: totalReimb,
		CreatedAt:        time.Now(),
		PayrollStart:     start,
		PayrollEnd:       end,
	}

	return payslip, nil
}

func (e *emplImplementation) getUserSalary(userID int64, db *sql.DB, ctx context.Context) (salary float64, err error) {
	query := `SELECT salary FROM users WHERE id = $1`
	err = db.QueryRowContext(ctx, query, userID).Scan(&salary)
	if err == sql.ErrNoRows {
		return 0, errors.New("user not found or salary not defined")
	}
	if err != nil {
		return 0, errors.New("failed to query user salary")
	}
	return salary, nil
}

func (e *emplImplementation) countAttendance(userID int64, db *sql.DB, ctx context.Context, start, end time.Time) (attendance int, totalWorkingDays int, err error) {
	query := `SELECT COUNT(*) FROM attendance WHERE user_id = $1 AND created_at BETWEEN $2 AND $3`
	err = db.QueryRowContext(ctx, query, userID, start, end).Scan(&attendance)
	if err != nil {
		return 0, 0, errors.New("failed to count attendance")
	}

	// calculate working days that doesn't include weekend
	totalWorkingDays = 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			totalWorkingDays++
		}
	}

	return attendance, totalWorkingDays, nil
}

func (e *emplImplementation) totalReimbursement(userID int64, db *sql.DB, ctx context.Context, start, end time.Time) (total float64, err error) {
	query := `SELECT COALESCE(SUM(reimbursement_amount), 0) FROM reimbursement WHERE user_id = $1 AND created_at BETWEEN $2 AND $3`
	err = db.QueryRowContext(ctx, query, userID, start, end).Scan(&total)
	if err != nil {
		return 0, errors.New("failed to sum reimbursements")
	}
	return total, nil
}

func (e *emplImplementation) overtimeDuration(userID int64, db *sql.DB, ctx context.Context, start, end time.Time) (totalHours float64, err error) {
	query := `SELECT COALESCE(SUM(overtime_duration), 0) FROM overtime WHERE employee_id = $1 AND created_at BETWEEN $2 AND $3`
	err = db.QueryRowContext(ctx, query, userID, start, end).Scan(&totalHours)
	if err != nil {
		return 0, errors.New("failed to sum overtime duration")
	}
	return totalHours, nil
}
