package file_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/file_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database"
	"github.com/ivan/storage-project-back/pkg/database/sql_database"
	"github.com/rs/zerolog/log"
)

type FileRepo struct {
	db sql_database.DBClient
}

func NewFileRepo(db *database.DatabaseModule) *FileRepo {
	return &FileRepo{
		db: db.SQLDB,
	}
}

func (f *FileRepo) UploadFile(fileModel *file_model.FileModel) (*file_model.FileModel, error) {
	cals, vals, phs, err := sql_builder.InsertArgs(fileModel)
	if err != nil {
		return nil, err
	}

	query := sql_builder.BuildInsertQuery(file_model.TableName, cals, phs)

	var inserted file_model.FileModel

	err = f.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.ID,
		&inserted.Name,
		&inserted.FolderID,
		&inserted.Size,
		&inserted.MimeType,
		&inserted.StorageKey,
		&inserted.CreatedAt,
		&inserted.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &inserted, nil

}

func (f *FileRepo) GetFileByID(id uuid.UUID) (*file_model.FileModel, error) {
	cals, err := sql_builder.GetStructCols(file_model.FileModel{})

	if err != nil {
		return nil, err
	}

	where := "id = $1"

	query := sql_builder.BuildSelectQuery(file_model.TableName, cals, &where)

	selected := file_model.FileModel{}

	err = f.db.QueryRow(context.Background(), query, id).Scan(
		&selected.ID,
		&selected.Name,
		&selected.FolderID,
		&selected.Size,
		&selected.MimeType,
		&selected.StorageKey,
		&selected.CreatedAt,
		&selected.DeletedAt,
	)

	if err != nil && f.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}

func (f *FileRepo) HardDelFile(id uuid.UUID) (int64, error) {
	where := "id = $1"

	query := sql_builder.BuildDeleteQuery(file_model.TableName, where)

	tags, err := f.db.Exec(context.Background(), query, id)

	if err != nil {
		return 0, err
	}

	return tags.RowsAffected(), nil
}

func (f *FileRepo) MarkFileAsDeleted(id uuid.UUID, deletedAt time.Time) (int64, error) {
	filds := []string{"deleted_at = $1"}

	where := "id = $2"

	query := sql_builder.BuildUpdateQuery(file_model.TableName, filds, where)

	tags, err := f.db.Exec(context.Background(), query, deletedAt, id)

	if err != nil {
		log.Error().Err(err).Msg("MarkFileAsDeleted")
		return 0, err
	}

	return tags.RowsAffected(), nil
}

func (f *FileRepo) GetMarkedToDelFiles() ([]file_model.FileModel, error) {
	cals, err := sql_builder.GetStructCols(file_model.FileModel{})

	if err != nil {
		return nil, err
	}

	where := "deleted_at IS NOT NULL"

	query := sql_builder.BuildSelectQuery(file_model.TableName, cals, &where)

	selected := []file_model.FileModel{}

	rows, err := f.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var fileModel file_model.FileModel
		err = rows.Scan(
			&fileModel.ID,
			&fileModel.Name,
			&fileModel.FolderID,
			&fileModel.Size,
			&fileModel.MimeType,
			&fileModel.StorageKey,
			&fileModel.CreatedAt,
			&fileModel.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		selected = append(selected, fileModel)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return selected, nil
}
