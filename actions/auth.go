package actions

import (
	"net/http"
	"strings"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"golang.org/x/crypto/bcrypt"
)

const sessionCurrentUserID = "current_user_id"

// SetCurrentUser loads the current user from the session and attaches it to the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
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
func Authorize(next buffalo.Handler) buffalo.Handler {
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
func AuthNew(c buffalo.Context) error {
    c.Set("user", models.User{})
    return c.Render(http.StatusOK, r.HTML("auth/new.plush.html"))
}

// AuthCreate handles sign-in.
func AuthCreate(c buffalo.Context) error {
    tx := c.Value("tx").(*pop.Connection)
    u := models.User{}
    if err := c.Bind(&u); err != nil {
        c.Flash().Add("danger", "Invalid request")
        return c.Redirect(http.StatusFound, "/signin")
    }

    u.Email = strings.ToLower(strings.TrimSpace(u.Email))

    dbUser := models.User{}
    if err := tx.Where("email = ?", u.Email).First(&dbUser); err != nil {
        c.Flash().Add("danger", "Invalid email or password")
        return c.Redirect(http.StatusFound, "/signin")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(u.Password)); err != nil {
        c.Flash().Add("danger", "Invalid email or password")
        return c.Redirect(http.StatusFound, "/signin")
    }

    c.Session().Set(sessionCurrentUserID, dbUser.ID.String())
    c.Flash().Add("success", "Welcome back!")
    return c.Redirect(http.StatusFound, "/")
}

// AuthDestroy signs the user out.
func AuthDestroy(c buffalo.Context) error {
    c.Session().Clear()
    c.Flash().Add("success", "Signed out")
    return c.Redirect(http.StatusFound, "/")
}
