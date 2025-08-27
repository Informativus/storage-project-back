package folder_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database/database"
)

type FldRepo struct {
	db database.DBClient
}

func NewFldRepo(db database.DBClient) *FldRepo {
	return &FldRepo{
		db: db,
	}
}

func (f *FldRepo) CreateFld(fldModel folder_model.FolderModel) (folder_model.FolderModel, error) {
	cals, vals, phs, err := sql_builder.InsertArgs(fldModel)

	if err != nil {
		return fldModel, err
	}

	query := sql_builder.BuildInsertQuery(folder_model.TableName, cals, phs)

	var inserted folder_model.FolderModel

	err = f.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.ID,
		&inserted.Name,
		&inserted.ParentID,
		&inserted.OwnerID,
		&inserted.MainFldId,
		&inserted.CreatedAt,
		&inserted.UpdatedAt,
	)

	if err != nil {
		return fldModel, err
	}

	return inserted, nil
}

func (f *FldRepo) InsertFolderAccess(fldAccessModel folder_model.FolderAccessModel) (folder_model.FolderAccessModel, error) {
	cals, vals, phs, err := sql_builder.InsertArgs(fldAccessModel)

	if err != nil {
		return folder_model.FolderAccessModel{}, err
	}

	query := sql_builder.BuildInsertQuery(folder_model.AccessTableName, cals, phs)

	var inserted folder_model.FolderAccessModel

	err = f.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.FolderID,
		&inserted.UserID,
		&inserted.RoleID,
	)

	if err != nil {
		return folder_model.FolderAccessModel{}, err
	}

	return inserted, nil
}

func (f *FldRepo) GetGeneralFolderByName(fldName string) (*folder_model.FolderModel, error) {
	cals, err := sql_builder.SelectArgs(folder_model.FolderModel{})

	if err != nil {
		return nil, err
	}

	where := "name = $1"

	query := sql_builder.BuildSelectQuery(folder_model.MainUserFolderViewName, cals, &where)

	var selected folder_model.FolderModel
	err = f.db.QueryRow(context.Background(), query, fldName).Scan(
		&selected.ID,
		&selected.Name,
		&selected.ParentID,
		&selected.OwnerID,
		&selected.MainFldId,
		&selected.CreatedAt,
		&selected.UpdatedAt,
	)

	if err != nil && f.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}

func (f *FldRepo) GetFldByName(fldName string) (*folder_model.FolderModel, error) {
	cals, err := sql_builder.SelectArgs(folder_model.FolderModel{})

	if err != nil {
		return nil, err
	}

	where := "name = $1"

	query := sql_builder.BuildSelectQuery(folder_model.TableName, cals, &where)

	selected := folder_model.FolderModel{}
	err = f.db.QueryRow(context.Background(), query, fldName).Scan(
		&selected.ID,
		&selected.Name,
		&selected.ParentID,
		&selected.OwnerID,
		&selected.MainFldId,
		&selected.CreatedAt,
		&selected.UpdatedAt,
	)

	if err != nil && f.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}

func (f *FldRepo) DelFld(id uuid.UUID) error {
	query := sql_builder.BuildDeleteQuery(folder_model.TableName, "id = $1")

	tag, err := f.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return err
	}

	return nil
}
