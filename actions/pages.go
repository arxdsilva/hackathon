package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// ScheduleIndex is public.
func ScheduleIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("schedule/index.plush.html"))
}
