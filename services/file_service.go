package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/dto"
	"github.com/RenuBhati/editor/models"
)

const repoBasePath = "./repos"

func CreateFile(req dto.CreateFileRequest) (models.File, error) {
	newFile := models.File{
		OwnerID:    req.OwnerID,
		Name:       req.Name,
		Content:    req.Content,
		FileType:   "owned",
		GitHistory: []string{},
	}

	if err := database.DB.Create(&newFile).Error; err != nil {
		return newFile, err
	}

	// create repo
	repoPath := filepath.Join(repoBasePath, fmt.Sprintf("%d", newFile.ID))
	if err := os.MkdirAll(repoPath, os.ModePerm); err != nil {
		return newFile, err
	}

	/*
		1 cd repoPath
		2 copy content of file from newFile(content)
		3 open file and write
		4 git init
		5 git add .
		6 git commit
	*/
	filePath2 := filepath.Join(repoPath, newFile.Name)
	if err := os.WriteFile(filePath2, []byte(newFile.Content), 0644); err != nil {
		return newFile, err
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return newFile, err
	}

	cmd = exec.Command("git", "add", newFile.Name)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return newFile, err
	}
	cmd = exec.Command("git", "commit", "-m", "created file")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return newFile, err
	}
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return newFile, err
	}

	commitHash := strings.TrimSpace(string(output))
	history := []string{commitHash}
	newFile.GitHistory = history
	if err := database.DB.Save(&newFile).Error; err != nil {
		return newFile, err
	}
	return newFile, nil

}
