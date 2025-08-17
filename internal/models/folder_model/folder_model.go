package folder_model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TableName              = "folders"
	MainUserFolderViewName = "main_user_folder"
	FolderNameLen          = 255
	PathLen                = 512
)

type FolderModel struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	UserID     uuid.UUID  `json:"user_id"`
	ParentID   *uuid.UUID `json:"parent_id"`
	Path       string     `json:"path"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUpdate time.Time  `json:"last_update"`
}
