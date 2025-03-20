package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/dto"
	"github.com/RenuBhati/editor/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	database.DB.AutoMigrate(&models.File{}, &models.SharedFile{})
}

func cleanupTest() {
	os.Remove("test.db")
	os.RemoveAll("repos")
}

func cleanupRepos(t *testing.T, fileID uint) {
	repoPath := filepath.Join(repoBasePath, "repos", fmt.Sprintf("%d", fileID))
	if err := os.RemoveAll(repoPath); err != nil {
		t.Fatalf("Failed to clean up repository directory: %v", err)
	}
}

func TestCreateFile(t *testing.T) {
	setupTestDB(t)
	defer cleanupTest()

	req := dto.CreateFileRequest{
		OwnerID: 1,
		Name:    "testfile.txt",
		Content: "Hello, Golang!",
	}

	createdFile, err := CreateFile(req)
	assert.NoError(t, err)
	assert.NotZero(t, createdFile.ID)
	assert.Equal(t, req.Name, createdFile.Name)
	assert.Equal(t, req.Content, createdFile.Content)

	cleanupRepos(t, createdFile.ID)
}

func TestGetFileByID(t *testing.T) {
	setupTestDB(t)
	defer cleanupTest()

	req := dto.CreateFileRequest{
		OwnerID: 1,
		Name:    "testfile.txt",
		Content: "Hello, Golang!",
	}
	createdFile, _ := CreateFile(req)

	file, err := GetFileByID(int(createdFile.ID))
	assert.NoError(t, err)
	assert.Equal(t, createdFile.ID, file.ID)
}

func TestListFiles(t *testing.T) {
	setupTestDB(t)
	defer cleanupTest()

	req := dto.CreateFileRequest{OwnerID: 1, Name: "test1.txt", Content: "Content"}
	CreateFile(req)
	CreateFile(req)

	files, total, err := ListFiles(1, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, int(total))
	assert.Len(t, files, 2)
}

func TestDeleteFiles(t *testing.T) {
	setupTestDB(t)
	defer cleanupTest()

	req := dto.CreateFileRequest{OwnerID: 1, Name: "test1.txt", Content: "Content"}
	createdFile, _ := CreateFile(req)

	err := DeleteFiles(1, int(createdFile.ID))
	assert.NoError(t, err)
}

func TestUpdateFiles(t *testing.T) {
	setupTestDB(t)
	defer cleanupTest()

	req := dto.CreateFileRequest{OwnerID: 1, Name: "test1.txt", Content: "Content"}
	createdFile, _ := CreateFile(req)

	updateReq := dto.UpdateFileRequest{UserID: 1, Content: "Updated content"}
	updatedFile, err := UpdatedFile(int(createdFile.ID), updateReq)

	assert.NoError(t, err)
	assert.Equal(t, "Updated content", updatedFile.Content)
}

func TestShareFile(t *testing.T) {
	setupTestDB(t)
	defer cleanupTest()

	req := dto.CreateFileRequest{OwnerID: 1, Name: "test.txt", Content: "Hello"}
	createdFile, _ := CreateFile(req)

	shareReq := dto.SharedWithRequest{ShareUserID: 2}
	_, err := ShareFile(int(createdFile.ID), 1, shareReq)
	assert.NoError(t, err)
}

func TestGetFileHistory(t *testing.T) {
	setupTestDB(t)
	defer cleanupTest()

	req := dto.CreateFileRequest{OwnerID: 1, Name: "test.txt", Content: "Hello"}
	createdFile, _ := CreateFile(req)

	history, err := GetFileHistory(int(createdFile.ID), 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, history)
}

func TestGetGitBlame(t *testing.T) {
	setupTestDB(t)
	defer cleanupTest()

	req := dto.CreateFileRequest{OwnerID: 1, Name: "test.txt", Content: "Hello"}
	createdFile, _ := CreateFile(req)

	blame, err := GetGitBlame(int(createdFile.ID), 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, blame)
}
