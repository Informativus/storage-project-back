package folder_model

import "github.com/google/uuid"

const (
	AccessTableName = "folder_access"
)

type FolderAccessModel struct {
	FolderID uuid.UUID `json:"folder_id"`
	UserID   uuid.UUID `json:"user_id"`
	RoleID   int8      `json:"role_id"`
}
