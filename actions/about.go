package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// AboutHandler is a public page with information about the service.
func (a *MyApp) AboutHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("about/index.plush.html"))
}
