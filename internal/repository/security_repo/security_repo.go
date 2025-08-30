package security_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database/database"
)

type SecurityRepo struct {
	db database.DBClient
}

func NewSecurityRepo(db database.DBClient) *SecurityRepo {
	return &SecurityRepo{
		db: db,
	}
}

func (s *SecurityRepo) GetUsrRoleForFolder(usrID uuid.UUID, fldID uuid.UUID) (*folder_model.FolderAccessModel, error) {
	col, err := sql_builder.SelectArgs(folder_model.FolderAccessModel{})

	if err != nil {
		return nil, err
	}

	where := "user_id = $1 and folder_id = $2"

	query := sql_builder.BuildSelectQuery(folder_model.AccessTableName, col, &where)

	selected := folder_model.FolderAccessModel{}

	err = s.db.QueryRow(context.Background(), query, usrID, fldID).Scan(
		&selected.FolderID,
		&selected.UserID,
		&selected.RoleID,
	)

	if err != nil && s.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}
