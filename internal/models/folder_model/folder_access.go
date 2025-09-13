package folder_model

import (
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/protection_group"
)

const (
	AccessTableName = "folder_access"
)

type FolderAccessModel struct {
	FolderID uuid.UUID `json:"folder_id"`
	UserID   uuid.UUID `json:"user_id"`
}

type FolderAccessDto struct {
	FolderID         uuid.UUID                               `json:"folder_id"`
	FolderName       string                                  `json:"folder_name"`
	ProtectionGroups []protection_group.ProtectionGroupsType `json:"protection_groups"`
}
