package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ProjectMembership represents a user joining a project
type ProjectMembership struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	ProjectID string   `json:"project_id" db:"project_id"`
	Project   *Project `json:"project,omitempty" belongs_to:"project"`

	UserID uuid.UUID `json:"user_id" db:"user_id"`
	User   *User     `json:"user,omitempty" belongs_to:"user"`
}

// String is not required by pop and may be deleted
func (m ProjectMembership) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// ProjectMemberships is not required by pop and may be deleted
type ProjectMemberships []ProjectMembership

// String is not required by pop and may be deleted
func (m ProjectMemberships) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate runs on Validate* calls
func (m *ProjectMembership) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.ProjectID, Name: "ProjectID"},
		&validators.UUIDIsPresent{Field: m.UserID, Name: "UserID"},
	), nil
}
