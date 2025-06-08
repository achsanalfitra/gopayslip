package internal

import (
	"database/sql"
	"errors"

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

func (a *App) Login(user, pass, role string) error {
	var hashedPassword string

	err := a.DB.QueryRow("SELECT password FROM users WHERE username=$1 and role=$2", user, role).Scan(&hashedPassword)
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
