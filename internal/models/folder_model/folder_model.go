package folder_model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TableName     = "folders"
	FolderNameLen = 255
	PathLen       = 512
)

type FolderModel struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id"`
	OwnerID   uuid.UUID  `json:"owner_id"`
	MainFldId *uuid.UUID `json:"main_folder_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
