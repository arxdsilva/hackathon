package actions

import (
	"net/http"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// ProfileShow displays the user's profile.
func ProfileShow(c buffalo.Context) error {
    user := c.Value("current_user").(models.User)
    c.Set("user", user)
    return c.Render(http.StatusOK, r.HTML("profile/show.plush.html"))
}

// ProfileEdit renders the profile edit form.
func ProfileEdit(c buffalo.Context) error {
    user := c.Value("current_user").(models.User)
    c.Set("user", user)
    return c.Render(http.StatusOK, r.HTML("profile/edit.plush.html"))
}

// ProfileUpdate handles profile updates.
func ProfileUpdate(c buffalo.Context) error {
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
