package models

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Hackathon represents a hackathon event
type Hackathon struct {
	ID          string    `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	EndDate     time.Time `json:"end_date" db:"end_date"`
	Status      string    `json:"status" db:"status"`
	OwnerID     uuid.UUID `json:"owner_id" db:"owner_id"`
	Schedule    string    `json:"schedule" db:"schedule"`
}

// String is not required by pop and may be deleted
func (h Hackathon) String() string {
	jh, _ := json.Marshal(h)
	return string(jh)
}

// Hackathons is not required by pop and may be deleted
type Hackathons []Hackathon

// String is not required by pop and may be deleted
func (h Hackathons) String() string {
	jh, _ := json.Marshal(h)
	return string(jh)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (h *Hackathon) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: h.Title, Name: "Title"},
		&validators.StringIsPresent{Field: h.Description, Name: "Description"},
		&validators.TimeIsPresent{Field: h.StartDate, Name: "StartDate"},
		&validators.TimeIsPresent{Field: h.EndDate, Name: "EndDate"},
		&validators.UUIDIsPresent{Field: h.OwnerID, Name: "OwnerID"},
	), nil
}

// generateUniqueHackathonID generates a unique 12-character hexadecimal string for hackathon ID
func generateUniqueHackathonID() string {
	bytes := make([]byte, 6) // 6 bytes = 12 hex characters
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405")))[:12]
	}
	return hex.EncodeToString(bytes)
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (h *Hackathon) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	if h.ID == "" {
		h.ID = generateUniqueHackathonID()
	}
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (h *Hackathon) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
