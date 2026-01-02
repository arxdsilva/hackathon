package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// File represents an uploaded file
type File struct {
	ID          string     `json:"id" db:"id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Filename    string     `json:"filename" db:"filename"`
	Data        []byte     `json:"data" db:"data"`
	ContentType string     `json:"content_type" db:"content_type"`
	Size        int        `json:"size" db:"size"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	User        *User      `json:"user,omitempty" belongs_to:"user" fk_id:"user_id"`
	HackathonID *string    `json:"hackathon_id" db:"hackathon_id"`
	Hackathon   *Hackathon `json:"hackathon,omitempty" belongs_to:"hackathon" fk_id:"hackathon_id"`
	ProjectID   *string    `json:"project_id" db:"project_id"`
	Project     *Project   `json:"project,omitempty" belongs_to:"project" fk_id:"project_id"`
}

// String is not required by pop and may be deleted
func (f File) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Files is not required by pop and may be deleted
type Files []File

// String is not required by pop and may be deleted
func (f Files) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (f *File) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: f.Filename, Name: "Filename"},
		&validators.StringIsPresent{Field: f.ContentType, Name: "ContentType"},
		&validators.IntIsPresent{Field: int(f.Size), Name: "Size"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (f *File) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	if f.ID == "" {
		f.ID = generateUniqueID()
	}
	return validate.NewErrors(), nil
}
