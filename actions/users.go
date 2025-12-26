package actions

import (
	"net/http"

	"hackathon/models"

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

	c.Session().Set(sessionCurrentUserID, u.ID.String())
	c.Flash().Add("success", "Account created! Welcome")
	return c.Redirect(http.StatusFound, "/")
}
