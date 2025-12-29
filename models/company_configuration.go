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
	CompanyName        string `json:"company_name" db:"company_name" form:"company_name"`
	CompanyLogoURL     string `json:"company_logo_url" db:"company_logo_url" form:"company_logo_url"`
	CompanyDescription string `json:"company_description" db:"company_description" form:"company_description"`
	ContactEmail       string `json:"contact_email" db:"contact_email" form:"contact_email"`
	WebsiteURL         string `json:"website_url" db:"website_url" form:"website_url"`
	SupportEmail       string `json:"support_email" db:"support_email" form:"support_email"`

	// Registration & Access
	AllowPublicRegistration  bool   `json:"allow_public_registration" db:"allow_public_registration" form:"allow_public_registration"`
	RequireEmailVerification bool   `json:"require_email_verification" db:"require_email_verification" form:"require_email_verification"`
	DefaultUserRole          string `json:"default_user_role" db:"default_user_role" form:"default_user_role"`
	AllowGuestAccess         bool   `json:"allow_guest_access" db:"allow_guest_access" form:"allow_guest_access"`

	// Hackathon Settings
	MaxProjectsPerUser            int  `json:"max_projects_per_user" db:"max_projects_per_user" form:"max_projects_per_user"`
	MaxTeamSize                   int  `json:"max_team_size" db:"max_team_size" form:"max_team_size"`
	MaxActiveHackathons           int  `json:"max_active_hackathons" db:"max_active_hackathons" form:"max_active_hackathons"`
	DefaultHackathonDurationHours int  `json:"default_hackathon_duration_hours" db:"default_hackathon_duration_hours" form:"default_hackathon_duration_hours"`
	RequireProjectApproval        bool `json:"require_project_approval" db:"require_project_approval" form:"require_project_approval"`

	// Security
	PasswordMinLength           int  `json:"password_min_length" db:"password_min_length" form:"password_min_length"`
	PasswordRequireUppercase    bool `json:"password_require_uppercase" db:"password_require_uppercase" form:"password_require_uppercase"`
	PasswordRequireNumbers      bool `json:"password_require_numbers" db:"password_require_numbers" form:"password_require_numbers"`
	PasswordRequireSpecialChars bool `json:"password_require_special_chars" db:"password_require_special_chars" form:"password_require_special_chars"`
	SessionTimeoutMinutes       int  `json:"session_timeout_minutes" db:"session_timeout_minutes" form:"session_timeout_minutes"`
	TwoFactorRequired           bool `json:"two_factor_required" db:"two_factor_required" form:"two_factor_required"`

	// Legal & Compliance
	TermsOfServiceURL string `json:"terms_of_service_url" db:"terms_of_service_url" form:"terms_of_service_url"`
	PrivacyPolicyURL  string `json:"privacy_policy_url" db:"privacy_policy_url" form:"privacy_policy_url"`
	DataRetentionDays int    `json:"data_retention_days" db:"data_retention_days" form:"data_retention_days"`

	// Feature Toggles
	FileUploadsEnabled    bool `json:"file_uploads_enabled" db:"file_uploads_enabled" form:"file_uploads_enabled"`
	ProjectImagesEnabled  bool `json:"project_images_enabled" db:"project_images_enabled" form:"project_images_enabled"`
	TeamFormationEnabled  bool `json:"team_formation_enabled" db:"team_formation_enabled" form:"team_formation_enabled"`
	PublicProfilesEnabled bool `json:"public_profiles_enabled" db:"public_profiles_enabled" form:"public_profiles_enabled"`
	AnalyticsEnabled      bool `json:"analytics_enabled" db:"analytics_enabled" form:"analytics_enabled"`
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

// UpdateChangedFields updates the existing configuration with changed fields from newConfig
// Returns true if any fields were changed
func (oldConfig *CompanyConfiguration) UpdateChangedFields(newConfig *CompanyConfiguration) bool {
	if oldConfig == nil || newConfig == nil {
		return false
	}

	updated := *oldConfig
	changed := false

	if newConfig.CompanyName != oldConfig.CompanyName {
		updated.CompanyName = newConfig.CompanyName
		changed = true
	}
	if newConfig.CompanyLogoURL != oldConfig.CompanyLogoURL {
		updated.CompanyLogoURL = newConfig.CompanyLogoURL
		changed = true
	}
	if newConfig.CompanyDescription != oldConfig.CompanyDescription {
		updated.CompanyDescription = newConfig.CompanyDescription
		changed = true
	}
	if newConfig.ContactEmail != oldConfig.ContactEmail {
		updated.ContactEmail = newConfig.ContactEmail
		changed = true
	}
	if newConfig.WebsiteURL != oldConfig.WebsiteURL {
		updated.WebsiteURL = newConfig.WebsiteURL
		changed = true
	}
	if newConfig.SupportEmail != oldConfig.SupportEmail {
		updated.SupportEmail = newConfig.SupportEmail
		changed = true
	}
	if newConfig.AllowPublicRegistration != oldConfig.AllowPublicRegistration {
		updated.AllowPublicRegistration = newConfig.AllowPublicRegistration
		changed = true
	}
	if newConfig.RequireEmailVerification != oldConfig.RequireEmailVerification {
		updated.RequireEmailVerification = newConfig.RequireEmailVerification
		changed = true
	}
	if newConfig.DefaultUserRole != oldConfig.DefaultUserRole {
		updated.DefaultUserRole = newConfig.DefaultUserRole
		changed = true
	}
	if newConfig.AllowGuestAccess != oldConfig.AllowGuestAccess {
		updated.AllowGuestAccess = newConfig.AllowGuestAccess
		changed = true
	}
	if newConfig.MaxProjectsPerUser != oldConfig.MaxProjectsPerUser {
		updated.MaxProjectsPerUser = newConfig.MaxProjectsPerUser
		changed = true
	}
	if newConfig.MaxTeamSize != oldConfig.MaxTeamSize {
		updated.MaxTeamSize = newConfig.MaxTeamSize
		changed = true
	}
	if newConfig.MaxActiveHackathons != oldConfig.MaxActiveHackathons {
		updated.MaxActiveHackathons = newConfig.MaxActiveHackathons
		changed = true
	}
	if newConfig.DefaultHackathonDurationHours != oldConfig.DefaultHackathonDurationHours {
		updated.DefaultHackathonDurationHours = newConfig.DefaultHackathonDurationHours
		changed = true
	}
	if newConfig.RequireProjectApproval != oldConfig.RequireProjectApproval {
		updated.RequireProjectApproval = newConfig.RequireProjectApproval
		changed = true
	}
	if newConfig.PasswordMinLength != oldConfig.PasswordMinLength {
		updated.PasswordMinLength = newConfig.PasswordMinLength
		changed = true
	}
	if newConfig.PasswordRequireUppercase != oldConfig.PasswordRequireUppercase {
		updated.PasswordRequireUppercase = newConfig.PasswordRequireUppercase
		changed = true
	}
	if newConfig.PasswordRequireNumbers != oldConfig.PasswordRequireNumbers {
		updated.PasswordRequireNumbers = newConfig.PasswordRequireNumbers
		changed = true
	}
	if newConfig.PasswordRequireSpecialChars != oldConfig.PasswordRequireSpecialChars {
		updated.PasswordRequireSpecialChars = newConfig.PasswordRequireSpecialChars
		changed = true
	}
	if newConfig.SessionTimeoutMinutes != oldConfig.SessionTimeoutMinutes {
		updated.SessionTimeoutMinutes = newConfig.SessionTimeoutMinutes
		changed = true
	}
	if newConfig.TwoFactorRequired != oldConfig.TwoFactorRequired {
		updated.TwoFactorRequired = newConfig.TwoFactorRequired
		changed = true
	}
	if newConfig.TermsOfServiceURL != oldConfig.TermsOfServiceURL {
		updated.TermsOfServiceURL = newConfig.TermsOfServiceURL
		changed = true
	}
	if newConfig.PrivacyPolicyURL != oldConfig.PrivacyPolicyURL {
		updated.PrivacyPolicyURL = newConfig.PrivacyPolicyURL
		changed = true
	}
	if newConfig.DataRetentionDays != oldConfig.DataRetentionDays {
		updated.DataRetentionDays = newConfig.DataRetentionDays
		changed = true
	}
	if newConfig.FileUploadsEnabled != oldConfig.FileUploadsEnabled {
		updated.FileUploadsEnabled = newConfig.FileUploadsEnabled
		changed = true
	}
	if newConfig.ProjectImagesEnabled != oldConfig.ProjectImagesEnabled {
		updated.ProjectImagesEnabled = newConfig.ProjectImagesEnabled
		changed = true
	}
	if newConfig.TeamFormationEnabled != oldConfig.TeamFormationEnabled {
		updated.TeamFormationEnabled = newConfig.TeamFormationEnabled
		changed = true
	}
	if newConfig.PublicProfilesEnabled != oldConfig.PublicProfilesEnabled {
		updated.PublicProfilesEnabled = newConfig.PublicProfilesEnabled
		changed = true
	}
	if newConfig.AnalyticsEnabled != oldConfig.AnalyticsEnabled {
		updated.AnalyticsEnabled = newConfig.AnalyticsEnabled
		changed = true
	}

	if changed {
		*oldConfig = updated
	}
	return changed
}

// GetDefaultConfig returns the default company configuration, loading from DB if exists
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
