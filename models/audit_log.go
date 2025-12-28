package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// AuditLog represents an audit log entry for admin actions.
type AuditLog struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	UserID       *uuid.UUID `db:"user_id" json:"user_id"`
	Action       string     `db:"action" json:"action"`
	ResourceType string     `db:"resource_type" json:"resource_type"`
	ResourceID   *string    `db:"resource_id" json:"resource_id"`
	Details      string     `db:"details" json:"details"`
	IPAddress    string     `db:"ip_address" json:"ip_address"`
	UserAgent    string     `db:"user_agent" json:"user_agent"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

// String returns the JSON representation of the audit log.
func (a AuditLog) String() string {
	ju, _ := json.Marshal(a)
	return string(ju)
}

// AuditLogs is a collection of AuditLog.
type AuditLogs []AuditLog

// String returns the JSON representation of the audit logs.
func (a AuditLogs) String() string {
	ju, _ := json.Marshal(a)
	return string(ju)
}
