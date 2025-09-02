package user_repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database/database"
	"github.com/rs/zerolog/log"
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
	log.Debug().Msg(query)

	var inserted user_model.UserModel
	err = ur.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.ID,
		&inserted.Name,
		&inserted.Blocked,
		&inserted.RoleID,
		&inserted.CreatedAt,
		&inserted.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return inserted, nil
}

func (ur *UserRepo) InsertUserToken(user user_model.UserTokensModel) (user_model.UserTokensModel, error) {
	cols, vals, phs, err := sql_builder.InsertArgs(user)

	if err != nil {
		return user_model.UserTokensModel{}, err
	}

	query := sql_builder.BuildInsertQuery(user_model.TokenTableName, cols, phs)

	var inserted user_model.UserTokensModel

	err = ur.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.ID,
		&inserted.UserID,
		&inserted.Token,
		&inserted.Revoked,
		&inserted.CreatedAt,
		&inserted.ExpiresAt,
	)

	if err != nil {
		return user_model.UserTokensModel{}, err
	}

	return inserted, nil
}

func (ur *UserRepo) DelUser(id uuid.UUID) (int64, error) {
	query := sql_builder.BuildDeleteQuery(user_model.TableName, "id = $1")

	tag, err := ur.db.Exec(context.Background(), query, id)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (ur *UserRepo) GetUserById(id uuid.UUID) (*user_model.UserModel, error) {
	cols, err := sql_builder.SelectArgs(user_model.UserModel{})

	if err != nil {
		return nil, err
	}

	where := "id = $1"

	query := sql_builder.BuildSelectQuery(user_model.TableName, cols, &where)

	selected := user_model.UserModel{}

	err = ur.db.QueryRow(context.Background(), query, id).Scan(
		&selected.ID,
		&selected.Name,
		&selected.Blocked,
		&selected.RoleID,
		&selected.CreatedAt,
		&selected.UpdatedAt,
	)

	if err != nil && ur.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil

}

func (ur *UserRepo) GetUserByName(name string) (*user_model.UserModel, error) {
	cols, err := sql_builder.SelectArgs(user_model.UserModel{})

	if err != nil {
		return nil, err
	}

	where := "name = $1"

	query := sql_builder.BuildSelectQuery(user_model.TableName, cols, &where)

	selected := user_model.UserModel{}

	err = ur.db.QueryRow(context.Background(), query, name).Scan(
		&selected.ID,
		&selected.Name,
		&selected.Blocked,
		&selected.RoleID,
		&selected.CreatedAt,
		&selected.UpdatedAt,
	)

	if err != nil && ur.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil

}

func (ur *UserRepo) GetUserAccessByToken(token string) (*user_model.UserTokensModel, error) {
	cols, err := sql_builder.SelectArgs(user_model.UserTokensModel{})

	if err != nil {
		return nil, err
	}

	where := "token = $1"
	query := sql_builder.BuildSelectQuery(user_model.TokenTableName, cols, &where)

	selected := user_model.UserTokensModel{}

	err = ur.db.QueryRow(context.Background(), query, token).Scan(
		&selected.ID,
		&selected.UserID,
		&selected.Token,
		&selected.Revoked,
		&selected.CreatedAt,
		&selected.ExpiresAt,
	)

	if err != nil && ur.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}

func (ur *UserRepo) DelExpiredTokens() (int64, error) {
	query := sql_builder.BuildDeleteQuery(user_model.TokenTableName, "expires_at < $1")

	tag, err := ur.db.Exec(context.Background(), query, time.Now())

	if err != nil {
		log.Error().Err(err).Msg("failed to delete expired tokens")
		return 0, err
	}

	return tag.RowsAffected(), nil
}
