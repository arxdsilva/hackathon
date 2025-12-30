// Package actions provides example tests demonstrating how to use repository mocks.
//
// To use mocks in your tests:
//
// 1. Create a mock controller: ctrl := gomock.NewController(t); defer ctrl.Finish()
// 2. Create a mock repository: mockRepo := repository.NewMockRepositoryInterface(ctrl)
// 3. Set up expectations: mockRepo.EXPECT().MethodName(args).Return(result, error).Times(1)
// 4. Inject the mock into your code (e.g., via dependency injection)
// 5. Call the code under test and verify expectations
//
// Example:
//
//	func TestMyHandler(t *testing.T) {
//		ctrl := gomock.NewController(t)
//		defer ctrl.Finish()
//
//		mockRepo := repository.NewMockRepositoryInterface(ctrl)
//		mockRepo.EXPECT().UserFindByID("123").Return(expectedUser, nil)
//
//		// Inject mockRepo into your handler/service
//		handler := NewHandler(mockRepo)
//
//		// Test your handler
//		result := handler.GetUser("123")
//		assert.Equal(t, expectedUser, result)
//	}
package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Example test showing how to use the repository mock
func TestAdminConfigIndex_WithMock(t *testing.T) {
	r := require.New(t)

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock some repository method calls that might be used in admin handlers
	mockRepo.EXPECT().
		UserCount().
		Return(10, nil).
		Times(1)

	// Actually call the method to satisfy the expectation
	count, err := mockRepo.UserCount()
	r.NoError(err)
	r.Equal(10, count)

	// This demonstrates the mock setup - in a real test you'd inject this
	// into your handler and test the actual logic
	r.NotNil(mockRepo)
}

// Example of mocking user operations
func TestUserOperations_WithMock(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock user count
	mockRepo.EXPECT().
		UserCount().
		Return(42, nil).
		Times(1)

	// Mock finding user by email
	expectedUser := &models.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockRepo.EXPECT().
		UserFindByEmail("test@example.com").
		Return(expectedUser, nil).
		Times(1)

	// Actually call the methods to satisfy the expectations
	count, err := mockRepo.UserCount()
	r.NoError(err)
	r.Equal(42, count)

	user, err := mockRepo.UserFindByEmail("test@example.com")
	r.NoError(err)
	r.Equal("test@example.com", user.Email)
	r.Equal("Test User", user.Name)
}

// Example of mocking hackathon operations
func TestHackathonOperations_WithMock(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock hackathon count
	mockRepo.EXPECT().
		HackathonCount().
		Return(5, nil).
		Times(1)

	// Mock finding hackathons by owner
	expectedHackathons := &models.Hackathons{}
	mockRepo.EXPECT().
		HackathonFindByOwnerID(gomock.Any()).
		Return(expectedHackathons, nil).
		Times(1)

	// Actually call the methods
	count, err := mockRepo.HackathonCount()
	r.NoError(err)
	r.Equal(5, count)

	hackathons, err := mockRepo.HackathonFindByOwnerID("owner-123")
	r.NoError(err)
	r.NotNil(hackathons)
}

// Example of mocking project operations
func TestProjectOperations_WithMock(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock project membership check
	mockRepo.EXPECT().
		ProjectIsUserMemberOfProject("project-123", "user-456").
		Return(true, nil).
		Times(1)

	// Mock finding projects by user
	expectedProjects := &models.Projects{}
	mockRepo.EXPECT().
		ProjectFindByUserID("user-456").
		Return(expectedProjects, nil).
		Times(1)

	// Actually call the methods
	isMember, err := mockRepo.ProjectIsUserMemberOfProject("project-123", "user-456")
	r.NoError(err)
	r.True(isMember)

	projects, err := mockRepo.ProjectFindByUserID("user-456")
	r.NoError(err)
	r.NotNil(projects)
}
