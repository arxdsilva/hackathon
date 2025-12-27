package actions

import (
	"net/http"
	"strconv"
	"time"

	"github.com/arxdsilva/hackathon/models"

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

	// Paginated projects for this hackathon (default 20 per page)
	page := 1
	if p := c.Param("page"); p != "" {
		if i, err := strconv.Atoi(p); err == nil && i > 0 {
			page = i
		}
	}

	projects := &models.Projects{}
	q := tx.Where("hackathon_id = ?", hackathon.ID).Order("created_at desc").Paginate(page, 20)
	if err := q.Eager("User").All(projects); err != nil {
		return err
	}

	// Count memberships for each project and check if current user is a member
	memberCounts := make(map[int]int)
	userMemberships := make(map[int]bool)
	currentUser := c.Value("current_user").(models.User)

	for _, project := range *projects {
		count, err := tx.Where("project_id = ?", project.ID).Count(&models.ProjectMembership{})
		if err == nil {
			memberCounts[project.ID] = count
		}

		// Check if current user is a member of this project
		memberCount, err := tx.Where("project_id = ? AND user_id = ?", project.ID, currentUser.ID).Count(&models.ProjectMembership{})
		if err == nil {
			userMemberships[project.ID] = memberCount > 0
		}
	}

	// Determine if current user can create a project (max 1 per hackathon)
	canCreate := true
	if cu, ok := c.Value("current_user").(models.User); ok {
		count, err := tx.Where("hackathon_id = ? AND user_id = ?", hackathon.ID, cu.ID).
			Count(&models.Project{})
		if err == nil {
			canCreate = count == 0
		}
	}

	// Calculate statistics for the template
	totalParticipants := 0
	totalTeams := 0
	durationDays := 1

	// Count total participants across all projects
	for _, project := range *projects {
		if count, exists := memberCounts[project.ID]; exists {
			totalParticipants += count
		}
		if project.UserID != nil {
			totalTeams++
		}
	}

	// Calculate duration in days
	if hackathon.EndDate.After(hackathon.StartDate) {
		duration := hackathon.EndDate.Sub(hackathon.StartDate)
		durationDays = int(duration.Hours() / 24)
		if durationDays < 1 {
			durationDays = 1
		}
	}

	// Calculate progress percentage
	totalProjects := len(*projects)
	activeProjects := 0
	for _, project := range *projects {
		if project.Status == "active" {
			activeProjects++
		}
	}
	progressPercentage := 0
	if totalProjects > 0 {
		progressPercentage = (activeProjects * 100) / totalProjects
	}

	c.Set("hackathon", hackathon)
	c.Set("projects", projects)
	c.Set("pagination", q.Paginator)
	c.Set("memberCounts", memberCounts)
	c.Set("userMemberships", userMemberships)
	c.Set("canCreateProject", canCreate)
	c.Set("totalParticipants", totalParticipants)
	c.Set("totalTeams", totalTeams)
	c.Set("durationDays", durationDays)
	c.Set("progressPercentage", progressPercentage)
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

	// Manually parse form fields to avoid binding ID
	hackathon.Title = c.Params().Get("Title")
	hackathon.Description = c.Params().Get("Description")
	hackathon.Status = c.Params().Get("Status")
	hackathon.Schedule = c.Params().Get("Schedule")

	// Parse dates
	if startStr := c.Params().Get("StartDate"); startStr != "" {
		if start, err := time.Parse("2006-01-02T15:04", startStr); err == nil {
			hackathon.StartDate = start
		}
	}
	if endStr := c.Params().Get("EndDate"); endStr != "" {
		if end, err := time.Parse("2006-01-02T15:04", endStr); err == nil {
			hackathon.EndDate = end
		}
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

	// Manually parse form fields
	hackathon.Title = c.Params().Get("Title")
	hackathon.Description = c.Params().Get("Description")
	hackathon.Status = c.Params().Get("Status")
	hackathon.Schedule = c.Params().Get("Schedule")

	// Parse dates
	if startStr := c.Params().Get("StartDate"); startStr != "" {
		if start, err := time.Parse("2006-01-02T15:04", startStr); err == nil {
			hackathon.StartDate = start
		}
	}
	if endStr := c.Params().Get("EndDate"); endStr != "" {
		if end, err := time.Parse("2006-01-02T15:04", endStr); err == nil {
			hackathon.EndDate = end
		}
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

// DashboardShow displays the owner dashboard with users and hackathon data
func DashboardShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	// Fetch all users
	users := &models.Users{}
	if err := tx.All(users); err != nil {
		return err
	}

	// Fetch all hackathons
	hackathons := &models.Hackathons{}
	if err := tx.All(hackathons); err != nil {
		return err
	}

	// Fetch all projects
	projects := &models.Projects{}
	if err := tx.All(projects); err != nil {
		return err
	}

	c.Set("users", users)
	c.Set("hackathons", hackathons)
	c.Set("projects", projects)

	return c.Render(http.StatusOK, r.HTML("dashboard/show.plush.html"))
}
