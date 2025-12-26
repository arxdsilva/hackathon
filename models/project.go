package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// Project represents a hackathon project submission
type Project struct {
	ID            int       `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	HackathonID   int       `json:"hackathon_id" db:"hackathon_id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	RepositoryURL string    `json:"repository_url" db:"repository_url"`
	DemoURL       string    `json:"demo_url" db:"demo_url"`
	Status        string    `json:"status" db:"status"`
	TeamMembers   int       `json:"team_members" db:"team_members"`
}

// String is not required by pop and may be deleted
func (p Project) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Projects is not required by pop and may be deleted
type Projects []Project

// String is not required by pop and may be deleted
func (p Projects) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *Project) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Name, Name: "Name"},
		&validators.StringIsPresent{Field: p.Description, Name: "Description"},
		&validators.IntIsPresent{Field: p.HackathonID, Name: "HackathonID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (p *Project) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (p *Project) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
