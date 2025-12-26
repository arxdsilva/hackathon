package actions

import (
	"net/http"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// ScheduleIndex shows schedules across all hackathons (public view).
func ScheduleIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	// Get all active/upcoming hackathons with schedules
	hackathons := &models.Hackathons{}
	if err := tx.Where("status IN (?, ?) AND schedule IS NOT NULL AND schedule != ''", "upcoming", "active").Order("start_date asc").All(hackathons); err != nil {
		return err
	}

	c.Set("hackathons", hackathons)
	return c.Render(http.StatusOK, r.HTML("schedule/index.plush.html"))
}
