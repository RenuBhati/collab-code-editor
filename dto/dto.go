package dto

type CreateFileRequest struct {
	Name    string `json:"name" validate:"required,min=1"`
	Content string `json:"content" validate:"required"`
	OwnerID int    `json:"owner_id" validate:"required,gt=0"`
}

type UpdateFileRequest struct {
	Content string `json:"content" validate:"required"`
	UserID  int    `json:"user_id" validate:"required,gt=0"`
}

type SharedWithRequest struct {
	ShareUserID int `json:"share_user_id" validate:"required,gt=0"`
}

type SaveFileRequest struct {
	Content string `json:"content" validate:"required"`
	UserID  int    `json:"user_id" validate:"required,gt=0"`
}
