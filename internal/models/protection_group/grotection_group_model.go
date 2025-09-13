package protection_group

import (
	"github.com/google/uuid"
)

const (
	ProtectionGroupsName       = "protection_groups"
	ProtectionGroupActionsName = "protection_group_actions"
	FolderProtectionGroupsName = "folder_protection_groups"
)

type ProtectionGroupsModel struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Descritpion string    `json:"description"`
}

type ProtectionGroupActionsModel struct {
	GroupID  uuid.UUID `json:"group_id"`
	ActionID uuid.UUID `json:"action_id"`
}

type FolderProtectionGroupsModel struct {
	FolderID uuid.UUID `json:"folder_id"`
	UserID   uuid.UUID `json:"user_id"`
	GroupID  uuid.UUID `json:"protection_group_id"`
}

type ProtectionGroupsType string

const (
	Owner  ProtectionGroupsType = "Owner"
	Reader ProtectionGroupsType = "Reader"
	Editor ProtectionGroupsType = "Editor"
)

var ProtectionGroupIDs = map[ProtectionGroupsType]uuid.UUID{
	Owner:  uuid.MustParse("87d8ac99-6755-44e6-b4a9-cbd69ab31a4a"),
	Editor: uuid.MustParse("4d81a112-169b-4659-9d1f-d6745aa96f1c"),
	Reader: uuid.MustParse("5f9352a3-4d78-409c-b42b-eeba7fbd779c"),
}

var ProtectionGroupNames = map[uuid.UUID]ProtectionGroupsType{}

func init() {
	for name, id := range ProtectionGroupIDs {
		ProtectionGroupNames[id] = name
	}
}
