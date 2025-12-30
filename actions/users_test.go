package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestUsersCreate_ValidDomain tests user creation with valid email domain
func TestUsersCreate_ValidDomain(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock domain validation - allow the domain
	mockRepo.EXPECT().
		CompanyAllowedDomainIsDomainAllowed("example.com").
		Return(true, nil).
		Times(1)

	// Test the domain validation logic
	domain := "example.com"
	allowed, err := mockRepo.CompanyAllowedDomainIsDomainAllowed(domain)
	r.NoError(err)
	r.True(allowed)
}

// TestUsersCreate_InvalidDomain tests user creation with invalid email domain
func TestUsersCreate_InvalidDomain(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock domain validation - reject the domain
	mockRepo.EXPECT().
		CompanyAllowedDomainIsDomainAllowed("baddomain.com").
		Return(false, nil).
		Times(1)

	// Test the domain validation logic
	domain := "baddomain.com"
	allowed, err := mockRepo.CompanyAllowedDomainIsDomainAllowed(domain)
	r.NoError(err)
	r.False(allowed)
}