package dto

type CreateFileRequest struct {
	Name    string `json:"name" validate:"required,min=1"`
	Content string `json:"content" validate:"required"`
	OwnerID int    `json:"owner_id" validate:"required,gt=0"`
}
