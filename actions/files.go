package actions

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"hackathon/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// FilesIndex lists all files
func FilesIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	files := &models.Files{}
	if err := tx.All(files); err != nil {
		return err
	}

	c.Set("files", files)
	return c.Render(http.StatusOK, r.HTML("files/index.plush.html"))
}

// FilesShow displays a single file
func FilesShow(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	file := &models.File{}
	if err := tx.Eager("User", "Hackathon", "Project").Find(file, c.Param("file_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("file", file)
	return c.Render(http.StatusOK, r.HTML("files/show.plush.html"))
}

// FilesNew renders the upload form
func FilesNew(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	// Get all hackathons for the dropdown
	hackathons := &models.Hackathons{}
	if err := tx.All(hackathons); err != nil {
		return err
	}

	// Get all projects for the dropdown
	projects := &models.Projects{}
	if err := tx.All(projects); err != nil {
		return err
	}

	c.Set("file", models.File{})
	c.Set("hackathons", hackathons)
	c.Set("projects", projects)
	return c.Render(http.StatusOK, r.HTML("files/new.plush.html"))
}

// FilesCreate handles file upload
func FilesCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(models.User)

	// Get the uploaded file
	uploadedFile, fileHeader, err := c.Request().FormFile("file")
	if err != nil {
		c.Flash().Add("danger", "No file uploaded")
		return c.Redirect(http.StatusFound, "/files/new")
	}
	defer uploadedFile.Close()

	// Validate file size (max 10MB)
	if fileHeader.Size > 10*1024*1024 {
		c.Flash().Add("danger", "File too large (max 10MB)")
		return c.Redirect(http.StatusFound, "/files/new")
	}

	// Read file data
	data, err := io.ReadAll(uploadedFile)
	if err != nil {
		c.Flash().Add("danger", "Failed to read file")
		return c.Redirect(http.StatusFound, "/files/new")
	}

	// Create file record
	fileRecord := &models.File{
		Filename:    fileHeader.Filename,
		Data:        data,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Size:        int(fileHeader.Size),
		UserID:      user.ID,
	}

	// Check for optional associations
	if hackathonID := c.Request().FormValue("hackathon_id"); hackathonID != "" {
		if id, err := strconv.Atoi(hackathonID); err == nil {
			fileRecord.HackathonID = &id
		}
	}

	if projectID := c.Request().FormValue("project_id"); projectID != "" {
		if id, err := strconv.Atoi(projectID); err == nil {
			fileRecord.ProjectID = &id
		}
	}

	verrs, err := tx.ValidateAndCreate(fileRecord)
	if err != nil {
		c.Flash().Add("danger", "Failed to create file record")
		return c.Redirect(http.StatusFound, "/files/new")
	}

	if verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("file", fileRecord)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("files/new.plush.html"))
	}

	c.Flash().Add("success", "File uploaded successfully")
	return c.Redirect(http.StatusFound, "/files")
}

// FilesDownload serves the file for download
func FilesDownload(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	file := &models.File{}
	if err := tx.Find(file, c.Param("file_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Filename))
	c.Response().Header().Set("Content-Type", file.ContentType)
	c.Response().Header().Set("Content-Length", strconv.FormatInt(int64(file.Size), 10))
	c.Response().Write(file.Data)
	return nil
}

// FilesDestroy deletes a file
func FilesDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	file := &models.File{}
	if err := tx.Find(file, c.Param("file_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Check ownership
	user := c.Value("current_user").(models.User)
	if file.UserID != user.ID && !user.IsOwner() {
		return c.Error(http.StatusForbidden, fmt.Errorf("not authorized"))
	}

	if err := tx.Destroy(file); err != nil {
		c.Flash().Add("danger", "Failed to delete file record")
		return c.Redirect(http.StatusFound, "/files")
	}

	c.Flash().Add("success", "File deleted successfully")
	return c.Redirect(http.StatusFound, "/files")
}
