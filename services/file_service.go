package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RenuBhati/editor/database"
	"github.com/RenuBhati/editor/dto"
	"github.com/RenuBhati/editor/models"
	"gorm.io/gorm"
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
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GIT_AUTHOR_NAME=User %d", req.OwnerID),
		fmt.Sprintf("GIT_AUTHOR_EMAIL=user%d@example.com", req.OwnerID),
	)
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

func ListFiles(userID, page, limit int) ([]models.File, int64, error) {
	var files []models.File
	var total int64
	offSet := (page - 1) * limit
	/*
			SELECT COUNT(*)
		FROM files
		WHERE owner_id = {userID}
		OR id IN (
		    SELECT file_id
		    FROM shared_files
		    WHERE user_id = {userID}
		);
	*/

	err := database.DB.Model(&models.File{}).
		Where("owner_id = ?", userID).
		Or("id IN (?)", database.DB.Model(&models.SharedFile{}).Select("file_id").Where("user_id = ?", userID)).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	/*
			SELECT *
		FROM files
		WHERE owner_id = {userID}
		OR id IN (
		    SELECT file_id
		    FROM shared_files
		    WHERE user_id = {userID}
		)
		LIMIT {limit} OFFSET {offset};
	*/

	err = database.DB.Model(&models.File{}).
		Where("owner_id = ?", userID).
		Or("id IN (?)", database.DB.Model(&models.SharedFile{}).Select("file_id").Where("user_id = ?", userID)).
		Offset(offSet).Limit(limit).
		Find(&files).Error
	if err != nil {
		return nil, 0, err
	}

	return files, total, nil

}

// GetFileByID retrieves a file by its ID.
func GetFileByID(fileID int) (models.File, error) {
	var file models.File
	err := database.DB.First(&file, fileID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return file, errors.New("file not found")
		}
		return file, err
	}
	return file, nil
}

// hasAccess checks whether the user is the owner
//
//	or the file is shared with the user.
func hasAccess(file models.File, userID int) bool {
	if file.OwnerID == userID {
		return true
	}
	var shared models.SharedFile
	err := database.DB.
		Where("file_id = ? AND user_id = ?", file.ID, userID).
		First(&shared).Error
	return err == nil
}

func GetFile(userID, fileID int) (models.File, error) {
	file, err := GetFileByID(fileID)
	if err != nil {
		return file, err
	}
	if !hasAccess(file, userID) {
		return file, errors.New("unauthorized access")
	}
	return file, nil
}
