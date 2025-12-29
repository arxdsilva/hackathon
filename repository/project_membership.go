package repository

import (
	"github.com/arxdsilva/hackathon/models"
	"github.com/gobuffalo/pop/v6"
)

// ProjectMembershipRepository handles project membership-related database operations
type ProjectMembershipRepository struct {
	*BaseRepository
}

// NewProjectMembershipRepository creates a new project membership repository
func NewProjectMembershipRepository(conn *pop.Connection) *ProjectMembershipRepository {
	return &ProjectMembershipRepository{
		BaseRepository: NewBaseRepository(conn),
	}
}

// FindByProjectIDAndUserID finds a membership by project and user ID
func (r *ProjectMembershipRepository) FindByProjectIDAndUserID(projectID, userID interface{}) (*models.ProjectMembership, error) {
	membership := &models.ProjectMembership{}
	err := r.conn.Where("project_id = ? AND user_id = ?", projectID, userID).First(membership)
	return membership, err
}

// CountByProjectID returns the number of memberships for a project
func (r *ProjectMembershipRepository) CountByProjectID(projectID interface{}) (int, error) {
	count, err := r.conn.Where("project_id = ?", projectID).Count(&models.ProjectMembership{})
	return count, err
}

// IsUserMember checks if a user is a member of a project
func (r *ProjectMembershipRepository) IsUserMember(projectID, userID interface{}) (bool, error) {
	count, err := r.conn.Where("project_id = ? AND user_id = ?", projectID, userID).Count(&models.ProjectMembership{})
	return count > 0, err
}
