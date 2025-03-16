package dto

type CreateFileRequest struct {
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
	OwnerID int    `json:"owner_id" validate:"required"`
}
