package actions

import (
	"net/http"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// ProjectsShow displays a single project
func ProjectsShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	hackathon := &models.Hackathon{}
	if err := tx.Find(hackathon, project.HackathonID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Check if current user is the project owner
	isOwner := false
	if cu, ok := c.Value("current_user").(models.User); ok && project.UserID != nil {
		isOwner = *project.UserID == cu.ID
	}

	c.Set("hackathon", hackathon)
	c.Set("project", project)
	c.Set("isProjectOwner", isOwner)
	return c.Render(http.StatusOK, r.HTML("projects/show.plush.html"))
}

// ProjectsNew renders the form for creating a new project
func ProjectsNew(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hackathon := &models.Hackathon{}

	if err := tx.Find(hackathon, c.Param("hackathon_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("hackathon", hackathon)
	c.Set("project", &models.Project{})
	return c.Render(http.StatusOK, r.HTML("projects/new.plush.html"))
}

// ProjectsCreate adds a new project to the DB
func ProjectsCreate(c buffalo.Context) error {
	project := &models.Project{}
	if err := c.Bind(project); err != nil {
		return err
	}

	// Set the user_id to current user
	currentUser := c.Value("current_user").(models.User)
	project.UserID = &currentUser.ID

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(project)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		hackathon := &models.Hackathon{}
		tx.Find(hackathon, project.HackathonID)
		c.Set("hackathon", hackathon)
		c.Set("project", project)
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("projects/new.plush.html"))
	}

	// Add the founder to project_memberships
	membership := &models.ProjectMembership{
		ProjectID: project.ID,
		UserID:    currentUser.ID,
	}
	if err := tx.Create(membership); err != nil {
		return err
	}

	c.Flash().Add("success", "Project created successfully!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%d/projects/%d", project.HackathonID, project.ID)
}

// ProjectsEdit renders the form for editing a project
func ProjectsEdit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Only the project owner can edit
	currentUser := c.Value("current_user").(models.User)
	if project.UserID == nil || *project.UserID != currentUser.ID {
		c.Flash().Add("danger", "You can only edit your own projects.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%d/projects/%d", project.HackathonID, project.ID)
	}

	hackathon := &models.Hackathon{}
	if err := tx.Find(hackathon, project.HackathonID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("hackathon", hackathon)
	c.Set("project", project)
	return c.Render(http.StatusOK, r.HTML("projects/edit.plush.html"))
}

// ProjectsUpdate updates a project in the DB
func ProjectsUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Only the project owner can update
	currentUser := c.Value("current_user").(models.User)
	if project.UserID == nil || *project.UserID != currentUser.ID {
		c.Flash().Add("danger", "You can only edit your own projects.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%d/projects/%d", project.HackathonID, project.ID)
	}

	if err := c.Bind(project); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(project)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		hackathon := &models.Hackathon{}
		tx.Find(hackathon, project.HackathonID)
		c.Set("hackathon", hackathon)
		c.Set("project", project)
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("projects/edit.plush.html"))
	}

	c.Flash().Add("success", "Project updated successfully!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%d/projects/%d", project.HackathonID, project.ID)
}

// ProjectsDestroy deletes a project from the DB
// ProjectsDestroy is disabled: projects are retained and cannot be deleted.
