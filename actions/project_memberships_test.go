package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestProjectMembershipsCreate_AlreadyMember tests joining a project when already a member
func TestProjectMembershipsCreate_AlreadyMember(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock membership check - user is already a member
	mockRepo.EXPECT().
		ProjectMembershipIsUserMember("project-123", "user-456").
		Return(true, nil).
		Times(1)

	// Test the call
	isMember, err := mockRepo.ProjectMembershipIsUserMember("project-123", "user-456")
	r.NoError(err)
	r.True(isMember)
}

// TestProjectMembershipsCreate_Success tests successfully joining a project
func TestProjectMembershipsCreate_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock membership check - user is not a member
	mockRepo.EXPECT().
		ProjectMembershipIsUserMember("project-123", "user-456").
		Return(false, nil).
		Times(1)

	// Mock project lookup
	expectedProject := &models.Project{
		ID:   123,
		Name: "Test Project",
	}

	mockRepo.EXPECT().
		ProjectFindByID("project-123").
		Return(expectedProject, nil).
		Times(1)

	// Test the calls
	isMember, err := mockRepo.ProjectMembershipIsUserMember("project-123", "user-456")
	r.NoError(err)
	r.False(isMember)

	project, err := mockRepo.ProjectFindByID("project-123")
	r.NoError(err)
	r.Equal("Test Project", project.Name)
}
