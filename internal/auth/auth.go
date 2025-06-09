package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/achsanalfitra/gopayslip/internal/app"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"user"`
	Password string `json:"pass"`
	Role     string `json:"role"`
}

type Auth interface {
	Login(user, pass, role string) error
	Register(user, pass string)
	GetToken(refreshToken string) (accessToken string)
	AllowAccess() (user string, err error)
}

func Login(user, pass, role string, ctx context.Context) error {
	var hashedPassword string

	db := app.GetDB(ctx)

	err := db.QueryRow("SELECT password FROM users WHERE username=$1 and role=$2", user, role).Scan(&hashedPassword)
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
