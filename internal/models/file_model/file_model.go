package file_model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TableName      = "files"
	MaxFileNameLen = 255
	MinFileNameLen = 3
)

type FileModel struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	FolderID   uuid.UUID `json:"folder_id"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	StorageKey string    `json:"storage_key"`
	CreatedAt  time.Time `json:"created_at"`
}
