package actions

import (
	"net/http"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
)

// RequireOwner middleware ensures only owners can access a route.
func RequireOwner(next buffalo.Handler) buffalo.Handler {
    return func(c buffalo.Context) error {
        user, ok := c.Value("current_user").(models.User)
        if !ok || !user.IsOwner() {
            c.Flash().Add("danger", "You must be an owner to access that page")
            return c.Redirect(http.StatusFound, "/")
        }
        return next(c)
    }
}
