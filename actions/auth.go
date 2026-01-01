package actions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arxdsilva/hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"golang.org/x/crypto/bcrypt"
)

const sessionCurrentUserID = "current_user_id"

// SetCurrentUser loads the current user from the session and attaches it to the context.
func (a *MyApp) SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid, ok := c.Session().Get(sessionCurrentUserID).(string); ok && uid != "" {
			tx, ok := c.Value("tx").(*pop.Connection)
			if ok {
				var u models.User
				if err := tx.Find(&u, uid); err == nil {
					c.Set("current_user", u)
				}
			}
		}
		return next(c)
	}
}

// Authorize ensures a user is signed in before proceeding.
func (a *MyApp) Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// allow public assets without auth
		if strings.HasPrefix(c.Request().URL.Path, "/assets/") {
			return next(c)
		}

		if _, ok := c.Value("current_user").(models.User); ok {
			return next(c)
		}

		c.Flash().Add("danger", "You must be signed in to access that page")
		return c.Redirect(http.StatusFound, "/signin")
	}
}

// RequirePasswordReset ensures users who need to reset their password can't access other pages.
func (a *MyApp) RequirePasswordReset(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// Check if user needs to reset password
		if currentUser, ok := c.Value("current_user").(models.User); ok {
			if currentUser.ForcePasswordReset {
				// Allow access to password reset page (handle trailing slash)
				if strings.HasPrefix(c.Request().URL.Path, "/reset-password") {
					return next(c)
				}
				// Redirect to reset password page for all other pages
				return c.Redirect(http.StatusFound, "/reset-password?required=true")
			}
		}

		// User doesn't need to reset password, allow normal access
		return next(c)
	}
}

// AuthNew renders the sign-in form.
func (a *MyApp) AuthNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(http.StatusOK, r.HTML("auth/new.plush.html"))
}

// AuthCreate handles sign-in.
func (a *MyApp) AuthCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	u := models.User{}
	if err := c.Bind(&u); err != nil {
		c.Flash().Add("danger", "Invalid request")
		return c.Redirect(http.StatusFound, "/signin")
	}

	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	repoManager := a.Repository(tx)
	dbUser, err := repoManager.UserFindByEmail(u.Email)
	if err != nil {
		c.Flash().Add("danger", "Invalid email or password")
		return c.Redirect(http.StatusFound, "/signin")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(u.Password)); err != nil {
		// Log failed login attempt
		logAuditEvent(tx, c, nil, "login_failed", "user", nil, fmt.Sprintf("Failed login attempt for email: %s", u.Email))

		// Find an admin to contact for password reset (order by creation date to get the first admin)
		adminUser := &models.User{}
		if err := tx.Where("role = ?", models.RoleOwner).Order("created_at ASC").First(adminUser); err == nil {
			c.Flash().Add("danger", fmt.Sprintf("Invalid email or password. If you've forgotten your password, please contact the administrator at %s to reset it.", adminUser.Email))
		} else {
			c.Flash().Add("danger", "Invalid email or password")
		}
		return c.Redirect(http.StatusFound, "/signin")
	}

	// Check if user needs to reset password
	if dbUser.ForcePasswordReset {
		// Clear session and set user for password reset flow
		c.Session().Clear()
		c.Session().Set(sessionCurrentUserID, dbUser.ID.String())
		return c.Redirect(http.StatusFound, "/reset-password?required=true")
	}

	// Log successful login
	logAuditEvent(tx, c, &dbUser.ID, "login", "user", &dbUser.ID, fmt.Sprintf("User logged in: %s (%s)", dbUser.Name, dbUser.Email))

	c.Session().Set(sessionCurrentUserID, dbUser.ID.String())
	c.Flash().Add("success", "Welcome back!")
	return c.Redirect(http.StatusFound, "/")
}

// AuthDestroy signs the user out.
func (a *MyApp) AuthDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	if currentUser, ok := c.Value("current_user").(models.User); ok {
		// Log logout
		logAuditEvent(tx, c, &currentUser.ID, "logout", "user", &currentUser.ID, fmt.Sprintf("User logged out: %s (%s)", currentUser.Name, currentUser.Email))
	}

	c.Session().Clear()
	c.Flash().Add("success", "Signed out")
	return c.Redirect(http.StatusFound, "/")
}

// ResetPasswordNew renders the forced password reset form.
func (a *MyApp) ResetPasswordNew(c buffalo.Context) error {
	// Check if user is logged in and needs password reset
	if currentUser, ok := c.Value("current_user").(models.User); ok {
		if !currentUser.ForcePasswordReset {
			// User doesn't need to reset password, redirect to home
			return c.Redirect(http.StatusFound, "/")
		}
		c.Set("user", currentUser)
		// Check if this is a required reset
		if c.Param("required") == "true" {
			c.Set("resetRequired", true)
		}
		return c.Render(http.StatusOK, r.HTML("auth/reset.plush.html"))
	}

	// Not logged in, redirect to signin
	c.Flash().Add("danger", "You must be signed in to reset your password")
	return c.Redirect(http.StatusFound, "/signin")
}

// ResetPasswordCreate handles forced password reset.
func (a *MyApp) ResetPasswordCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	// Get current user
	currentUser, ok := c.Value("current_user").(models.User)
	if !ok {
		c.Flash().Add("danger", "You must be signed in to reset your password")
		return c.Redirect(http.StatusFound, "/signin")
	}

	if !currentUser.ForcePasswordReset {
		// User doesn't need to reset password, redirect to home
		return c.Redirect(http.StatusFound, "/")
	}

	// Bind form data
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		c.Flash().Add("danger", "Invalid request")
		return c.Redirect(http.StatusFound, "/reset-password")
	}

	// Validate passwords
	if user.Password == "" {
		c.Flash().Add("danger", "Password is required")
		return c.Redirect(http.StatusFound, "/reset-password")
	}

	if user.Password != user.PasswordConfirmation {
		c.Flash().Add("danger", "Password confirmation does not match")
		return c.Redirect(http.StatusFound, "/reset-password")
	}

	// Validate password against company policy
	currentUser.Password = user.Password
	currentUser.PasswordConfirmation = user.PasswordConfirmation
	if verrs, err := currentUser.ValidateCreate(tx); err != nil || verrs.HasAny() {
		if err != nil {
			c.Flash().Add("danger", "Validation error")
			return c.Redirect(http.StatusFound, "/reset-password")
		}
		for _, verr := range verrs.Errors {
			for _, msg := range verr {
				c.Flash().Add("danger", msg)
			}
		}
		return c.Redirect(http.StatusFound, "/reset-password")
	}

	// Hash new password
	ph, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Flash().Add("danger", "Failed to process password")
		return c.Redirect(http.StatusFound, "/reset-password")
	}

	// Update user password and clear force reset flag
	currentUser.PasswordHash = string(ph)
	currentUser.ForcePasswordReset = false

	if err := tx.Update(&currentUser); err != nil {
		c.Flash().Add("danger", "Failed to update password")
		return c.Redirect(http.StatusFound, "/reset-password")
	}

	// Log password reset
	logAuditEvent(tx, c, &currentUser.ID, "password_reset", "user", &currentUser.ID, fmt.Sprintf("User reset password: %s (%s)", currentUser.Name, currentUser.Email))

	c.Flash().Add("success", "Password updated successfully!")
	return c.Redirect(http.StatusFound, "/")
}
