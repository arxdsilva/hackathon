package actions

import (
	"net/http"

	"github.com/arxdsilva/hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

// AdminIndex renders the main admin overview page
func AdminIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	// Get statistics
	var userCount, hackathonCount, projectCount, activeProjectCount int

	tx.RawQuery("SELECT COUNT(*) FROM users").First(&userCount)
	tx.RawQuery("SELECT COUNT(*) FROM hackathons").First(&hackathonCount)
	tx.RawQuery("SELECT COUNT(*) FROM projects").First(&projectCount)
	tx.RawQuery("SELECT COUNT(*) FROM projects WHERE status = 'active'").First(&activeProjectCount)

	// Get recent users
	recentUsers := &models.Users{}
	if err := tx.Order("created_at DESC").Limit(5).All(recentUsers); err != nil {
		return err
	}

	// Get recent hackathons
	recentHackathons := &models.Hackathons{}
	if err := tx.Order("created_at DESC").Limit(5).All(recentHackathons); err != nil {
		return err
	}

	// Get recent projects
	recentProjects := &models.Projects{}
	if err := tx.Order("created_at DESC").Limit(5).All(recentProjects); err != nil {
		return err
	}

	c.Set("stats", map[string]int{
		"users":          userCount,
		"hackathons":     hackathonCount,
		"projects":       projectCount,
		"activeProjects": activeProjectCount,
	})
	c.Set("recentUsers", recentUsers)
	c.Set("recentHackathons", recentHackathons)
	c.Set("recentProjects", recentProjects)

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
		q = q.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
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

	c.Set("pageTitle", "Users Management")
	return c.Render(http.StatusOK, r.HTML("admin/users/index.plush.html", "admin/layout.plush.html"))
}

// AdminUsersShow displays a specific user
func AdminUsersShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	user := &models.User{}
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Get user's projects
	projects := &models.Projects{}
	if err := tx.Where("user_id = ?", user.ID).All(projects); err != nil {
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

	if err := tx.Destroy(user); err != nil {
		c.Flash().Add("danger", "Failed to delete user")
		return c.Redirect(http.StatusFound, "/admin/users")
	}

	c.Flash().Add("success", "User deleted successfully")
	return c.Redirect(http.StatusFound, "/admin/users")
}

// AdminHackathonsIndex lists all hackathons
func AdminHackathonsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	hackathons := &models.Hackathons{}
	q := tx.PaginateFromParams(c.Params())

	// Handle search
	if search := c.Param("search"); search != "" {
		q = q.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
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
		q = q.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
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

// AdminConfigUpdate updates company configuration settings
func AdminConfigUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	config := &models.CompanyConfiguration{}
	if err := c.Bind(config); err != nil {
		return err
	}

	// Try to find existing config
	existingConfig := &models.CompanyConfiguration{}
	err := tx.First(existingConfig)
	if err != nil {
		// If no config exists, create a new one
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
	} else {
		// Update existing config
		config.ID = existingConfig.ID
		config.CreatedAt = existingConfig.CreatedAt
		verrs, err := tx.ValidateAndUpdate(config)
		if err != nil {
			return err
		}
		if verrs.HasAny() {
			c.Set("errors", verrs)
			c.Set("config", config)
			return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/config/index.plush.html", "admin/layout.plush.html"))
		}
	}

	c.Flash().Add("success", "Company configuration updated successfully!")
	return c.Redirect(http.StatusFound, "/admin/config")
}

// AdminPasswordsIndex manages password reset functionality
func AdminPasswordsIndex(c buffalo.Context) error {
	c.Set("pageTitle", "Password Reset Management")
	return c.Render(http.StatusOK, r.HTML("admin/passwords/index.plush.html", "admin/layout.plush.html"))
}
