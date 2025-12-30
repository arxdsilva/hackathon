package actions

import (
	"fmt"
	"net/http"

	"github.com/arxdsilva/hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// ProjectMembershipsCreate allows a user to join a project
func (a *MyApp) ProjectMembershipsCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(models.User)
	repoManager := a.Repository(tx)

	projectID := c.Param("project_id")
	hackathonID := c.Param("hackathon_id")

	// Check if already a member
	isMember, err := repoManager.ProjectMembershipIsUserMember(projectID, currentUser.ID)
	if err != nil {
		return err
	}
	if isMember {
		c.Flash().Add("warning", "You are already a member of this project.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%s", hackathonID)
	}

	// Find the project
	project, err := repoManager.ProjectFindByID(projectID)
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Create membership
	membership := &models.ProjectMembership{
		ProjectID: project.ID,
		UserID:    currentUser.ID,
	}

	if err := tx.Create(membership); err != nil {
		return err
	}

	// Log project membership creation
	logAuditEvent(tx, c, &currentUser.ID, "join", "project_membership", &membership.ID, fmt.Sprintf("User joined project: %s", project.Name))

	c.Flash().Add("success", "You joined the project!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%s", hackathonID)
}

// ProjectMembershipsDestroy allows a user to leave a project
func (a *MyApp) ProjectMembershipsDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	currentUser := c.Value("current_user").(models.User)
	repoManager := a.Repository(tx)

	projectID := c.Param("project_id")
	hackathonID := c.Param("hackathon_id")

	// Check if user is the project owner
	project, err := repoManager.ProjectFindByID(projectID)
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if project.UserID != nil && *project.UserID == currentUser.ID {
		c.Flash().Add("danger", "You cannot leave a project you own.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%s", hackathonID)
	}

	// Find and delete membership
	membership, err := repoManager.ProjectMembershipFindByProjectIDAndUserID(projectID, currentUser.ID)
	if err != nil {
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
