package models

type SharedFile struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	FileID uint `gorm:"not null" json:"file_id"`
	UserID int  `gorm:"not null" json:"user_id"`
}
