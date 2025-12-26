package actions

import (
	"net/http"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// ProjectsIndex lists all projects for a hackathon
func ProjectsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hackathon := &models.Hackathon{}

	if err := tx.Find(hackathon, c.Param("hackathon_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	projects := &models.Projects{}
	if err := tx.Where("hackathon_id = ?", hackathon.ID).All(projects); err != nil {
		return err
	}

	c.Set("hackathon", hackathon)
	c.Set("projects", projects)

	// Compute create permission: user can create only if they have no project in this hackathon
	canCreate := true
	if cu, ok := c.Value("current_user").(models.User); ok {
		count, err := tx.Where("hackathon_id = ? AND user_id = ?", hackathon.ID, cu.ID).Count(&models.Project{})
		if err == nil {
			canCreate = count == 0
		}
	}
	c.Set("canCreateProject", canCreate)

	return c.Render(http.StatusOK, r.HTML("projects/index.plush.html"))
}

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

	c.Set("hackathon", hackathon)
	c.Set("project", project)
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
