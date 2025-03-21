package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/dto"
	"github.com/RenuBhati/editor/models"
	"github.com/RenuBhati/editor/services"
	"github.com/gofiber/fiber/v2"
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

func cleanupTestDB() {
	os.Remove("test.db")
	os.RemoveAll("repos")
}

func TestCreateFileController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	app := fiber.New()
	app.Post("/files", CreateFiles)

	reqBody := `{"owner_id":1,"name":"testfile.txt","content":"Hello, Golang!"}`
	req := httptest.NewRequest(http.MethodPost, "/files", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var createdFile models.File
	json.NewDecoder(resp.Body).Decode(&createdFile)
	assert.NotZero(t, createdFile.ID)
}

func TestListFilesController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	services.CreateFile(dto.CreateFileRequest{
		OwnerID: 1, Name: "file1.txt", Content: "Content1",
	})

	app := fiber.New()
	app.Get("/files", ListFiles)

	req := httptest.NewRequest(http.MethodGet, "/files?user_id=1&page=1&limit=10", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetFileController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	createdFile, _ := services.CreateFile(dto.CreateFileRequest{
		OwnerID: 1, Name: "testfile.txt", Content: "Hello, Golang!",
	})

	app := fiber.New()
	app.Get("/files/:id", GetFile)

	reqURL := "/files/" + strconv.Itoa(int(createdFile.ID)) + "?user_id=1"
	req := httptest.NewRequest(http.MethodGet, reqURL, nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateFileController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	createdFile, _ := services.CreateFile(dto.CreateFileRequest{
		OwnerID: 1, Name: "testfile.txt", Content: "Hello, Golang!",
	})

	app := fiber.New()
	app.Put("/files/:id", UpdateFiles)

	reqBody := `{"content":"Updated Content","user_id":1}`
	req := httptest.NewRequest(http.MethodPut, "/files/"+strconv.Itoa(int(createdFile.ID)), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteFileController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	createdFile, _ := services.CreateFile(dto.CreateFileRequest{
		OwnerID: 1, Name: "test1.txt", Content: "Content",
	})

	app := fiber.New()
	app.Delete("/files/:id", DeleteFiles)

	reqURL := "/files/" + strconv.Itoa(int(createdFile.ID)) + "?user_id=1"
	req := httptest.NewRequest(http.MethodDelete, reqURL, nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestShareFileController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	createdFile, _ := services.CreateFile(dto.CreateFileRequest{
		OwnerID: 1, Name: "testfile.txt", Content: "Hello, Golang!",
	})

	app := fiber.New()
	app.Post("/files/:id/share", ShareFile)

	reqBody := `{"share_user_id":2}`
	req := httptest.NewRequest(http.MethodPost, "/files/"+strconv.Itoa(int(createdFile.ID))+"/share?user_id=1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSaveFileController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	createdFile, _ := services.CreateFile(dto.CreateFileRequest{
		OwnerID: 1, Name: "testfile.txt", Content: "Hello, Golang!",
	})

	app := fiber.New()
	app.Post("/files/:id/save", SaveFile)

	reqBody := `{"content":"Updated Content","user_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/files/"+strconv.Itoa(int(createdFile.ID))+"/save", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFileHistoryController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	createdFile, _ := services.CreateFile(dto.CreateFileRequest{
		OwnerID: 1, Name: "testfile.txt", Content: "Hello, Golang!",
	})

	app := fiber.New()
	app.Get("/files/:id/history", FileHistory)

	req := httptest.NewRequest(http.MethodGet, "/files/"+strconv.Itoa(int(createdFile.ID))+"/history?user_id=1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGitBlameController(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB()

	createdFile, _ := services.CreateFile(dto.CreateFileRequest{
		OwnerID: 1, Name: "testfile.txt", Content: "Hello, Golang!",
	})

	app := fiber.New()
	app.Get("/files/:id/blame", GitBlame)

	req := httptest.NewRequest(http.MethodGet, "/files/"+strconv.Itoa(int(createdFile.ID))+"/blame?user_id=1", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
