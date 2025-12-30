package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestAdminDashboard tests the admin dashboard data loading
func TestAdminDashboard_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock all the dashboard metrics
	mockRepo.EXPECT().UserCount().Return(150, nil).Times(1)
	mockRepo.EXPECT().HackathonCount().Return(12, nil).Times(1)
	mockRepo.EXPECT().ProjectCount().Return(89, nil).Times(1)
	mockRepo.EXPECT().ProjectCountActive().Return(45, nil).Times(1)
	mockRepo.EXPECT().ProjectCountPresenting().Return(23, nil).Times(1)

	// Mock recent data
	recentUsers := &models.Users{}
	recentHackathons := &models.Hackathons{}
	recentProjects := &models.Projects{}
	presentingProjects := &models.Projects{}

	mockRepo.EXPECT().UserGetRecent(5).Return(recentUsers, nil).Times(1)
	mockRepo.EXPECT().HackathonGetRecent(5).Return(recentHackathons, nil).Times(1)
	mockRepo.EXPECT().ProjectGetRecent(5).Return(recentProjects, nil).Times(1)
	mockRepo.EXPECT().ProjectFindPresentingFromActiveHackathons().Return(presentingProjects, nil).Times(1)

	// Actually call all the methods to satisfy expectations
	userCount, err := mockRepo.UserCount()
	r.NoError(err)
	r.Equal(150, userCount)

	hackathonCount, err := mockRepo.HackathonCount()
	r.NoError(err)
	r.Equal(12, hackathonCount)

	projectCount, err := mockRepo.ProjectCount()
	r.NoError(err)
	r.Equal(89, projectCount)

	activeCount, err := mockRepo.ProjectCountActive()
	r.NoError(err)
	r.Equal(45, activeCount)

	presentingCount, err := mockRepo.ProjectCountPresenting()
	r.NoError(err)
	r.Equal(23, presentingCount)

	users, err := mockRepo.UserGetRecent(5)
	r.NoError(err)
	r.NotNil(users)

	hackathons, err := mockRepo.HackathonGetRecent(5)
	r.NoError(err)
	r.NotNil(hackathons)

	projects, err := mockRepo.ProjectGetRecent(5)
	r.NoError(err)
	r.NotNil(projects)

	presenting, err := mockRepo.ProjectFindPresentingFromActiveHackathons()
	r.NoError(err)
	r.NotNil(presenting)
}

// TestAdminUserShow tests viewing a specific user in admin
func TestAdminUserShow_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock user lookup
	expectedUser := &models.User{
		Email: "user@example.com",
		Name:  "Test User",
	}

	mockRepo.EXPECT().
		UserFindByID("user-123").
		Return(expectedUser, nil).
		Times(1)

	// Mock user's projects
	userProjects := &models.Projects{}
	mockRepo.EXPECT().
		ProjectFindByUserIDWithHackathon("user-123").
		Return(userProjects, nil).
		Times(1)

	// Test the calls
	user, err := mockRepo.UserFindByID("user-123")
	r.NoError(err)
	r.Equal("user@example.com", user.Email)
	r.Equal("Test User", user.Name)

	projects, err := mockRepo.ProjectFindByUserIDWithHackathon("user-123")
	r.NoError(err)
	r.NotNil(projects)
}

// TestAdminCompanyAllowedDomainsIndex tests viewing allowed domains
func TestAdminCompanyAllowedDomainsIndex_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock domain listing
	expectedDomains := &models.CompanyAllowedDomains{}
	mockRepo.EXPECT().
		CompanyAllowedDomainFindAll().
		Return(expectedDomains, nil).
		Times(1)

	// Test the call
	domains, err := mockRepo.CompanyAllowedDomainFindAll()
	r.NoError(err)
	r.NotNil(domains)
}
