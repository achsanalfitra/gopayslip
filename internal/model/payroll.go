package model

import (
	"time"
)

type Payroll struct {
	IsRun       bool      `json:"is_run"`
	ID          int64     `json:"id"`
	CreatedBy   int64     `json:"created_by"`
	UpdatedBy   int64     `json:"updated_by"`
	StartPeriod time.Time `json:"start_period"`
	EndPeriod   time.Time `json:"end_period"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
