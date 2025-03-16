package models

import "time"

type File struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name"`
	OwnerID    uint      `json:"owner_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	SharedWith []uint    `json:"shared_with" gorm:"type:json"`
	FileType   string    `json:"file_type"`
	GitHistory []string  `json:"git_history" gorm:"type:json"`
}
