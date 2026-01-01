package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User role constants
const (
	RoleOwner  = "owner"
	RoleHacker = "hacker"
)

// User represents a registered account in the system.
type User struct {
	ID                   uuid.UUID `db:"id" json:"id"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
	Email                string    `db:"email" json:"email"`
	Name                 string    `db:"name" json:"name"`
	CompanyTeam          string    `db:"company_team" json:"company_team"`
	Role                 string    `db:"role" json:"role"`
	PasswordHash         string    `db:"password_hash" json:"-"`
	Password             string    `db:"-" json:"password"`
	PasswordConfirmation string    `db:"-" json:"password_confirmation"`
	ForcePasswordReset   bool      `db:"force_password_reset" json:"force_password_reset"`
}

// IsOwner returns true if the user is an owner.
func (u User) IsOwner() bool {
	return u.Role == RoleOwner
}

// IsHacker returns true if the user is a hacker (participant).
func (u User) IsHacker() bool {
	return u.Role == RoleHacker
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
				// Debug logging
				fmt.Printf("Email uniqueness check: Email='%s', exists=%v\n", u.Email, exists)
				return !exists
			},
		},
	), err
}

// ValidateCreate checks password fields on create.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error

	// Get company configuration for password policy
	config, configErr := GetDefaultConfig(tx)

	// Basic validators
	baseValidators := []validate.Validator{
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	}

	// Add password policy validators if config is available
	if configErr == nil {
		// Minimum length check
		minLength := config.PasswordMinLength
		if minLength < 6 {
			minLength = 6 // Enforce minimum of 6 as requested
		}

		policyValidators := []validate.Validator{
			&validators.FuncValidator{
				Name:    "Password",
				Message: "Password does not meet minimum length requirements",
				Fn: func() bool {
					return len(u.Password) >= minLength
				},
			},
		}

		// Character requirement checks
		if config.PasswordRequireUppercase {
			policyValidators = append(policyValidators, &validators.FuncValidator{
				Name:    "Password",
				Message: "Password must contain at least one uppercase letter",
				Fn: func() bool {
					for _, char := range u.Password {
						if char >= 'A' && char <= 'Z' {
							return true
						}
					}
					return false
				},
			})
		}

		if config.PasswordRequireNumbers {
			policyValidators = append(policyValidators, &validators.FuncValidator{
				Name:    "Password",
				Message: "Password must contain at least one number",
				Fn: func() bool {
					for _, char := range u.Password {
						if char >= '0' && char <= '9' {
							return true
						}
					}
					return false
				},
			})
		}

		if config.PasswordRequireSpecialChars {
			policyValidators = append(policyValidators, &validators.FuncValidator{
				Name:    "Password",
				Message: "Password must contain at least one special character",
				Fn: func() bool {
					specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
					for _, char := range u.Password {
						if strings.ContainsRune(specialChars, char) {
							return true
						}
					}
					return false
				},
			})
		}

		// Combine all validators
		allValidators := append(baseValidators, policyValidators...)
		return validate.Validate(allValidators...), err
	} else {
		// Fallback to basic minimum length of 6
		fallbackValidators := append(baseValidators, &validators.FuncValidator{
			Name:    "Password",
			Message: "Password must be at least 6 characters long",
			Fn: func() bool {
				return len(u.Password) >= 6
			},
		})
		return validate.Validate(fallbackValidators...), err
	}
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

	// Check if this is the first user - if so, make them owner
	if u.Role == "" {
		count, err := tx.Count(&User{})
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
		if count == 0 {
			u.Role = RoleOwner
		} else {
			u.Role = RoleHacker
		}
	}

	return tx.ValidateAndCreate(u)
}
