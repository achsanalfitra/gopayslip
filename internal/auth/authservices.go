package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/achsanalfitra/gopayslip/hlp"
	"github.com/achsanalfitra/gopayslip/internal/app"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(user, pass, role string, ctx context.Context) error
	Register(user, pass, role string, ctx context.Context) error
}

func Login(user, pass, role string, ctx context.Context) error {
	var hashedPassword string

	// connect to database
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	err = db.QueryRowContext(ctx, "SELECT password FROM users WHERE username=$1 and role=$2", user, role).Scan(&hashedPassword)
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

func Register(user, pass, role string, ctx context.Context) error {
	// connect to database
	db, err := hlp.GetDB(ctx, app.PQ)
	if err != nil {
		return err
	}

	// check if user already exists, this early exit increases performance
	var existingUser string
	err = db.QueryRowContext(ctx, "SELECT username FROM users WHERE username=$1 and role=$2", user, role).Scan(&existingUser)
	if err == nil {
		return errors.New("user already exists")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashedPassword := string(hashedPasswordBytes)

	_, err = db.ExecContext(ctx, "INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", user, hashedPassword, role)
	if err != nil {
		return errors.New("failed to register user")
	}

	return nil
}
