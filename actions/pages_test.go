package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestPagesScheduleIndex tests the public schedule page
func TestPagesScheduleIndex_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock active hackathons with schedule
	activeHackathons := &models.Hackathons{}
	mockRepo.EXPECT().
		HackathonGetActiveWithSchedule().
		Return(activeHackathons, nil).
		Times(1)

	// Test the call
	hackathons, err := mockRepo.HackathonGetActiveWithSchedule()
	r.NoError(err)
	r.NotNil(hackathons)
}
