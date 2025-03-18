package services

import (
	"encoding/json"
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
		GitHistory: "",
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
	historyBytes, err := json.Marshal(history)
	if err != nil {
		return newFile, err
	}
	newFile.GitHistory = string(historyBytes)
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

func UpdatedFile(fileID int, req dto.UpdateFileRequest) (models.File, error) {

	file, err := GetFileByID(fileID)
	if err != nil {
		return file, err
	}

	if !hasAccess(file, req.UserID) {
		return file, errors.New("unauthorized access")
	}

	file.Content = req.Content

	if err := database.DB.Save(&file).Error; err != nil {
		return file, err
	}

	repoPath := filepath.Join(repoBasePath, fmt.Sprintf("%d", file.ID))
	if err := os.MkdirAll(repoPath, os.ModePerm); err != nil {
		return file, err
	}
	filePath2 := filepath.Join(repoPath, file.Name)
	if err := os.WriteFile(filePath2, []byte(file.Content), 0644); err != nil {
		return file, err
	}

	cmd := exec.Command("git", "add", file.Name)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return file, err
	}
	cmd = exec.Command("git", "commit", "-m", "updated file")
	cmd.Dir = repoPath
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GIT_AUTHOR_NAME=User %d", req.UserID),
		fmt.Sprintf("GIT_AUTHOR_EMAIL=user%d@example.com", req.UserID),
	)
	if err := cmd.Run(); err != nil {
		return file, err
	}
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return file, err
	}

	commitHash := strings.TrimSpace(string(output))

	if err := appendCommitHash(&file, commitHash); err != nil {
		return file, err
	}

	return file, nil

	// file entry exists in DB else return error

	// if ownerId and userId are the same else return check in shared table entry else return error "unauthorised user"

	//use GetFile

	//file.content = req.content

	//update row update with new content - database.DB.Save

	//repoPath := filepath.Join(repoBasePath, fmt.Sprintf("%d", newFile.ID))

	/*filePath2 := filepath.Join(repoPath, newFile.Name)
	  if err := os.WriteFile(filePath2, []byte(newFile.Content), 0644); err != nil {
	  	return newFile, err
	  }
	*/

	// git add fileName

	//git commit -m "updated file", set env (author just userId)

	// git rev-parse

	//newFile.GitHistory= append(newFile.Githistory,commitHAsh).....I dont need to create new hostory array here

	//update DB with Save
}

func appendCommitHash(file *models.File, newHash string) error {
	var history []string
	if file.GitHistory != "" {
		if err := json.Unmarshal([]byte(file.GitHistory), &history); err != nil {
			history = []string{}
		}
	}
	history = append(history, newHash)
	bytes, err := json.Marshal(history)
	if err != nil {
		return err
	}
	file.GitHistory = string(bytes)
	return database.DB.Model(file).Update("git_history", file.GitHistory).Error

}

func DeleteFiles(userID, fileID int) error {

	file, err := GetFileByID(fileID)
	if err != nil {
		return err
	}

	if file.OwnerID != userID {
		return errors.New("unauthorized: only owner can delete file")
	}

	if err := database.DB.Delete(&models.File{}, fileID).Error; err != nil {
		return err
	}

	repoPath := filepath.Join(repoBasePath, fmt.Sprintf("%d", file.ID))
	return os.RemoveAll(repoPath)

}
