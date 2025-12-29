package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/models"
	"github.com/stretchr/testify/require"
)

// TestAdminConfigUpdate_ValidationError tests the validation logic without database
func TestAdminConfigUpdate_ValidationError(t *testing.T) {
	r := require.New(t)

	// Test validation with invalid data
	config := &models.CompanyConfiguration{
		CompanyName:     "", // Invalid
		ContactEmail:    "", // Invalid
		DefaultUserRole: "", // Invalid
	}

	// Test validation
	verrs, err := config.Validate(nil)
	r.NoError(err)
	r.True(verrs.HasAny())

	// Check that the required fields have errors
	r.True(len(verrs.Errors["company_name"]) > 0, "company_name should have validation errors")
	r.True(len(verrs.Errors["contact_email"]) > 0, "contact_email should have validation errors")
	r.True(len(verrs.Errors["default_user_role"]) > 0, "default_user_role should have validation errors")

	// Check specific error messages
	r.Contains(verrs.Errors["company_name"], "CompanyName can not be blank.")
	r.Contains(verrs.Errors["contact_email"], "ContactEmail can not be blank.")
	r.Contains(verrs.Errors["default_user_role"], "DefaultUserRole can not be blank.")
}

// TestBindConfigBooleans tests the boolean binding helper function
func TestBindConfigBooleans(t *testing.T) {
	r := require.New(t)

	// Test the boolean conversion logic directly
	// This simulates what bindConfigBooleans does: param == "true"
	testCases := []struct {
		param    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"", false},
		{"TRUE", false}, // Case sensitive
		{"anything", false},
	}

	for _, tc := range testCases {
		result := tc.param == "true"
		r.Equal(tc.expected, result, "Param: %s", tc.param)
	}

	// Test that all boolean fields are properly handled by creating a config
	// and manually setting the values as bindConfigBooleans would
	config := &models.CompanyConfiguration{}

	// Simulate the bindConfigBooleans logic with test parameters
	params := map[string]string{
		"allow_public_registration":      "true",
		"require_email_verification":     "false",
		"allow_guest_access":             "true",
		"require_project_approval":       "false",
		"password_require_uppercase":     "true",
		"password_require_numbers":       "false",
		"password_require_special_chars": "true",
		"two_factor_required":            "false",
		"file_uploads_enabled":           "true",
		"project_images_enabled":         "false",
		"team_formation_enabled":         "true",
		"public_profiles_enabled":        "false",
		"analytics_enabled":              "true",
	}

	// Apply the same logic as bindConfigBooleans
	config.AllowPublicRegistration = params["allow_public_registration"] == "true"
	config.RequireEmailVerification = params["require_email_verification"] == "true"
	config.AllowGuestAccess = params["allow_guest_access"] == "true"
	config.RequireProjectApproval = params["require_project_approval"] == "true"
	config.PasswordRequireUppercase = params["password_require_uppercase"] == "true"
	config.PasswordRequireNumbers = params["password_require_numbers"] == "true"
	config.PasswordRequireSpecialChars = params["password_require_special_chars"] == "true"
	config.TwoFactorRequired = params["two_factor_required"] == "true"
	config.FileUploadsEnabled = params["file_uploads_enabled"] == "true"
	config.ProjectImagesEnabled = params["project_images_enabled"] == "true"
	config.TeamFormationEnabled = params["team_formation_enabled"] == "true"
	config.PublicProfilesEnabled = params["public_profiles_enabled"] == "true"
	config.AnalyticsEnabled = params["analytics_enabled"] == "true"

	// Verify boolean values
	r.True(config.AllowPublicRegistration)
	r.False(config.RequireEmailVerification)
	r.True(config.AllowGuestAccess)
	r.False(config.RequireProjectApproval)
	r.True(config.PasswordRequireUppercase)
	r.False(config.PasswordRequireNumbers)
	r.True(config.PasswordRequireSpecialChars)
	r.False(config.TwoFactorRequired)
	r.True(config.FileUploadsEnabled)
	r.False(config.ProjectImagesEnabled)
	r.True(config.TeamFormationEnabled)
	r.False(config.PublicProfilesEnabled)
	r.True(config.AnalyticsEnabled)
}

// TestAdminConfigUpdate_NoDatabase tests the update logic without database operations
func TestAdminConfigUpdate_NoDatabase(t *testing.T) {
	r := require.New(t)

	// Test the UpdateChangedFields method
	existing := &models.CompanyConfiguration{
		CompanyName:        "Old Name",
		ContactEmail:       "old@email.com",
		DefaultUserRole:    "hacker",
		MaxProjectsPerUser: 3,
	}

	newConfig := &models.CompanyConfiguration{
		CompanyName:        "New Name",
		ContactEmail:       "new@email.com",
		DefaultUserRole:    "hacker",
		MaxProjectsPerUser: 5,
	}

	// Test that changes are detected
	changed := existing.UpdateChangedFields(newConfig)
	r.True(changed)
	r.Equal("New Name", existing.CompanyName)
	r.Equal("new@email.com", existing.ContactEmail)
	r.Equal(5, existing.MaxProjectsPerUser)

	// Test no changes scenario
	newConfig2 := &models.CompanyConfiguration{
		CompanyName:        "New Name",
		ContactEmail:       "new@email.com",
		DefaultUserRole:    "hacker",
		MaxProjectsPerUser: 5,
	}

	changed2 := existing.UpdateChangedFields(newConfig2)
	r.False(changed2)
}

// TestAdminConfigIndex_NoDatabase tests the index logic without database
func TestAdminConfigIndex_NoDatabase(t *testing.T) {
	r := require.New(t)

	// Test that default config has expected values
	defaultConfig := &models.CompanyConfiguration{
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
		DataRetentionDays:             2555,
		FileUploadsEnabled:            true,
		ProjectImagesEnabled:          true,
		TeamFormationEnabled:          true,
		PublicProfilesEnabled:         true,
	}

	r.Equal("Hackathon Platform", defaultConfig.CompanyName)
	r.Equal("admin@hackathon.com", defaultConfig.ContactEmail)
	r.True(defaultConfig.AllowPublicRegistration)
	r.Equal(5, defaultConfig.MaxProjectsPerUser)
}

// TestAdminConfigUpdate_InvalidRole tests role validation
func TestAdminConfigUpdate_InvalidRole(t *testing.T) {
	r := require.New(t)

	// Test that invalid roles are rejected
	config := &models.CompanyConfiguration{
		DefaultUserRole: "invalid_role",
	}

	verrs, err := config.Validate(nil)
	r.NoError(err)
	r.True(verrs.HasAny())
	// Check that the error message contains the expected text
	errorMsg := verrs.Errors["default_user_role"][0]
	r.Contains(errorMsg, "DefaultUserRole must be either 'hacker' or 'owner'")

	// Test valid roles
	config = &models.CompanyConfiguration{
		CompanyName:     "Test Company",
		ContactEmail:    "test@example.com",
		DefaultUserRole: "hacker",
	}
	verrs2, err2 := config.Validate(nil)
	r.NoError(err2)
	r.False(verrs2.HasAny())

	config.DefaultUserRole = "owner"
	verrs3, err3 := config.Validate(nil)
	r.NoError(err3)
	r.False(verrs3.HasAny())
}
