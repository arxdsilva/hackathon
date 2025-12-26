package actions

import (
	"net/http"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// HackathonsIndex lists all hackathons
func HackathonsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hackathons := &models.Hackathons{}

	q := tx.PaginateFromParams(c.Params())
	if err := q.Order("start_date desc").All(hackathons); err != nil {
		return err
	}

	c.Set("hackathons", hackathons)
	c.Set("pagination", q.Paginator)

	return c.Render(http.StatusOK, r.HTML("hackathons/index.plush.html"))
}

// HackathonsShow displays a single hackathon
func HackathonsShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hackathon := &models.Hackathon{}

	if err := tx.Find(hackathon, c.Param("hackathon_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("hackathon", hackathon)
	return c.Render(http.StatusOK, r.HTML("hackathons/show.plush.html"))
}

// HackathonsNew renders the form for creating a new hackathon (owner-only)
func HackathonsNew(c buffalo.Context) error {
	c.Set("hackathon", &models.Hackathon{})
	return c.Render(http.StatusOK, r.HTML("hackathons/new.plush.html"))
}

// HackathonsCreate adds a new hackathon to the DB (owner-only)
func HackathonsCreate(c buffalo.Context) error {
	hackathon := &models.Hackathon{}
	if err := c.Bind(hackathon); err != nil {
		return err
	}

	// Set the owner to current user
	currentUser := c.Value("current_user").(models.User)
	hackathon.OwnerID = currentUser.ID

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(hackathon)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("hackathon", hackathon)
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("hackathons/new.plush.html"))
	}

	c.Flash().Add("success", "Hackathon created successfully!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%d", hackathon.ID)
}

// HackathonsEdit renders the form for editing a hackathon (owner-only)
func HackathonsEdit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hackathon := &models.Hackathon{}

	if err := tx.Find(hackathon, c.Param("hackathon_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("hackathon", hackathon)
	return c.Render(http.StatusOK, r.HTML("hackathons/edit.plush.html"))
}

// HackathonsUpdate updates a hackathon in the DB (owner-only)
func HackathonsUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hackathon := &models.Hackathon{}

	if err := tx.Find(hackathon, c.Param("hackathon_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := c.Bind(hackathon); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(hackathon)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("hackathon", hackathon)
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("hackathons/edit.plush.html"))
	}

	c.Flash().Add("success", "Hackathon updated successfully!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%d", hackathon.ID)
}

// HackathonsDestroy deletes a hackathon from the DB (owner-only)
func HackathonsDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hackathon := &models.Hackathon{}

	if err := tx.Find(hackathon, c.Param("hackathon_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(hackathon); err != nil {
		return err
	}

	c.Flash().Add("success", "Hackathon deleted successfully!")
	return c.Redirect(http.StatusSeeOther, "/hackathons")
}
