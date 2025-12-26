package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User represents a registered account in the system.
type User struct {
    ID                   uuid.UUID `db:"id" json:"id"`
    CreatedAt            time.Time `db:"created_at" json:"created_at"`
    UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
    Email                string    `db:"email" json:"email"`
    PasswordHash         string    `db:"password_hash" json:"-"`
    Password             string    `db:"-" json:"password"`
    PasswordConfirmation string    `db:"-" json:"password_confirmation"`
}

// String returns the JSON representation of the user.
func (u User) String() string {
    ju, _ := json.Marshal(u)
    return string(ju)
}

// Users is a slice of User.
type Users []User

// String returns the JSON representation of users.
func (u Users) String() string {
    ju, _ := json.Marshal(u)
    return string(ju)
}

// Validate validates common fields.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
    var err error
    return validate.Validate(
        &validators.StringIsPresent{Field: u.Email, Name: "Email"},
        &validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
        // ensure email uniqueness
        &validators.FuncValidator{
            Field:   u.Email,
            Name:    "Email",
            Message: "%s is already taken",
            Fn: func() bool {
                var exists bool
                q := tx.Where("email = ?", u.Email)
                if u.ID != uuid.Nil {
                    q = q.Where("id != ?", u.ID)
                }
                exists, err = q.Exists(u)
                if err != nil {
                    return false
                }
                return !exists
            },
        },
    ), err
}

// ValidateCreate checks password fields on create.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
    var err error
    return validate.Validate(
        &validators.StringIsPresent{Field: u.Password, Name: "Password"},
        &validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
    ), err
}

// ValidateUpdate currently does no extra validation.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
    return validate.NewErrors(), nil
}

// Create hashes the password and saves the user.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
    u.Email = strings.ToLower(strings.TrimSpace(u.Email))
    ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return validate.NewErrors(), errors.WithStack(err)
    }
    u.PasswordHash = string(ph)
    return tx.ValidateAndCreate(u)
}
