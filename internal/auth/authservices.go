package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/achsanalfitra/gopayslip/hlp"
	"github.com/achsanalfitra/gopayslip/internal/app"
	"github.com/achsanalfitra/gopayslip/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(user, pass, role string, ctx context.Context) error
	Register(user, pass, role string, salary float64, ctx context.Context) error
}

func Login(user, pass, role string, ctx context.Context) error {
	var hashedPassword string

	// connect to database
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	err = db.QueryRowContext(ctx, "SELECT password FROM users WHERE username=$1 and user_role=$2", user, role).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass)); err != nil {
		return errors.New("invalid password")
	}

	return nil
}

func Register(user, pass, role string, salary float64, ctx context.Context) error {
	// connect to database
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	// check if user already exists, this early exit increases performance
	var tempID int64
	err = db.QueryRowContext(ctx, "SELECT id FROM users WHERE username=$1", user).Scan(&tempID)
	if err == nil {
		return errors.New("user already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.New("database query error")
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashedPassword := string(hashedPasswordBytes)

	// this services asks the whole model

	// populate the data
	createdAt := time.Now()

	var initialCreatedUpdatedBy int64 = 0

	userToInsert := model.User{
		Username:  user,
		Password:  hashedPassword,
		UserRole:  model.Role(role),
		Salary:    salary,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		CreatedBy: initialCreatedUpdatedBy, // initial placeholder
		UpdatedBy: initialCreatedUpdatedBy,
	}

	insertQuery := `INSERT INTO users (username, password, user_role, salary, created_at, updated_at, created_by, updated_by)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var newUserID int64
	err = db.QueryRowContext(
		ctx, insertQuery,
		userToInsert.Username,
		userToInsert.Password,
		userToInsert.UserRole,
		userToInsert.Salary,
		userToInsert.CreatedAt,
		userToInsert.UpdatedAt,
		userToInsert.CreatedBy,
		userToInsert.UpdatedBy,
	).Scan(&newUserID)

	if err != nil {
		return errors.New("failed to insert user")
	}

	_, err = db.ExecContext(ctx, "UPDATE users SET created_by=$1, updated_by=$1 WHERE id=$2", newUserID, newUserID)
	if err != nil {
		return errors.New("failed to update created_by/updated_by for new user")
	}

	return nil
}
