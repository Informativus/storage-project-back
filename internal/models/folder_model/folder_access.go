package folder_model

import (
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
)

const (
	AccessTableName = "folder_access"
)

type FolderAccessModel struct {
	FolderID uuid.UUID        `json:"folder_id"`
	UserID   uuid.UUID        `json:"user_id"`
	RoleID   roles_model.Role `json:"role_id"`
}
