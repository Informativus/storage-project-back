package folder_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database/database"
	"github.com/rs/zerolog/log"
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
		&inserted.UserID,
		&inserted.ParentID,
		&inserted.Path,
		&inserted.CreatedAt,
		&inserted.LastUpdate,
	)

	if err != nil {
		return fldModel, err
	}

	return inserted, nil
}

func (f *FldRepo) GetGeneralFolderByName(fldName string) (folder_model.FolderModel, error) {
	cals, err := sql_builder.SelectArgs(folder_model.FolderModel{})

	if err != nil {
		return folder_model.FolderModel{}, err
	}

	where := "name = $1"

	query := sql_builder.BuildSelectQuery(folder_model.MainUserFolderViewName, cals, &where)

	var selected folder_model.FolderModel
	err = f.db.QueryRow(context.Background(), query, fldName).Scan(
		&selected.ID,
		&selected.Name,
		&selected.UserID,
		&selected.ParentID,
		&selected.Path,
		&selected.CreatedAt,
		&selected.LastUpdate,
	)

	if err != nil && f.db.IsErrNoRows(err) {
		return folder_model.FolderModel{}, err
	}

	return selected, nil
}

func (f *FldRepo) DelMainFld(id uuid.UUID) error {
	query := sql_builder.BuildDeleteQuery(folder_model.TableName, "id = $1")

	tag, err := f.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		log.Error().Msg(fmt.Sprintf("failed to delete folder with id %s", id.String()))
		return errors.New("failed_delete")
	}

	return nil
}
