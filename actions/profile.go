package actions

import (
	"net/http"

	"github.com/arxdsilva/hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"golang.org/x/crypto/bcrypt"
)

// ProfileShow displays the user's profile.
func (a *MyApp) ProfileShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(models.User)
	repoManager := a.Repository(tx)

	// Fetch hackathons owned by this user
	ownedHackathons, err := repoManager.HackathonFindByOwnerID(user.ID)
	if err != nil {
		return err
	}

	// Fetch projects created by this user
	createdProjects, err := repoManager.ProjectFindByUserID(user.ID)
	if err != nil {
		return err
	}

	// For member projects, we still need the complex logic to get projects through memberships
	// This is a bit complex for the repository pattern, so we'll keep it for now
	var memberships models.ProjectMemberships
	if err := tx.Where("user_id = ?", user.ID).Eager("Project").All(&memberships); err != nil {
		return err
	}

	// Extract member projects
	var memberProjects models.Projects
	for _, membership := range memberships {
		if membership.Project != nil {
			memberProjects = append(memberProjects, *membership.Project)
		}
	}

	// Combine and deduplicate projects
	allProjects := append(*createdProjects, memberProjects...)

	c.Set("user", user)
	c.Set("ownedHackathons", ownedHackathons)
	c.Set("projects", allProjects)
	return c.Render(http.StatusOK, r.HTML("profile/show.plush.html"))
}

// ProfileEdit renders the profile edit form.
func (a *MyApp) ProfileEdit(c buffalo.Context) error {
	user := c.Value("current_user").(models.User)
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("profile/edit.plush.html"))
}

// ProfileUpdate handles profile updates.
func (a *MyApp) ProfileUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(models.User)

	if err := c.Bind(&user); err != nil {
		c.Flash().Add("danger", "Unable to read form input")
		return c.Redirect(http.StatusFound, "/profile")
	}

	verrs, err := tx.ValidateAndUpdate(&user)
	if err != nil {
		c.Flash().Add("danger", "Error updating profile")
		return c.Redirect(http.StatusFound, "/profile")
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("profile/edit.plush.html"))
	}

	// Update the session with the new user data
	c.Session().Set(sessionCurrentUserID, user.ID.String())
	c.Flash().Add("success", "Profile updated!")
	return c.Redirect(http.StatusFound, "/profile")
}

// ProfileChangePassword handles password change requests.
func (a *MyApp) ProfileChangePassword(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(models.User)

	// Get form data
	currentPassword := c.Param("CurrentPassword")
	newPassword := c.Param("Password")
	newPasswordConfirmation := c.Param("PasswordConfirmation")

	// Validate current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		c.Flash().Add("danger", "Current password is incorrect")
		c.Set("user", user)
		c.Set("passwordErrors", map[string][]string{
			"CurrentPassword": {"Current password is incorrect"},
		})
		return c.Render(http.StatusUnprocessableEntity, r.HTML("profile/edit.plush.html"))
	}

	// Validate new password is not empty
	if newPassword == "" {
		c.Flash().Add("danger", "New password is required")
		c.Set("user", user)
		c.Set("passwordErrors", map[string][]string{
			"Password": {"New password is required"},
		})
		return c.Render(http.StatusUnprocessableEntity, r.HTML("profile/edit.plush.html"))
	}

	// Validate confirmation matches
	if newPassword != newPasswordConfirmation {
		c.Flash().Add("danger", "Password confirmation does not match")
		c.Set("user", user)
		c.Set("passwordErrors", map[string][]string{
			"PasswordConfirmation": {"Passwords do not match"},
		})
		return c.Render(http.StatusUnprocessableEntity, r.HTML("profile/edit.plush.html"))
	}

	// Validate new password is different from current
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPassword)); err == nil {
		c.Flash().Add("danger", "New password must be different from your current password")
		c.Set("user", user)
		c.Set("passwordErrors", map[string][]string{
			"Password": {"New password must be different from your current password"},
		})
		return c.Render(http.StatusUnprocessableEntity, r.HTML("profile/edit.plush.html"))
	}

	// Validate against company password policy
	tempUser := &models.User{
		Email:                user.Email,
		Password:             newPassword,
		PasswordConfirmation: newPasswordConfirmation,
	}
	if verrs, err := tempUser.ValidateCreate(tx); err != nil || verrs.HasAny() {
		if verrs.HasAny() {
			c.Set("user", user)
			c.Set("passwordErrors", verrs)
			return c.Render(http.StatusUnprocessableEntity, r.HTML("profile/edit.plush.html"))
		}
		c.Flash().Add("danger", "Validation error")
		return c.Redirect(http.StatusFound, "/profile/edit")
	}

	// Hash new password
	ph, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Flash().Add("danger", "Failed to process password")
		return c.Redirect(http.StatusFound, "/profile/edit")
	}

	// Update user password
	user.PasswordHash = string(ph)
	if err := tx.Update(&user); err != nil {
		c.Flash().Add("danger", "Failed to update password")
		return c.Redirect(http.StatusFound, "/profile/edit")
	}

	c.Flash().Add("success", "Password updated successfully!")
	return c.Redirect(http.StatusFound, "/profile")
}
