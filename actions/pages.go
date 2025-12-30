package actions

import (
	"net/http"

	"github.com/arxdsilva/hackathon/repository"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// ScheduleIndex shows schedules across all hackathons (public view).
func (a *MyApp) ScheduleIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	repoManager := repository.NewRepositoryManager(tx)

	// Get all active/upcoming hackathons with schedules
	hackathons, err := repoManager.Hackathon().GetActiveWithSchedule()
	if err != nil {
		return err
	}

	c.Set("hackathons", hackathons)
	return c.Render(http.StatusOK, r.HTML("schedule/index.plush.html"))
}
