package folder_model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TableName     = "folders"
	FolderNameLen = 255
)

type FolderModel struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	UserID     uuid.UUID `json:"user_id"`
	ParentID   uuid.UUID `json:"parent_id"`
	CreatedAt  time.Time `json:"created_at"`
	LastUpdate time.Time `json:"last_update"`
}
