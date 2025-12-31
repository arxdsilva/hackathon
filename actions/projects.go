package actions

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/arxdsilva/hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// ProjectsImage serves the project image
func (a *MyApp) ProjectsImage(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if len(project.ImageData) == 0 || project.ImageContentType == nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("no image"))
	}

	c.Response().Header().Set("Content-Type", *project.ImageContentType)
	c.Response().Write(project.ImageData)
	return nil
}

// ProjectsUpdateImage updates only the project image
func (a *MyApp) ProjectsUpdateImage(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Only the project owner can update
	currentUser := c.Value("current_user").(models.User)
	if project.UserID == nil || *project.UserID != currentUser.ID {
		c.Flash().Add("danger", "You can only edit your own projects.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
	}

	// Handle image upload
	if uploadedImage, imageHeader, err := c.Request().FormFile("image"); err == nil {
		defer uploadedImage.Close()

		// Validate file size (max 5MB)
		if imageHeader.Size > 5*1024*1024 {
			c.Flash().Add("danger", "Image too large (max 5MB)")
			return c.Redirect(http.StatusFound, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
		}

		// Validate content type
		contentType := imageHeader.Header.Get("Content-Type")
		if contentType == "" || !strings.HasPrefix(contentType, "image/") {
			c.Flash().Add("danger", "Invalid image file")
			return c.Redirect(http.StatusFound, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
		}

		// Read image data
		imageData, err := io.ReadAll(uploadedImage)
		if err != nil {
			c.Flash().Add("danger", "Failed to read image")
			return c.Redirect(http.StatusFound, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
		}

		project.ImageData = imageData
		project.ImageContentType = &contentType

		verrs, err := tx.ValidateAndUpdate(project)
		if err != nil {
			return err
		}

		if verrs.HasAny() {
			c.Flash().Add("danger", "Failed to update image")
			return c.Redirect(http.StatusFound, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
		}

		c.Flash().Add("success", "Project image updated successfully!")
	} else {
		c.Flash().Add("danger", "No image selected")
	}

	return c.Redirect(http.StatusFound, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
}

// ProjectsIndex lists all projects for a hackathon
func (a *MyApp) ProjectsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	hackathon := &models.Hackathon{}

	if err := tx.Find(hackathon, c.Param("hackathon_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	projects := &models.Projects{}
	if err := tx.Where("hackathon_id = ?", hackathon.ID).Eager("User").All(projects); err != nil {
		return err
	}

	// Count memberships for each project
	memberCounts := make(map[string]int)
	for _, project := range *projects {
		count, err := tx.Where("project_id = ?", project.ID).Count(&models.ProjectMembership{})
		if err == nil {
			memberCounts[project.ID] = count
		}
	}

	c.Set("hackathon", hackathon)
	c.Set("projects", projects)
	c.Set("memberCounts", memberCounts)
	return c.Render(http.StatusOK, r.HTML("projects/index.plush.html"))
}

// ProjectsShow displays a single project
func (a *MyApp) ProjectsShow(c buffalo.Context) error {
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

	// Load files for this project
	files := &models.Files{}
	if err := tx.Where("project_id = ?", project.ID).All(files); err != nil {
		return err
	}

	// Load project members
	memberships := &models.ProjectMemberships{}
	if err := tx.Where("project_id = ?", project.ID).All(memberships); err != nil {
		return err
	}

	// Load users for the memberships
	userIDs := make([]interface{}, len(*memberships))
	for i, membership := range *memberships {
		userIDs[i] = membership.UserID
	}

	projectUsers := &models.Users{}
	if len(userIDs) > 0 {
		if err := tx.Where("id IN (?)", userIDs...).All(projectUsers); err != nil {
			return err
		}
	}

	// Check if current user is a member
	isMember := false
	if cu, ok := c.Value("current_user").(models.User); ok {
		count, err := tx.Where("project_id = ? AND user_id = ?", project.ID, cu.ID).Count(&models.ProjectMembership{})
		if err == nil && count > 0 {
			isMember = true
		}
	}

	c.Set("hackathon", hackathon)
	c.Set("project", project)
	c.Set("isProjectOwner", isOwner)
	c.Set("files", files)
	c.Set("projectUsers", projectUsers)
	c.Set("isMember", isMember)
	return c.Render(http.StatusOK, r.HTML("projects/show.plush.html"))
}

// ProjectsNew renders the form for creating a new project
func (a *MyApp) ProjectsNew(c buffalo.Context) error {
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
func (a *MyApp) ProjectsCreate(c buffalo.Context) error {
	project := &models.Project{}
	if err := c.Bind(project); err != nil {
		return err
	}

	// Set the user_id to current user
	currentUser := c.Value("current_user").(models.User)
	project.UserID = &currentUser.ID

	// Handle image upload
	if uploadedImage, imageHeader, err := c.Request().FormFile("image"); err == nil {
		defer uploadedImage.Close()

		// Validate file size (max 5MB)
		if imageHeader.Size > 5*1024*1024 {
			c.Flash().Add("danger", "Image too large (max 5MB)")
			return c.Redirect(http.StatusFound, "/hackathons/%s/projects/new", c.Param("hackathon_id"))
		}

		// Validate content type
		contentType := imageHeader.Header.Get("Content-Type")
		if contentType == "" || !strings.HasPrefix(contentType, "image/") {
			c.Flash().Add("danger", "Invalid image file")
			return c.Redirect(http.StatusFound, "/hackathons/%s/projects/new", c.Param("hackathon_id"))
		}

		// Read image data
		imageData, err := io.ReadAll(uploadedImage)
		if err != nil {
			c.Flash().Add("danger", "Failed to read image")
			return c.Redirect(http.StatusFound, "/hackathons/%s/projects/new", c.Param("hackathon_id"))
		}

		project.ImageData = imageData
		project.ImageContentType = &contentType
	}

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

	// Log project creation
	logAuditEvent(tx, c, &currentUser.ID, "create", "project", &project.ID, fmt.Sprintf("Project created: %s", project.Name))

	c.Flash().Add("success", "Project created successfully!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
}

// ProjectsEdit renders the form for editing a project
func (a *MyApp) ProjectsEdit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Only the project owner can edit
	currentUser := c.Value("current_user").(models.User)
	if project.UserID == nil || *project.UserID != currentUser.ID {
		c.Flash().Add("danger", "You can only edit your own projects.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
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
func (a *MyApp) ProjectsUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Only the project owner can update
	currentUser := c.Value("current_user").(models.User)
	if project.UserID == nil || *project.UserID != currentUser.ID {
		c.Flash().Add("danger", "You can only edit your own projects.")
		return c.Redirect(http.StatusSeeOther, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
	}

	// Handle image upload
	if uploadedImage, imageHeader, err := c.Request().FormFile("image"); err == nil {
		defer uploadedImage.Close()

		// Validate file size (max 5MB)
		if imageHeader.Size > 5*1024*1024 {
			c.Flash().Add("danger", "Image too large (max 5MB)")
			return c.Redirect(http.StatusFound, "/hackathons/%d/projects/%d/edit", project.HackathonID, project.ID)
		}

		// Validate content type
		contentType := imageHeader.Header.Get("Content-Type")
		if contentType == "" || !strings.HasPrefix(contentType, "image/") {
			c.Flash().Add("danger", "Invalid image file")
			return c.Redirect(http.StatusFound, "/hackathons/%d/projects/%d/edit", project.HackathonID, project.ID)
		}

		// Read image data
		imageData, err := io.ReadAll(uploadedImage)
		if err != nil {
			c.Flash().Add("danger", "Failed to read image")
			return c.Redirect(http.StatusFound, "/hackathons/%d/projects/%d/edit", project.HackathonID, project.ID)
		}

		project.ImageData = imageData
		project.ImageContentType = &contentType
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

	// Log project update
	logAuditEvent(tx, c, &currentUser.ID, "update", "project", &project.ID, fmt.Sprintf("Project updated: %s", project.Name))

	c.Flash().Add("success", "Project updated successfully!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
}

// ProjectsDestroy deletes a project from the DB
// ProjectsDestroy is disabled: projects are retained and cannot be deleted.

// ProjectsTogglePresenting toggles the presenting status of a project
func (a *MyApp) ProjectsTogglePresenting(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	project := &models.Project{}

	if err := tx.Find(project, c.Param("project_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Check if current user is the project owner
	isOwner := false
	if cu, ok := c.Value("current_user").(models.User); ok && project.UserID != nil {
		isOwner = *project.UserID == cu.ID
	}

	if !isOwner {
		return c.Error(http.StatusForbidden, fmt.Errorf("only project owner can toggle presenting status"))
	}

	// Toggle presenting status
	project.Presenting = !project.Presenting

	// Set or clear presentation order
	if project.Presenting {
		now := time.Now()
		project.PresentationOrder = &now
	} else {
		project.PresentationOrder = nil
	}

	// Update the project
	if err := tx.Update(project); err != nil {
		return err
	}

	// Log the action
	if cu, ok := c.Value("current_user").(models.User); ok {
		action := "set_presenting"
		if !project.Presenting {
			action = "unset_presenting"
		}
		logAuditEvent(tx, c, &cu.ID, action, "project", &project.ID, fmt.Sprintf("Project presenting status changed: %s", project.Name))
	}

	c.Flash().Add("success", "Project presenting status updated!")
	return c.Redirect(http.StatusSeeOther, "/hackathons/%s/projects/%s", project.HackathonID, project.ID)
}
