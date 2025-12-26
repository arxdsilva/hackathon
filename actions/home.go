package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("home/index.plush.html"))
}

// RoutesHandler renders a list of defined routes.
func RoutesHandler(c buffalo.Context) error {
	c.Set("routes", app.Routes())
	return c.Render(http.StatusOK, r.HTML("routes/index.plush.html"))
}
