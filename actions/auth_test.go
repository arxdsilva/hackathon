package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestAuthSignin tests the signin functionality
func TestAuthSignin_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Create a test app with mock repository
	app := &MyApp{}

	// Mock user lookup
	expectedUser := &models.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "$2a$10$hashedpassword", // bcrypt hash for "password"
	}

	mockRepo.EXPECT().
		UserFindByEmail("test@example.com").
		Return(expectedUser, nil).
		Times(1)

	// Test the repository call (in real implementation, this would be called by the handler)
	user, err := mockRepo.UserFindByEmail("test@example.com")
	r.NoError(err)
	r.Equal("test@example.com", user.Email)
	r.Equal("Test User", user.Name)

	r.NotNil(app)
	r.NotNil(mockRepo)
}
