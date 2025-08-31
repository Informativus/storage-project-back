package file_repo

import (
	"context"

	"github.com/ivan/storage-project-back/internal/models/file_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database/database"
)

type FileRepo struct {
	db database.DBClient
}

func NewFileRepo(db database.DBClient) *FileRepo {
	return &FileRepo{
		db: db,
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
	)

	if err != nil {
		return nil, err
	}

	return &inserted, nil

}
