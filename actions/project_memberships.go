package actions

import (
	"fmt"
	"net/http"

	"github.com/arxdsilva/hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// ProjectMembershipsCreate allows a user to join a project
func ProjectMembershipsCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(models.User)

	projectID := c.Param("project_id")
	hackathonID := c.Param("hackathon_id")

	// Check if already a member
	count, err := tx.Where("project_id = ? AND user_id = ?", projectID, currentUser.ID).Count(&models.ProjectMembership{})
	if err != nil {
		return err
	}
	if count > 0 {
		c.Flash().Add("warning", "You are already a member of this project.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%s", hackathonID)
	}

	// Create membership
	membership := &models.ProjectMembership{
		ProjectID: 0, // will be set by binding
		UserID:    currentUser.ID,
	}

	// Parse projectID as int
	var project models.Project
	if err := tx.Find(&project, projectID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	membership.ProjectID = project.ID

	if err := tx.Create(membership); err != nil {
		return err
	}

	// Log project membership creation
	logAuditEvent(tx, c, &currentUser.ID, "join", "project_membership", &membership.ID, fmt.Sprintf("User joined project: %s", project.Name))

	c.Flash().Add("success", "You joined the project!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%s", hackathonID)
}

// ProjectMembershipsDestroy allows a user to leave a project
func ProjectMembershipsDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(models.User)

	projectID := c.Param("project_id")
	hackathonID := c.Param("hackathon_id")

	// Check if user is the project owner
	var project models.Project
	if err := tx.Find(&project, projectID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if project.UserID != nil && *project.UserID == currentUser.ID {
		c.Flash().Add("danger", "You cannot leave a project you own.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%s", hackathonID)
	}

	// Find and delete membership
	membership := &models.ProjectMembership{}
	if err := tx.Where("project_id = ? AND user_id = ?", projectID, currentUser.ID).First(membership); err != nil {
		c.Flash().Add("warning", "You are not a member of this project.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%s", hackathonID)
	}

	if err := tx.Destroy(membership); err != nil {
		return err
	}

	// Log project membership deletion
	logAuditEvent(tx, c, &currentUser.ID, "leave", "project_membership", &membership.ID, fmt.Sprintf("User left project: %s", project.Name))

	c.Flash().Add("success", "You left the project.")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%s", hackathonID)
}
