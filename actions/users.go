package actions

import (
	"net/http"
	"strings"

	"github.com/arxdsilva/hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// UsersNew renders the sign-up form.
func UsersNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(http.StatusOK, r.HTML("users/new.plush.html"))
}

// UsersCreate handles user registration.
func UsersCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	u := &models.User{}
	if err := c.Bind(u); err != nil {
		c.Flash().Add("danger", "Unable to read form input")
		return c.Redirect(http.StatusFound, "/users/new")
	}

	// Validate email domain is allowed
	if strings.Contains(u.Email, "@") {
		parts := strings.Split(u.Email, "@")
		if len(parts) == 2 {
			domain := strings.ToLower(strings.TrimSpace(parts[1]))
			allowed, err := models.IsDomainAllowed(tx, domain)
			if err != nil {
				c.Flash().Add("danger", "Could not validate email domain")
				return c.Redirect(http.StatusFound, "/users/new")
			}
			if !allowed {
				c.Flash().Add("danger", "Email domain is not allowed for registration")
				return c.Redirect(http.StatusFound, "/users/new")
			}
		}
	}

	verrs, err := u.Create(tx)
	if err != nil {
		c.Flash().Add("danger", "Could not create user")
		return c.Redirect(http.StatusFound, "/users/new")
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("user", u)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("users/new.plush.html"))
	}

	// Only set session if no user is currently logged in (prevents admin session switching)
	if _, ok := c.Session().Get(sessionCurrentUserID).(string); !ok {
		c.Session().Set(sessionCurrentUserID, u.ID.String())
		c.Flash().Add("success", "Account created! Welcome")
		return c.Redirect(http.StatusFound, "/")
	}

	// If admin created the user, redirect back to admin with success message
	c.Flash().Add("success", "User created successfully")
	return c.Redirect(http.StatusFound, "/admin/users")
}

// UsersEdit renders the user edit form for role management.
func UsersEdit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	userID := c.Param("user_id")
	u := &models.User{}
	if err := tx.Find(u, userID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("user", u)
	return c.Render(http.StatusOK, r.HTML("users/edit.plush.html"))
}

// UsersUpdate handles updating a user's role.
func UsersUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	userID := c.Param("user_id")
	u := &models.User{}
	if err := tx.Find(u, userID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	newRole := c.Request().FormValue("role")
	if newRole != models.RoleOwner && newRole != models.RoleHacker {
		c.Flash().Add("danger", "Invalid role")
		return c.Redirect(http.StatusFound, "/admin")
	}

	u.Role = newRole
	if err := tx.Update(u); err != nil {
		c.Flash().Add("danger", "Could not update user")
		return c.Redirect(http.StatusFound, "/admin")
	}

	c.Flash().Add("success", "User updated")
	return c.Redirect(http.StatusFound, "/admin")
}
