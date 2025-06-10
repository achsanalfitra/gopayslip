package model

import (
	"time"

	"github.com/google/uuid"
)

type ActionType string

const (
	CREATE ActionType = "CREATE"
	UPDATE ActionType = "UPDATE"
	READ   ActionType = "READ"
	DELETE ActionType = "DELETE"
)

type AuditLog struct {
	ID               int64      `json:"id"`
	CreatedBy        int64      `json:"created_by"`
	AffectedRecordID int64      `json:"affected_record_id"`
	RequestId        uuid.UUID  `json:"request_id"`
	ActionType       ActionType `json:"action_type"`
	EventType        string     `json:"event_type"` // always implement this using a strict typing
	AffectedRecord   Table      `json:"affected_table"`
	OldData          string     `json:"old_data"`
	NewData          string     `json:"new_data"`
	IPAddress        string     `json:"ip_address"`
	CreatedAt        time.Time  `json:"created_at"`
}
