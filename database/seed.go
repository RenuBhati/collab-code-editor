package database

import (
	"log"

	"github.com/RenuBhati/editor/models"
)

func SeedDB() {
	DB.Exec("DELETE FROM shared_files")
	DB.Exec("DELETE FROM files")

	files := []models.File{
		{
			Name:       "hello.go",
			OwnerID:    1,
			Content:    "package main\n\nfunc main() { println(\"Hello, World!\") }",
			FileType:   "owned",
			GitHistory: []string{},
		},
		{
			Name:       "app.js",
			OwnerID:    2,
			Content:    "console.log('Hello from app.js');",
			FileType:   "owned",
			GitHistory: []string{},
		},
		{
			Name:       "main.py",
			OwnerID:    3,
			Content:    "print('Hello from Python')",
			FileType:   "owned",
			GitHistory: []string{},
		},
		{
			Name:       "index.html",
			OwnerID:    1,
			Content:    "<html><body><h1>Hello HTML</h1></body></html>",
			FileType:   "owned",
			GitHistory: []string{},
		},
		{
			Name:       "script.rb",
			OwnerID:    2,
			Content:    "puts 'Hello from Ruby'",
			FileType:   "owned",
			GitHistory: []string{},
		},
	}

	for i := range files {
		if err := DB.Create(&files[i]).Error; err != nil {
			log.Printf("Error seeding file %s: %v", files[i].Name, err)
		} else {
			log.Printf("Seeded file %s with ID: %d", files[i].Name, files[i].ID)
		}
	}

	sharedRecords := []models.SharedFile{
		// Share app.js (file owned by user 2) with user 1 and user 3.
		{FileID: files[1].ID, UserID: 1},
		{FileID: files[1].ID, UserID: 3},
		// Share main.py (file owned by user 3) with user 1.
		{FileID: files[2].ID, UserID: 1},
		// Share index.html (file owned by user 1) with user 2 and user 3.
		{FileID: files[3].ID, UserID: 2},
		{FileID: files[3].ID, UserID: 3},
		// Share script.rb (file owned by user 2) with user 1, user 3, and user 4.
		{FileID: files[4].ID, UserID: 1},
		{FileID: files[4].ID, UserID: 3},
		{FileID: files[4].ID, UserID: 4},
	}

	for _, rec := range sharedRecords {
		if err := DB.Create(&rec).Error; err != nil {
			log.Printf("Error seeding shared record for file_id %d and user_id %d: %v", rec.FileID, rec.UserID, err)
		} else {
			log.Printf("Seeded shared record: file_id %d shared with user_id %d", rec.FileID, rec.UserID)
		}
	}
}
