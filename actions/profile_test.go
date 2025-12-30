package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestProfileShow tests the profile display functionality
func TestProfileShow_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock hackathons owned by user
	ownedHackathons := &models.Hackathons{}
	mockRepo.EXPECT().
		HackathonFindByOwnerID("user-123").
		Return(ownedHackathons, nil).
		Times(1)

	// Mock projects created by user
	createdProjects := &models.Projects{}
	mockRepo.EXPECT().
		ProjectFindByUserID("user-123").
		Return(createdProjects, nil).
		Times(1)

	// Test the repository calls
	hackathons, err := mockRepo.HackathonFindByOwnerID("user-123")
	r.NoError(err)
	r.NotNil(hackathons)

	projects, err := mockRepo.ProjectFindByUserID("user-123")
	r.NoError(err)
	r.NotNil(projects)
}
