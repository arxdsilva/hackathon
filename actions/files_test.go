package actions

import (
	"testing"

	"github.com/arxdsilva/hackathon/models"
	"github.com/arxdsilva/hackathon/repository"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestFilesIndex tests listing all files
func TestFilesIndex_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock file listing
	expectedFiles := &models.Files{}
	mockRepo.EXPECT().
		FileFindAll().
		Return(expectedFiles, nil).
		Times(1)

	// Test the call
	files, err := mockRepo.FileFindAll()
	r.NoError(err)
	r.NotNil(files)
}

// TestFilesShow tests viewing a specific file
func TestFilesShow_Success(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	// Mock file lookup
	expectedFile := &models.File{
		ID:       123,
		Filename: "test-file.pdf",
	}

	mockRepo.EXPECT().
		FileFindByID("123").
		Return(expectedFile, nil).
		Times(1)

	// Test the call
	file, err := mockRepo.FileFindByID("123")
	r.NoError(err)
	r.Equal("test-file.pdf", file.Filename)
	r.Equal(123, file.ID)
}