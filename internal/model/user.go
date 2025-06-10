package model

import (
	"time"
)

type Role string

const (
	ADMIN    Role = "ADMIN"
	EMPLOYEE Role = "EMPLOYEE"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
	Salary    float64   `json:"salary"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	UserRole  Role      `json:"user_role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
