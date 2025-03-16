package dto

type CreateFileRequest struct {
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
}

