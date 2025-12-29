package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Project represents a hackathon project submission
type Project struct {
	ID                int        `json:"id" db:"id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	HackathonID       int        `json:"hackathon_id" db:"hackathon_id"`
	Hackathon         *Hackathon `json:"hackathon,omitempty" belongs_to:"hackathon" fk_id:"hackathon_id"`
	UserID            *uuid.UUID `json:"user_id" db:"user_id"`
	User              *User      `json:"user,omitempty" belongs_to:"user" fk_id:"user_id"`
	Name              string     `json:"name" db:"name"`
	Description       string     `json:"description" db:"description"`
	RepositoryURL     string     `json:"repository_url" db:"repository_url"`
	DemoURL           string     `json:"demo_url" db:"demo_url"`
	Status            string     `json:"status" db:"status"`
	ImageData         []byte     `json:"image_data" db:"image_data"`
	ImageContentType  *string    `json:"image_content_type" db:"image_content_type"`
	Presenting        bool       `json:"presenting" db:"presenting"`
	PresentationOrder *time.Time `json:"presentation_order" db:"presentation_order"`
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
	verrs := validate.NewErrors()
	// Ensure a user is set on creation
	if p.UserID == nil {
		verrs.Add("UserID", "must be present")
		return verrs, nil
	}
	// Enforce one project per user per hackathon
	existing := &Projects{}
	if err := tx.Where("hackathon_id = ? AND user_id = ?", p.HackathonID, p.UserID).Limit(1).All(existing); err != nil {
		return verrs, err
	}
	if len(*existing) > 0 {
		verrs.Add("UserID", "already has a project in this hackathon")
	}
	return verrs, nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (p *Project) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
