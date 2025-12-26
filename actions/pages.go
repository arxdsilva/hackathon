package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// TeamsIndex shows teams (protected).
func TeamsIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("teams/index.plush.html"))
}

// LeaderboardIndex is public.
func LeaderboardIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("leaderboard/index.plush.html"))
}

// ScheduleIndex is public.
func ScheduleIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("schedule/index.plush.html"))
}
