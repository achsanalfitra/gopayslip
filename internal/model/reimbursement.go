package model

import (
	"time"

	"github.com/google/uuid"
)

type Reimbursement struct {
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	CreatedBy           int64     `json:"created_by"`
	UpdatedBy           int64     `json:"updated_by"`
	ReimbursementAmount float64   `json:"reimbursement_amount"`
	RequestId           uuid.UUID `json:"request_id"`
	Description         string    `json:"description"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
