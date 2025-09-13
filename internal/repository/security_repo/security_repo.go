package security_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/models/protection_group"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database"
	"github.com/ivan/storage-project-back/pkg/database/sql_database"
)

type SecurityRepo struct {
	db sql_database.DBClient
}

func NewSecurityRepo(db *database.DatabaseModule) *SecurityRepo {
	return &SecurityRepo{
		db: db.SQLDB,
	}
}

func (sr *SecurityRepo) GetUsrRoleForFolder(usrID uuid.UUID, fldID uuid.UUID) (*folder_model.FolderAccessModel, error) {
	col, err := sql_builder.GetStructCols(folder_model.FolderAccessModel{})

	if err != nil {
		return nil, err
	}

	where := "user_id = $1 and folder_id = $2"

	query := sql_builder.BuildSelectQuery(folder_model.AccessTableName, col, &where)

	selected := folder_model.FolderAccessModel{}

	err = sr.db.QueryRow(context.Background(), query, usrID, fldID).Scan(
		&selected.FolderID,
		&selected.UserID,
	)

	if err != nil && sr.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}

func (sr *SecurityRepo) InsertFolderAccess(fldAccessModel *folder_model.FolderAccessModel) (*folder_model.FolderAccessModel, error) {
	cals, vals, phs, err := sql_builder.InsertArgs(fldAccessModel)

	if err != nil {
		return nil, err
	}

	query := sql_builder.BuildInsertQuery(folder_model.AccessTableName, cals, phs)

	var inserted folder_model.FolderAccessModel

	err = sr.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.FolderID,
		&inserted.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &inserted, nil
}

func (sr *SecurityRepo) InsertFolderProtectionGroups(fldProtectGroupMdl *protection_group.FolderProtectionGroupsModel) (*protection_group.FolderProtectionGroupsModel, error) {
	cals, vals, phs, err := sql_builder.InsertArgs(fldProtectGroupMdl)

	if err != nil {
		return nil, err
	}

	query := sql_builder.BuildInsertQuery(protection_group.FolderProtectionGroupsName, cals, phs)

	var inserted protection_group.FolderProtectionGroupsModel

	err = sr.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.FolderID,
		&inserted.UserID,
		&inserted.GroupID,
	)

	if err != nil {
		return nil, err
	}

	return &inserted, nil
}

func (sr *SecurityRepo) GetUsrFoldersProtectionGroups(usrID uuid.UUID) ([]protection_group.FolderProtectionGroupsModel, error) {
	cols, err := sql_builder.GetStructCols(protection_group.FolderProtectionGroupsModel{})

	if err != nil {
		return nil, err
	}

	where := "user_id = $1"

	query := sql_builder.BuildSelectQuery(protection_group.FolderProtectionGroupsName, cols, &where)

	selected := []protection_group.FolderProtectionGroupsModel{}

	rows, err := sr.db.Query(context.Background(), query, usrID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var folderProtectionGroupsModel protection_group.FolderProtectionGroupsModel
		err = rows.Scan(
			&folderProtectionGroupsModel.FolderID,
			&folderProtectionGroupsModel.UserID,
			&folderProtectionGroupsModel.GroupID,
		)

		if err != nil {
			return nil, err
		}

		selected = append(selected, folderProtectionGroupsModel)
	}

	return selected, nil
}
