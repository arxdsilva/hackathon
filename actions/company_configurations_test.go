package actions

import (
	"net/http"
	"net/url"

	"github.com/arxdsilva/hackathon/models"
)

// NOTE: These are integration tests that require a database connection.
// For unit tests that don't require database, see company_configurations_unit_test.go

func (as *ActionSuite) Test_AdminConfigIndex_Integration() {
	// Create a test user with owner role
	user := &models.User{
		Email: "admin@test.com",
		Name:  "Admin User",
		Role:  models.RoleOwner,
	}
	as.NoError(as.DB.Create(user))

	// Set up session with the owner user
	as.Session.Set("current_user_id", user.ID.String())

	res := as.HTML("/admin/config").Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "Company Configuration")
}

func (as *ActionSuite) Test_AdminConfigUpdate_Integration() {
	// Create a test user with owner role
	user := &models.User{
		Email: "admin@test.com",
		Name:  "Admin User",
		Role:  models.RoleOwner,
	}
	as.NoError(as.DB.Create(user))

	// Set up session with the owner user
	as.Session.Set("current_user_id", user.ID.String())

	// Test updating company configuration
	form := url.Values{
		"company_name":                     {"Test Company"},
		"company_description":              {"A test company description"},
		"contact_email":                    {"contact@test.com"},
		"website_url":                      {"https://test.com"},
		"support_email":                    {"support@test.com"},
		"allow_public_registration":        {"true"},
		"require_email_verification":       {"true"},
		"default_user_role":                {"hacker"},
		"allow_guest_access":               {"false"},
		"max_projects_per_user":            {"5"},
		"max_team_size":                    {"4"},
		"max_active_hackathons":            {"10"},
		"default_hackathon_duration_hours": {"48"},
		"require_project_approval":         {"false"},
		"password_min_length":              {"8"},
		"password_require_uppercase":       {"true"},
		"password_require_numbers":         {"true"},
		"password_require_special_chars":   {"false"},
		"session_timeout_minutes":          {"480"},
		"two_factor_required":              {"false"},
		"terms_of_service_url":             {"https://test.com/tos"},
		"privacy_policy_url":               {"https://test.com/privacy"},
		"data_retention_days":              {"2555"},
		"file_uploads_enabled":             {"true"},
		"project_images_enabled":           {"true"},
		"team_formation_enabled":           {"true"},
		"public_profiles_enabled":          {"true"},
		"analytics_enabled":                {"false"},
	}

	res := as.HTML("/admin/config").Put(form)
	as.Equal(http.StatusFound, res.Code) // Should redirect after successful update

	// Verify the configuration was updated
	config := &models.CompanyConfiguration{}
	err := as.DB.First(config)
	as.NoError(err)
	as.Equal("Test Company", config.CompanyName)
	as.Equal("A test company description", config.CompanyDescription)
	as.Equal("contact@test.com", config.ContactEmail)
	as.Equal("https://test.com", config.WebsiteURL)
	as.Equal("support@test.com", config.SupportEmail)
	as.True(config.AllowPublicRegistration)
	as.True(config.RequireEmailVerification)
	as.Equal("hacker", config.DefaultUserRole)
	as.False(config.AllowGuestAccess)
	as.Equal(5, config.MaxProjectsPerUser)
	as.Equal(4, config.MaxTeamSize)
	as.Equal(10, config.MaxActiveHackathons)
	as.Equal(48, config.DefaultHackathonDurationHours)
	as.False(config.RequireProjectApproval)
	as.Equal(8, config.PasswordMinLength)
	as.True(config.PasswordRequireUppercase)
	as.True(config.PasswordRequireNumbers)
	as.False(config.PasswordRequireSpecialChars)
	as.Equal(480, config.SessionTimeoutMinutes)
	as.False(config.TwoFactorRequired)
	as.Equal("https://test.com/tos", config.TermsOfServiceURL)
	as.Equal("https://test.com/privacy", config.PrivacyPolicyURL)
	as.Equal(2555, config.DataRetentionDays)
	as.True(config.FileUploadsEnabled)
	as.True(config.ProjectImagesEnabled)
	as.True(config.TeamFormationEnabled)
	as.True(config.PublicProfilesEnabled)
	as.False(config.AnalyticsEnabled)
}

func (as *ActionSuite) Test_AdminConfigUpdate_ValidationError_Integration() {
	// Create a test user with owner role
	user := &models.User{
		Email: "admin@test.com",
		Name:  "Admin User",
		Role:  models.RoleOwner,
	}
	as.NoError(as.DB.Create(user))

	// Set up session with the owner user
	as.Session.Set("current_user_id", user.ID.String())

	// Test updating with invalid data (missing required fields)
	form := url.Values{
		"company_name":      {""}, // Required field left empty
		"contact_email":     {""}, // Required field left empty
		"default_user_role": {""}, // Required field left empty
	}

	res := as.HTML("/admin/config").Put(form)
	as.Equal(http.StatusUnprocessableEntity, res.Code) // Should return validation error
	as.Contains(res.Body.String(), "CompanyName can not be blank")
	as.Contains(res.Body.String(), "ContactEmail can not be blank")
	as.Contains(res.Body.String(), "DefaultUserRole can not be blank")
}

func (as *ActionSuite) Test_AdminConfigUpdate_NoChanges_Integration() {
	// Create a test user with owner role
	user := &models.User{
		Email: "admin@test.com",
		Name:  "Admin User",
		Role:  models.RoleOwner,
	}
	as.NoError(as.DB.Create(user))

	// Set up session with the owner user
	as.Session.Set("current_user_id", user.ID.String())

	// First create a config with default values
	defaultConfig, err := models.GetDefaultConfig(as.DB)
	as.NoError(err)
	as.NoError(as.DB.Create(defaultConfig))

	// Try to update with the same values (no changes)
	form := url.Values{
		"company_name":                     {"Hackathon Platform"},
		"company_description":              {"A comprehensive hackathon management platform"},
		"contact_email":                    {"admin@hackathon.com"},
		"allow_public_registration":        {"true"},
		"require_email_verification":       {"true"},
		"default_user_role":                {"hacker"},
		"max_projects_per_user":            {"5"},
		"max_team_size":                    {"4"},
		"max_active_hackathons":            {"10"},
		"default_hackathon_duration_hours": {"48"},
		"password_min_length":              {"8"},
		"password_require_uppercase":       {"true"},
		"password_require_numbers":         {"true"},
		"session_timeout_minutes":          {"480"},
		"data_retention_days":              {"2555"},
		"file_uploads_enabled":             {"true"},
		"project_images_enabled":           {"true"},
		"team_formation_enabled":           {"true"},
		"public_profiles_enabled":          {"true"},
	}

	res := as.HTML("/admin/config").Put(form)
	as.Equal(http.StatusFound, res.Code) // Should redirect
	as.Contains(as.Session.Get("flash").(map[string][]string)["info"][0], "No changes were made")
}

func (as *ActionSuite) Test_AdminConfigUpdate_AccessDenied_Integration() {
	// Create a test user with hacker role (not owner)
	user := &models.User{
		Email: "hacker@test.com",
		Name:  "Hacker User",
		Role:  models.RoleHacker,
	}
	as.NoError(as.DB.Create(user))

	// Set up session with the hacker user
	as.Session.Set("current_user_id", user.ID.String())

	// Try to access admin config (should be denied)
	res := as.HTML("/admin/config").Get()
	as.Equal(http.StatusFound, res.Code) // Should redirect to login or access denied
}
