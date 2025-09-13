package user_model

import (
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
)

const (
	TableName = "users"
)

type UserModel struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	Blocked   bool             `json:"blocked"`
	RoleID    roles_model.Role `json:"role_id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt *time.Time       `json:"deleted_at"`
}

type UserDto struct {
	ID              uuid.UUID                      `json:"id"`
	Name            string                         `json:"name"`
	Blocked         bool                           `json:"blocked"`
	RoleID          roles_model.Role               `json:"role_id"`
	FoldersToAccess []folder_model.FolderAccessDto `json:"folders_to_access"`
	CreatedAt       time.Time                      `json:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at"`
	DeletedAt       *time.Time                     `json:"deleted_at"`
}
