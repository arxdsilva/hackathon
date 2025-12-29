package repository

import (
	"github.com/arxdsilva/hackathon/models"
	"github.com/gobuffalo/pop/v6"
)

// ProjectRepository handles project-related database operations
type ProjectRepository struct {
	*BaseRepository
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(conn *pop.Connection) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: NewBaseRepository(conn),
	}
}

// Count returns the total number of projects
func (r *ProjectRepository) Count() (int, error) {
	var count int
	err := r.conn.RawQuery("SELECT COUNT(*) FROM projects").First(&count)
	return count, err
}

// CountActive returns the number of active projects
func (r *ProjectRepository) CountActive() (int, error) {
	var count int
	err := r.conn.RawQuery("SELECT COUNT(*) FROM projects WHERE status = 'active'").First(&count)
	return count, err
}

// CountPresenting returns the number of projects that are presenting from active hackathons
func (r *ProjectRepository) CountPresenting() (int, error) {
	var count int
	err := r.conn.RawQuery("SELECT COUNT(*) FROM projects WHERE presenting = true AND hackathon_id IN (SELECT id FROM hackathons WHERE status IN ('active', 'upcoming'))").First(&count)
	return count, err
}

// FindByID finds a project by ID
func (r *ProjectRepository) FindByID(id interface{}) (*models.Project, error) {
	project := &models.Project{}
	err := r.conn.Find(project, id)
	return project, err
}

// FindByHackathonID finds all projects for a specific hackathon
func (r *ProjectRepository) FindByHackathonID(hackathonID interface{}) (*models.Projects, error) {
	projects := &models.Projects{}
	err := r.conn.Where("hackathon_id = ?", hackathonID).Eager("User").All(projects)
	return projects, err
}

// FindByUserID finds all projects created by a specific user
func (r *ProjectRepository) FindByUserID(userID interface{}) (*models.Projects, error) {
	projects := &models.Projects{}
	err := r.conn.Where("user_id = ?", userID).All(projects)
	return projects, err
}

// FindByUserIDWithHackathon finds all projects created by a specific user with hackathon data
func (r *ProjectRepository) FindByUserIDWithHackathon(userID interface{}) (*models.Projects, error) {
	projects := &models.Projects{}
	err := r.conn.Eager("Hackathon").Where("user_id = ?", userID).All(projects)
	return projects, err
}

// FindPresentingByHackathonID finds presenting projects for a specific hackathon
func (r *ProjectRepository) FindPresentingByHackathonID(hackathonID interface{}) (*models.Projects, error) {
	projects := &models.Projects{}
	err := r.conn.Where("hackathon_id = ? AND presenting = ?", hackathonID, true).Order("presentation_order asc").Eager("User").All(projects)
	return projects, err
}

// FindPresentingFromActiveHackathons finds all presenting projects from active/upcoming hackathons
func (r *ProjectRepository) FindPresentingFromActiveHackathons() (*models.Projects, error) {
	projects := &models.Projects{}
	err := r.conn.Where("presenting = ? AND hackathon_id IN (SELECT id FROM hackathons WHERE status IN (?, ?))", true, "active", "upcoming").Order("presentation_order ASC").Eager("User", "Hackathon").All(projects)
	return projects, err
}

// GetRecent returns the most recently created projects (limited)
func (r *ProjectRepository) GetRecent(limit int) (*models.Projects, error) {
	projects := &models.Projects{}
	err := r.conn.Order("created_at DESC").Limit(limit).All(projects)
	return projects, err
}

// GetFilesByProjectID finds all files for a specific project
func (r *ProjectRepository) GetFilesByProjectID(projectID interface{}) (*models.Files, error) {
	files := &models.Files{}
	err := r.conn.Where("project_id = ?", projectID).All(files)
	return files, err
}

// GetMembershipsByProjectID finds all memberships for a specific project
func (r *ProjectRepository) GetMembershipsByProjectID(projectID interface{}) (*models.ProjectMemberships, error) {
	memberships := &models.ProjectMemberships{}
	err := r.conn.Where("project_id = ?", projectID).All(memberships)
	return memberships, err
}

// CountMembershipsByProjectID returns the number of memberships for a project
func (r *ProjectRepository) CountMembershipsByProjectID(projectID interface{}) (int, error) {
	count, err := r.conn.Where("project_id = ?", projectID).Count(&models.ProjectMembership{})
	return count, err
}

// IsUserMemberOfProject checks if a user is a member of a project
func (r *ProjectRepository) IsUserMemberOfProject(projectID, userID interface{}) (bool, error) {
	count, err := r.conn.Where("project_id = ? AND user_id = ?", projectID, userID).Count(&models.ProjectMembership{})
	return count > 0, err
}
