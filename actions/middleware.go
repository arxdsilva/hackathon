package actions

import (
	"net/http"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// RequireLogin middleware ensures the user is logged in.
func RequireLogin(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		_, ok := c.Value("current_user").(models.User)
		if !ok {
			c.Flash().Add("danger", "You must be logged in to access that page")
			return c.Redirect(http.StatusFound, "/signin")
		}
		return next(c)
	}
}

// RequireRoleOwner middleware ensures the user has the "owner" role.
func RequireRoleOwner(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user, ok := c.Value("current_user").(models.User)
		if !ok || user.Role != models.RoleOwner {
			c.Flash().Add("danger", "You must have owner role to access that page")
			return c.Redirect(http.StatusFound, "/")
		}
		return next(c)
	}
}

// RequireHackathonOwner middleware ensures only the owner of the hackathon can access a route.
func RequireHackathonOwner(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user, ok := c.Value("current_user").(models.User)
		if !ok {
			c.Flash().Add("danger", "You must be logged in")
			return c.Redirect(http.StatusFound, "/")
		}

		tx := c.Value("tx").(*pop.Connection)
		hackathon := &models.Hackathon{}
		if err := tx.Find(hackathon, c.Param("hackathon_id")); err != nil {
			return c.Error(http.StatusNotFound, err)
		}

		if hackathon.OwnerID != user.ID {
			c.Flash().Add("danger", "You must be the owner of this hackathon to access that page")
			return c.Redirect(http.StatusFound, "/")
		}

		return next(c)
	}
}
