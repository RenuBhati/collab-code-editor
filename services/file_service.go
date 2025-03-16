package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/dto"
	"github.com/RenuBhati/editor/models"
)

const repoBasePath = "./repos"

func createFile(req dto.CreateFileRequest) (models.File, error) {
	newFile := models.File{
		OwnerID:    req.OwnerID,
		Name:       req.Name,
		Content:    req.Content,
		FileType:   "owned",
		GitHistory: []string{},
	}

	myDb := database.DB
	if result := myDb.Create(&newFile); result.Error != nil {
		return newFile, result.Error
	}

	// create repo
	repoPath := filepath.Join(repoBasePath, fmt.Sprintf("%d", newFile.ID))
	if err := os.MkdirAll(repoPath, os.ModePerm); err != nil {
		return newFile, err
	}

	/*
		1 we will do cd repoPath
		2 copy content of file from newFile(content)
		3 open file and write
		4 git init
		5 git add .
		6 git commit 
	*/

}
