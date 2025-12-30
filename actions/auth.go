package actions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"

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

	repoManager := repository.NewRepositoryManager(tx)
	dbUser, err := repoManager.User().FindByEmail(u.Email)
	if err != nil {
		c.Flash().Add("danger", "Invalid email or password")
		return c.Redirect(http.StatusFound, "/signin")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(u.Password)); err != nil {
		// Log failed login attempt
		logAuditEvent(tx, c, nil, "login_failed", "user", nil, fmt.Sprintf("Failed login attempt for email: %s", u.Email))
		c.Flash().Add("danger", "Invalid email or password")
		return c.Redirect(http.StatusFound, "/signin")
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
