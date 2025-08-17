package user_repo

import (
	"context"

	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database/database"
)

type UserRepo struct {
	db database.DBClient
}

func NewUserRepo(db database.DBClient) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) CreateUser(user user_model.UserModel) (user_model.UserModel, error) {
	cols, vals, phs, err := sql_builder.InsertArgs(user)
	if err != nil {
		return user, err
	}

	query := sql_builder.BuildInsertQuery(user_model.TableName, cols, phs)

	var inserted user_model.UserModel
	err = ur.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.ID,
		&inserted.Token,
		&inserted.Blocked,
	)
	if err != nil {
		return user, err
	}

	return inserted, nil
}
