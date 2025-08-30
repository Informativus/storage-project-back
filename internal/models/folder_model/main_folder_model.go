package folder_model

import (
	"github.com/google/uuid"
)

const (
	MainUserFolderViewName = "main_user_folder"
)

type MainFolderModel struct {
	UserID   uuid.UUID `json:"user_id"`
	FolderID uuid.UUID `json:"folder_id"`
	Name     string    `json:"name"`
}
