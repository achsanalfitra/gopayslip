package model

import (
	"time"

	"github.com/google/uuid"
)

type Overtime struct {
	ID        int64         `json:"id"`
	UserID    int64         `json:"user_id"`
	CreatedBy int64         `json:"created_by"`
	UpdatedBy int64         `json:"updated_by"`
	RequestId uuid.UUID     `json:"request_id"`
	Interval  time.Duration `json:"overtime_duration"`
	Date      time.Time     `json:"overtime_date"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
