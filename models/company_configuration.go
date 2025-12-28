package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// CompanyConfiguration represents company-wide settings and configurations
type CompanyConfiguration struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Branding & Identity
	CompanyName        string `json:"company_name" db:"company_name"`
	CompanyLogoURL     string `json:"company_logo_url" db:"company_logo_url"`
	CompanyDescription string `json:"company_description" db:"company_description"`
	ContactEmail       string `json:"contact_email" db:"contact_email"`
	WebsiteURL         string `json:"website_url" db:"website_url"`
	SupportEmail       string `json:"support_email" db:"support_email"`

	// Registration & Access
	AllowPublicRegistration  bool   `json:"allow_public_registration" db:"allow_public_registration"`
	RequireEmailVerification bool   `json:"require_email_verification" db:"require_email_verification"`
	DefaultUserRole          string `json:"default_user_role" db:"default_user_role"`
	AllowGuestAccess         bool   `json:"allow_guest_access" db:"allow_guest_access"`

	// Hackathon Settings
	MaxProjectsPerUser            int  `json:"max_projects_per_user" db:"max_projects_per_user"`
	MaxTeamSize                   int  `json:"max_team_size" db:"max_team_size"`
	MaxActiveHackathons           int  `json:"max_active_hackathons" db:"max_active_hackathons"`
	DefaultHackathonDurationHours int  `json:"default_hackathon_duration_hours" db:"default_hackathon_duration_hours"`
	RequireProjectApproval        bool `json:"require_project_approval" db:"require_project_approval"`

	// Email & Notifications
	SMTPServer                string `json:"smtp_server" db:"smtp_server"`
	SMTPPort                  int    `json:"smtp_port" db:"smtp_port"`
	SMTPUsername              string `json:"smtp_username" db:"smtp_username"`
	SMTPPassword              string `json:"smtp_password" db:"smtp_password"`
	FromEmailAddress          string `json:"from_email_address" db:"from_email_address"`
	EmailNotificationsEnabled bool   `json:"email_notifications_enabled" db:"email_notifications_enabled"`

	// Security
	PasswordMinLength           int  `json:"password_min_length" db:"password_min_length"`
	PasswordRequireUppercase    bool `json:"password_require_uppercase" db:"password_require_uppercase"`
	PasswordRequireNumbers      bool `json:"password_require_numbers" db:"password_require_numbers"`
	PasswordRequireSpecialChars bool `json:"password_require_special_chars" db:"password_require_special_chars"`
	SessionTimeoutMinutes       int  `json:"session_timeout_minutes" db:"session_timeout_minutes"`
	TwoFactorRequired           bool `json:"two_factor_required" db:"two_factor_required"`

	// Legal & Compliance
	TermsOfServiceURL string `json:"terms_of_service_url" db:"terms_of_service_url"`
	PrivacyPolicyURL  string `json:"privacy_policy_url" db:"privacy_policy_url"`
	DataRetentionDays int    `json:"data_retention_days" db:"data_retention_days"`

	// Feature Toggles
	FileUploadsEnabled    bool `json:"file_uploads_enabled" db:"file_uploads_enabled"`
	ProjectImagesEnabled  bool `json:"project_images_enabled" db:"project_images_enabled"`
	TeamFormationEnabled  bool `json:"team_formation_enabled" db:"team_formation_enabled"`
	PublicProfilesEnabled bool `json:"public_profiles_enabled" db:"public_profiles_enabled"`
	AnalyticsEnabled      bool `json:"analytics_enabled" db:"analytics_enabled"`
}

// String returns the JSON representation of the company configuration
func (c CompanyConfiguration) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// CompanyConfigurations is a collection of company configurations
type CompanyConfigurations []CompanyConfiguration

// String returns the JSON representation of the company configurations
func (c CompanyConfigurations) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" method
func (c *CompanyConfiguration) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.CompanyName, Name: "CompanyName"},
		&validators.StringIsPresent{Field: c.ContactEmail, Name: "ContactEmail"},
		&validators.StringIsPresent{Field: c.DefaultUserRole, Name: "DefaultUserRole"},
	), nil
}

// GetDefaultConfig returns the default company configuration
func GetDefaultConfig(tx *pop.Connection) (*CompanyConfiguration, error) {
	config := &CompanyConfiguration{}
	err := tx.First(config)
	if err != nil {
		// If no config exists, return a default one
		return &CompanyConfiguration{
			CompanyName:                   "Hackathon Platform",
			CompanyDescription:            "A comprehensive hackathon management platform",
			ContactEmail:                  "admin@hackathon.com",
			AllowPublicRegistration:       true,
			RequireEmailVerification:      true,
			DefaultUserRole:               "hacker",
			MaxProjectsPerUser:            5,
			MaxTeamSize:                   4,
			MaxActiveHackathons:           10,
			DefaultHackathonDurationHours: 48,
			SMTPPort:                      587,
			FromEmailAddress:              "noreply@hackathon.com",
			EmailNotificationsEnabled:     true,
			PasswordMinLength:             8,
			PasswordRequireUppercase:      true,
			PasswordRequireNumbers:        true,
			SessionTimeoutMinutes:         480,
			DataRetentionDays:             2555, // ~7 years
			FileUploadsEnabled:            true,
			ProjectImagesEnabled:          true,
			TeamFormationEnabled:          true,
			PublicProfilesEnabled:         true,
		}, nil
	}
	return config, nil
}
