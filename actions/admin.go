package actions

import (
	"fmt"
	"net/http"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

// logAuditEvent logs a user action to the audit_logs table
func logAuditEvent(tx *pop.Connection, c buffalo.Context, userID *uuid.UUID, action, resourceType string, resourceID interface{}, details string) {
	var resourceIDStr *string
	if resourceID != nil {
		switch v := resourceID.(type) {
		case *uuid.UUID:
			if v != nil {
				str := v.String()
				resourceIDStr = &str
			}
		case uuid.UUID:
			str := v.String()
			resourceIDStr = &str
		case *int:
			if v != nil {
				str := fmt.Sprintf("%d", *v)
				resourceIDStr = &str
			}
		case int:
			str := fmt.Sprintf("%d", v)
			resourceIDStr = &str
		case string:
			resourceIDStr = &v
		}
	}

	auditLog := &models.AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceIDStr,
		Details:      details,
		IPAddress:    c.Request().RemoteAddr,
		UserAgent:    c.Request().UserAgent(),
	}

	// Don't fail the main operation if audit logging fails
	if err := tx.Create(auditLog); err != nil {
		c.Logger().Errorf("Failed to create audit log: %v", err)
	}
}

// AdminIndex renders the main admin overview page
func AdminIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	repoManager := repository.NewRepositoryManager(tx)

	// Get statistics
	userCount, err := repoManager.User().Count()
	if err != nil {
		return err
	}

	hackathonCount, err := repoManager.Hackathon().Count()
	if err != nil {
		return err
	}

	projectCount, err := repoManager.Project().Count()
	if err != nil {
		return err
	}

	activeProjectCount, err := repoManager.Project().CountActive()
	if err != nil {
		return err
	}

	presentingProjectCount, err := repoManager.Project().CountPresenting()
	if err != nil {
		return err
	}

	// Get recent users
	recentUsers, err := repoManager.User().GetRecent(5)
	if err != nil {
		return err
	}

	// Get recent hackathons
	recentHackathons, err := repoManager.Hackathon().GetRecent(5)
	if err != nil {
		return err
	}

	// Get recent projects
	recentProjects, err := repoManager.Project().GetRecent(5)
	if err != nil {
		return err
	}

	// Get presenting projects
	presentingProjects, err := repoManager.Project().FindPresentingFromActiveHackathons()
	if err != nil {
		return err
	}

	c.Set("stats", map[string]int{
		"users":              userCount,
		"hackathons":         hackathonCount,
		"projects":           projectCount,
		"activeProjects":     activeProjectCount,
		"presentingProjects": presentingProjectCount,
	})
	c.Set("recentUsers", recentUsers)
	c.Set("recentHackathons", recentHackathons)
	c.Set("recentProjects", recentProjects)
	c.Set("presentingProjects", presentingProjects)

	c.Set("pageTitle", "Overview")
	return c.Render(http.StatusOK, r.HTML("admin/index.plush.html", "admin/layout.plush.html"))
}

// AdminUsersIndex lists all users for admin management
func AdminUsersIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	users := &models.Users{}
	q := tx.PaginateFromParams(c.Params())

	// Handle search
	if search := c.Param("search"); search != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}

	// Handle role filter
	if role := c.Param("role"); role != "" {
		q = q.Where("role = ?", role)
	}

	if err := q.All(users); err != nil {
		return err
	}

	c.Set("users", users)
	c.Set("pagination", q.Paginator)
	c.Set("search", c.Param("search"))
	c.Set("roleFilter", c.Param("role"))
	c.Set("user", models.User{}) // For the add user modal form

	c.Set("pageTitle", "Users Management")
	return c.Render(http.StatusOK, r.HTML("admin/users/index.plush.html", "admin/layout.plush.html"))
}

// AdminUsersShow displays a specific user
func AdminUsersShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	repoManager := repository.NewRepositoryManager(tx)

	user, err := repoManager.User().FindByID(c.Param("user_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Get user's projects with hackathon information
	projects, err := repoManager.Project().FindByUserIDWithHackathon(user.ID)
	if err != nil {
		return err
	}

	c.Set("user", user)
	c.Set("projects", projects)

	c.Set("pageTitle", "User Details")
	return c.Render(http.StatusOK, r.HTML("admin/users/show.plush.html", "admin/layout.plush.html"))
}

// AdminUsersEdit renders the user edit form
func AdminUsersEdit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("user", user)
	c.Set("pageTitle", "Edit User")
	return c.Render(http.StatusOK, r.HTML("admin/users/edit.plush.html", "admin/layout.plush.html"))
}

// AdminUsersUpdate handles user updates
func AdminUsersUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(user); err != nil {
		return err
	}

	// Validate role
	if user.Role != models.RoleOwner && user.Role != models.RoleHacker {
		c.Flash().Add("danger", "Invalid role specified")
		return c.Redirect(http.StatusFound, c.Request().Referer())
	}

	if err := tx.Update(user); err != nil {
		c.Flash().Add("danger", "Failed to update user")
		return c.Redirect(http.StatusFound, c.Request().Referer())
	}

	// Log the user update
	userID := c.Value("current_user").(models.User).ID
	logAuditEvent(tx, c, &userID, "update", "user", &user.ID, fmt.Sprintf("Updated user %s (%s)", user.Name, user.Email))

	c.Flash().Add("success", "User updated successfully")
	return c.Redirect(http.StatusFound, "/admin/users")
}

// AdminUsersDestroy deletes a user
func AdminUsersDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Prevent deletion of any owner accounts
	if user.IsOwner() {
		c.Flash().Add("danger", "Cannot delete owner accounts")
		return c.Redirect(http.StatusFound, "/admin/users")
	}

	if err := tx.Destroy(user); err != nil {
		c.Logger().Errorf("Failed to delete user: %v", err)
		c.Flash().Add("danger", "Failed to delete user")
		return c.Redirect(http.StatusFound, "/admin/users")
	}

	// Log the user deletion
	currentUser := c.Value("current_user").(models.User)
	logAuditEvent(tx, c, &currentUser.ID, "delete", "user", &user.ID, fmt.Sprintf("Deleted user %s (%s)", user.Name, user.Email))

	c.Logger().Infof("User deleted successfully: %s", user.Name)
	c.Flash().Add("success", "User deleted successfully")
	return c.Redirect(http.StatusFound, "/admin/users")
}

// AdminUsersNew renders the form to create a new user
func AdminUsersNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	c.Set("pageTitle", "Create New User")
	return c.Render(http.StatusOK, r.HTML("admin/users/new.plush.html", "admin/layout.plush.html"))
}

// AdminUsersCreate creates a new user from admin panel
func AdminUsersCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	u := &models.User{}
	if err := c.Bind(u); err != nil {
		c.Flash().Add("danger", "Unable to read form input")
		return c.Redirect(http.StatusFound, "/admin/users/new")
	}

	// Set default role if not provided
	if u.Role == "" {
		u.Role = models.RoleHacker
	}

	// Debug logging
	c.Logger().Infof("Attempting to create user: Email=%s, Name=%s, Role=%s", u.Email, u.Name, u.Role)

	// Skip domain validation for admin-created users
	verrs, err := tx.ValidateAndCreate(u)
	if err != nil {
		c.Flash().Add("danger", "Could not create user")
		return c.Redirect(http.StatusFound, "/admin/users/new")
	}

	if verrs.HasAny() {
		c.Logger().Infof("Validation errors: %v", verrs.Errors)
		c.Set("errors", verrs)
		c.Set("user", u)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/users/new.plush.html", "admin/layout.plush.html"))
	}

	// Log the user creation
	userID := c.Value("current_user").(models.User).ID
	logAuditEvent(tx, c, &userID, "create", "user", &u.ID, fmt.Sprintf("Created user %s (%s) with role %s", u.Name, u.Email, u.Role))

	c.Flash().Add("success", "User created successfully")
	return c.Redirect(http.StatusFound, "/admin/users")
}

// AdminHackathonsIndex lists all hackathons
func AdminHackathonsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	hackathons := &models.Hackathons{}
	q := tx.PaginateFromParams(c.Params())

	// Handle search
	if search := c.Param("search"); search != "" {
		q = q.Where("LOWER(title) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}

	// Handle status filter
	if status := c.Param("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.All(hackathons); err != nil {
		return err
	}

	c.Set("hackathons", hackathons)
	c.Set("pagination", q.Paginator)
	c.Set("search", c.Param("search"))
	c.Set("statusFilter", c.Param("status"))

	c.Set("pageTitle", "Hackathons Management")
	return c.Render(http.StatusOK, r.HTML("admin/hackathons/index.plush.html", "admin/layout.plush.html"))
}

// AdminProjectsIndex lists all projects
func AdminProjectsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	projects := &models.Projects{}
	q := tx.PaginateFromParams(c.Params())

	// Handle search
	if search := c.Param("search"); search != "" {
		q = q.Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}

	// Handle status filter
	if status := c.Param("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.All(projects); err != nil {
		return err
	}

	c.Set("projects", projects)
	c.Set("pagination", q.Paginator)
	c.Set("search", c.Param("search"))
	c.Set("statusFilter", c.Param("status"))

	c.Set("pageTitle", "Projects Management")
	return c.Render(http.StatusOK, r.HTML("admin/projects/index.plush.html", "admin/layout.plush.html"))
}

// AdminEmailsIndex manages allowed email domains
func AdminEmailsIndex(c buffalo.Context) error {
	c.Set("pageTitle", "Email Domains Management")
	return c.Render(http.StatusOK, r.HTML("admin/emails/index.plush.html", "admin/layout.plush.html"))
}

// AdminConfigIndex manages company configuration settings
func AdminConfigIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	config, err := models.GetDefaultConfig(tx)
	if err != nil {
		return err
	}

	c.Set("config", config)
	c.Set("pageTitle", "Company Configuration")
	return c.Render(http.StatusOK, r.HTML("admin/config/index.plush.html", "admin/layout.plush.html"))
}

// bindConfigBooleans manually sets boolean fields from form params since c.Bind doesn't parse "true" correctly
func bindConfigBooleans(c buffalo.Context, config *models.CompanyConfiguration) {
	config.AllowPublicRegistration = c.Param("allow_public_registration") == "true"
	config.RequireEmailVerification = c.Param("require_email_verification") == "true"
	config.AllowGuestAccess = c.Param("allow_guest_access") == "true"
	config.RequireProjectApproval = c.Param("require_project_approval") == "true"
	config.PasswordRequireUppercase = c.Param("password_require_uppercase") == "true"
	config.PasswordRequireNumbers = c.Param("password_require_numbers") == "true"
	config.PasswordRequireSpecialChars = c.Param("password_require_special_chars") == "true"
	config.TwoFactorRequired = c.Param("two_factor_required") == "true"
	config.FileUploadsEnabled = c.Param("file_uploads_enabled") == "true"
	config.ProjectImagesEnabled = c.Param("project_images_enabled") == "true"
	config.TeamFormationEnabled = c.Param("team_formation_enabled") == "true"
	config.PublicProfilesEnabled = c.Param("public_profiles_enabled") == "true"
	config.AnalyticsEnabled = c.Param("analytics_enabled") == "true"
}

// AdminConfigUpdate updates company configuration settings
func AdminConfigUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	// Find the canonical config (smallest ID)
	existingConfig := &models.CompanyConfiguration{}
	err := tx.First(existingConfig)
	if err != nil {
		// If no config exists, create a new one
		config := &models.CompanyConfiguration{}
		if err := c.Bind(config); err != nil {
			return err
		}
		// Manually handle boolean fields
		bindConfigBooleans(c, config)
		config.ID = uuid.Must(uuid.NewV4())
		verrs, err := tx.ValidateAndCreate(config)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			c.Set("errors", verrs)
			c.Set("config", config)
			return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/config/index.plush.html", "admin/layout.plush.html"))
		}
		return c.Redirect(http.StatusFound, "/admin/config")
	}

	// Update the existing config - only update fields that changed
	newConfig := &models.CompanyConfiguration{}
	if err := c.Bind(newConfig); err != nil {
		return err
	}

	// Manually handle boolean fields since c.Bind doesn't parse "true" correctly for bools
	bindConfigBooleans(c, newConfig)

	// Check if any fields changed and update them
	changed := existingConfig.UpdateChangedFields(newConfig)
	if !changed {
		c.Flash().Add("info", "No changes were made to the configuration.")
		return c.Redirect(http.StatusFound, "/admin/config")
	}

	verrs, err := tx.ValidateAndUpdate(existingConfig)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("config", existingConfig)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/config/index.plush.html", "admin/layout.plush.html"))
	}

	c.Flash().Add("success", "Company configuration updated successfully!")
	return c.Redirect(http.StatusFound, "/admin/config")
}

// AdminPasswordsIndex manages password reset functionality
func AdminPasswordsIndex(c buffalo.Context) error {
	c.Set("pageTitle", "Password Reset Management")
	return c.Render(http.StatusOK, r.HTML("admin/passwords/index.plush.html", "admin/layout.plush.html"))
}

// AdminDomainsIndex lists all company allowed domains
func AdminDomainsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	repoManager := repository.NewRepositoryManager(tx)

	domains, err := repoManager.CompanyAllowedDomain().FindAll()
	if err != nil {
		return err
	}

	c.Set("domains", domains)
	c.Set("pageTitle", "Allowed Domains Management")
	return c.Render(http.StatusOK, r.HTML("admin/domains/index.plush.html", "admin/layout.plush.html"))
}

// AdminDomainsCreate creates a new allowed domain
func AdminDomainsCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	domain := &models.CompanyAllowedDomain{}

	// Debug logging before binding
	c.Logger().Info("Before binding - domain struct:", domain)

	if err := c.Bind(domain); err != nil {
		c.Flash().Add("danger", "Unable to read form input")
		return c.Redirect(http.StatusFound, "/admin/domains")
	}

	// Debug logging after binding
	c.Logger().Info("After binding - domain:", domain.Domain)
	c.Logger().Info("After binding - description:", domain.Description)
	c.Logger().Info("After binding - is_active:", domain.IsActive)

	// Set default values
	if domain.IsActive == false {
		domain.IsActive = true // Default to active
	}

	verrs, err := tx.ValidateAndCreate(domain)
	if err != nil {
		c.Flash().Add("danger", "Could not create domain")
		return c.Redirect(http.StatusFound, "/admin/domains")
	}

	if verrs.HasAny() {
		c.Flash().Add("danger", verrs.String())
		return c.Redirect(http.StatusFound, "/admin/domains")
	}

	c.Flash().Add("success", "Domain added successfully!")
	return c.Redirect(http.StatusFound, "/admin/domains")
}

// AdminDomainsUpdate updates an existing allowed domain
func AdminDomainsUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	domainID := c.Param("domain_id")
	domain := &models.CompanyAllowedDomain{}
	if err := tx.Find(domain, domainID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind the updated values
	if err := c.Bind(domain); err != nil {
		c.Flash().Add("danger", "Unable to read form input")
		return c.Redirect(http.StatusFound, "/admin/domains")
	}

	verrs, err := tx.ValidateAndUpdate(domain)
	if err != nil {
		c.Flash().Add("danger", "Could not update domain")
		return c.Redirect(http.StatusFound, "/admin/domains")
	}

	if verrs.HasAny() {
		c.Flash().Add("danger", verrs.String())
		return c.Redirect(http.StatusFound, "/admin/domains")
	}

	c.Flash().Add("success", "Domain updated successfully!")
	return c.Redirect(http.StatusFound, "/admin/domains")
}

// AdminDomainsDestroy deletes an allowed domain
func AdminDomainsDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	domainID := c.Param("domain_id")
	domain := &models.CompanyAllowedDomain{}
	if err := tx.Find(domain, domainID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(domain); err != nil {
		c.Flash().Add("danger", "Could not delete domain")
		return c.Redirect(http.StatusFound, "/admin/domains")
	}

	c.Flash().Add("success", "Domain deleted successfully!")
	return c.Redirect(http.StatusFound, "/admin/domains")
}

// AdminAuditLogsIndex displays audit logs
func AdminAuditLogsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	auditLogs := &models.AuditLogs{}
	q := tx.PaginateFromParams(c.Params())

	// Handle search
	if search := c.Param("search"); search != "" {
		q = q.Where("LOWER(action) LIKE LOWER(?) OR LOWER(resource_type) LIKE LOWER(?) OR LOWER(details) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Handle filters
	if action := c.Param("action"); action != "" {
		q = q.Where("action = ?", action)
	}

	if resourceType := c.Param("resource_type"); resourceType != "" {
		q = q.Where("resource_type = ?", resourceType)
	}

	if err := q.Order("created_at DESC").All(auditLogs); err != nil {
		return err
	}

	// Fetch user information for logs that have a user_id
	userMap := make(map[string]string)
	if len(*auditLogs) > 0 {
		userIDs := make([]string, 0)
		for _, log := range *auditLogs {
			if log.UserID != nil {
				userIDs = append(userIDs, log.UserID.String())
			}
		}
		if len(userIDs) > 0 {
			users := &models.Users{}
			if err := tx.Where("id IN (?)", userIDs).All(users); err == nil {
				for _, user := range *users {
					userMap[user.ID.String()] = user.Name
				}
			}
		}
	}

	c.Set("auditLogs", auditLogs)
	c.Set("pagination", q.Paginator)
	c.Set("search", c.Param("search"))
	c.Set("actionFilter", c.Param("action"))
	c.Set("resourceTypeFilter", c.Param("resource_type"))
	c.Set("userMap", userMap)
	c.Set("pageTitle", "Audit Logs")
	return c.Render(http.StatusOK, r.HTML("admin/audit_logs/index.plush.html", "admin/layout.plush.html"))
}
