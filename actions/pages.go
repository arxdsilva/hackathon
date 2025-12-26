package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HackathonsIndex shows the hackathons list (protected).
func HackathonsIndex(c buffalo.Context) error {
    return c.Render(http.StatusOK, r.HTML("hackathons/index.plush.html"))
}

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
