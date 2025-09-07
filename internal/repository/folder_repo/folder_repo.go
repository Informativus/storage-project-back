package folder_repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database"
	"github.com/ivan/storage-project-back/pkg/database/sql_database"
)

type IFldRepo interface {
	CreateFld(fldModel *folder_model.FolderModel) (*folder_model.FolderModel, error)
	InsertFolderAccess(fldAccessModel folder_model.FolderAccessModel) (folder_model.FolderAccessModel, error)
	GetGeneralFolderByName(fldName string) (*folder_model.MainFolderModel, error)
	GetGeneralFolderById(fldID uuid.UUID) (*folder_model.MainFolderModel, error)
	GetGeneralFolderByUsrId(id uuid.UUID) (*folder_model.MainFolderModel, error)
	GetGeneralFolderBySubFldId(id uuid.UUID) (*folder_model.FolderModel, error)
	GetFldByNameAndMainFldId(fldName string, mainID uuid.UUID) (*folder_model.FolderModel, error)
	GetFldById(fldId uuid.UUID) (*folder_model.FolderModel, error)
	GetFldByIdAndMainFldId(fldID uuid.UUID, mainID uuid.UUID) (*folder_model.FolderModel, error)
	DelFld(id uuid.UUID) error
}

type FldRepo struct {
	db sql_database.DBClient
}

func NewFldRepo(db *database.DatabaseModule) *FldRepo {
	return &FldRepo{
		db: db.SQLDB,
	}
}

func (f *FldRepo) CreateFld(fldModel *folder_model.FolderModel) (*folder_model.FolderModel, error) {
	cals, vals, phs, err := sql_builder.InsertArgs(fldModel)

	if err != nil {
		return nil, err
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

	return &inserted, nil
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

func (f *FldRepo) GetGeneralFolderByName(fldName string) (*folder_model.MainFolderModel, error) {
	cals, err := sql_builder.GetStructCols(folder_model.MainFolderModel{})

	if err != nil {
		return nil, err
	}

	where := "name = $1"

	query := sql_builder.BuildSelectQuery(folder_model.MainUserFolderViewName, cals, &where)

	var selected folder_model.MainFolderModel
	err = f.db.QueryRow(context.Background(), query, fldName).Scan(
		&selected.UserID,
		&selected.FolderID,
		&selected.Name,
	)

	if err != nil && f.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}

func (f *FldRepo) GetGeneralFolderById(fldID uuid.UUID) (*folder_model.FolderModel, error) {
	cals, err := sql_builder.GetStructCols(folder_model.FolderModel{})

	if err != nil {
		return nil, err
	}

	where := "id = $1 and main_folder_id is null"

	query := sql_builder.BuildSelectQuery(folder_model.TableName, cals, &where)

	var selected folder_model.FolderModel
	err = f.db.QueryRow(context.Background(), query, fldID).Scan(
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

func (f *FldRepo) GetGeneralFolderByUsrId(id uuid.UUID) (*folder_model.MainFolderModel, error) {
	cals, err := sql_builder.GetStructCols(folder_model.MainFolderModel{})

	if err != nil {
		return nil, err
	}

	where := "user_id = $1"

	query := sql_builder.BuildSelectQuery(folder_model.MainUserFolderViewName, cals, &where)

	var selected folder_model.MainFolderModel
	err = f.db.QueryRow(context.Background(), query, id).Scan(
		&selected.UserID,
		&selected.FolderID,
		&selected.Name,
	)

	if err != nil && f.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}

func (f *FldRepo) GetGeneralFolderBySubFldId(id uuid.UUID) (*folder_model.FolderModel, error) {
	cals, err := sql_builder.GetStructCols(folder_model.FolderModel{})

	if err != nil {
		return nil, err
	}

	where := "id = (select main_folder_id from folders where id = $1)"

	query := sql_builder.BuildSelectQuery(folder_model.TableName, cals, &where)

	selected := folder_model.FolderModel{}
	err = f.db.QueryRow(context.Background(), query, id).Scan(
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

func (f *FldRepo) GetFldByNameAndMainFldId(fldName string, mainID uuid.UUID) (*folder_model.FolderModel, error) {
	cals, err := sql_builder.GetStructCols(folder_model.FolderModel{})

	if err != nil {
		return nil, err
	}

	where := "name = $1 and main_folder_id = $2"

	query := sql_builder.BuildSelectQuery(folder_model.TableName, cals, &where)

	selected := folder_model.FolderModel{}
	err = f.db.QueryRow(context.Background(), query, fldName, mainID).Scan(
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

func (f *FldRepo) GetFldByIdAndMainFldId(fldID uuid.UUID, mainID uuid.UUID) (*folder_model.FolderModel, error) {
	cals, err := sql_builder.GetStructCols(folder_model.FolderModel{})

	if err != nil {
		return nil, err
	}

	where := "id = $1 and main_folder_id = $2"

	query := sql_builder.BuildSelectQuery(folder_model.TableName, cals, &where)

	selected := folder_model.FolderModel{}
	err = f.db.QueryRow(context.Background(), query, fldID, mainID).Scan(
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

func (f *FldRepo) GetFldById(fldId uuid.UUID) (*folder_model.FolderModel, error) {
	cals, err := sql_builder.GetStructCols(folder_model.FolderModel{})

	if err != nil {
		return nil, err
	}

	where := "id = $1"

	query := sql_builder.BuildSelectQuery(folder_model.TableName, cals, &where)

	selected := folder_model.FolderModel{}
	err = f.db.QueryRow(context.Background(), query, fldId).Scan(
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

func (f *FldRepo) DelFld(id uuid.UUID) (int64, error) {
	query := sql_builder.BuildDeleteQuery(folder_model.TableName, "id = $1")

	tag, err := f.db.Exec(context.Background(), query, id)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}
