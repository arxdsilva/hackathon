package repository

import (
	"github.com/arxdsilva/hackathon/models"
	"github.com/gobuffalo/pop/v6"
)

// FileRepository handles file-related database operations
type FileRepository struct {
	*BaseRepository
}

// NewFileRepository creates a new file repository
func NewFileRepository(conn *pop.Connection) *FileRepository {
	return &FileRepository{
		BaseRepository: NewBaseRepository(conn),
	}
}

// FindByID finds a file by ID with eager loading
func (r *FileRepository) FindByID(id interface{}) (*models.File, error) {
	file := &models.File{}
	err := r.conn.Eager("User", "Hackathon", "Project").Find(file, id)
	return file, err
}

// FindAll finds all files
func (r *FileRepository) FindAll() (*models.Files, error) {
	files := &models.Files{}
	err := r.conn.All(files)
	return files, err
}

// FindAllHackathons finds all hackathons (for file upload context)
func (r *FileRepository) FindAllHackathons() (*models.Hackathons, error) {
	hackathons := &models.Hackathons{}
	err := r.conn.Where("status != ?", "hidden").All(hackathons)
	return hackathons, err
}

// FindAllProjects finds all projects (for file upload context)
func (r *FileRepository) FindAllProjects() (*models.Projects, error) {
	projects := &models.Projects{}
	err := r.conn.All(projects)
	return projects, err
}
