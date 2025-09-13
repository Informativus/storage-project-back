package db_user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/utils/sql_builder"
	"github.com/ivan/storage-project-back/pkg/database"
	"github.com/ivan/storage-project-back/pkg/database/no_sql_database"
	"github.com/ivan/storage-project-back/pkg/database/sql_database"
	"github.com/rs/zerolog/log"
)

type SqlUserRepo struct {
	db      sql_database.DBClient
	cacheDb *no_sql_database.RedisClient
}

func NewSqlUserRepo(db *database.DatabaseModule) *SqlUserRepo {
	return &SqlUserRepo{
		db: db.SQLDB,
	}
}

func (ur *SqlUserRepo) CreateUser(user user_model.UserModel) (*user_model.UserModel, error) {
	cols, vals, phs, err := sql_builder.InsertArgs(user)
	if err != nil {
		return nil, err
	}

	query := sql_builder.BuildInsertQuery(user_model.TableName, cols, phs)

	var inserted user_model.UserModel
	err = ur.db.QueryRow(context.Background(), query, vals...).Scan(
		&inserted.ID,
		&inserted.Name,
		&inserted.Blocked,
		&inserted.RoleID,
		&inserted.CreatedAt,
		&inserted.UpdatedAt,
		&inserted.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &inserted, nil
}

func (ur *SqlUserRepo) InsertUserToken(user user_model.UserTokensModel) (user_model.UserTokensModel, error) {
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

func (ur *SqlUserRepo) GetUserById(id uuid.UUID) (*user_model.UserModel, error) {
	cols, err := sql_builder.GetStructCols(user_model.UserModel{})

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
		&selected.DeletedAt,
	)

	if err != nil && ur.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil
}

func (ur *SqlUserRepo) GetUserByName(name string) (*user_model.UserModel, error) {
	cols, err := sql_builder.GetStructCols(user_model.UserModel{})

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
		&selected.DeletedAt,
	)

	if err != nil && ur.db.IsErrNoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &selected, nil

}

func (ur *SqlUserRepo) GetUserAccessByToken(token string) (*user_model.UserTokensModel, error) {
	cols, err := sql_builder.GetStructCols(user_model.UserTokensModel{})

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

func (ur *SqlUserRepo) UpdateBlockUserInf(blocked bool, id uuid.UUID) (*user_model.UserModel, error) {
	cols, err := sql_builder.GetStructCols(user_model.UserModel{})

	if err != nil {
		return nil, err
	}

	setClauses := sql_builder.BuildSetClauses([]string{"blocked"})

	query := sql_builder.BuildUpdateQueryReturn(user_model.TableName, setClauses, fmt.Sprintf("id = $%d", len(setClauses)+1), cols)

	var updated user_model.UserModel

	err = ur.db.QueryRow(context.Background(), query, blocked, id).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Blocked,
		&updated.RoleID,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	)

	return &updated, err
}

func (ur *SqlUserRepo) DelExpiredTokens() (int64, error) {
	query := sql_builder.BuildDeleteQuery(user_model.TokenTableName, "expires_at < $1")

	tag, err := ur.db.Exec(context.Background(), query, time.Now())

	if err != nil {
		log.Error().Err(err).Msg("failed to delete expired tokens")
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (ur *SqlUserRepo) MarkUsrAsDeleted(deletedAt time.Time, id uuid.UUID) (int64, error) {
	filds := []string{"deleted_at = $1"}

	where := "id = $2"

	query := sql_builder.BuildUpdateQuery(user_model.TableName, filds, where)

	tags, err := ur.db.Exec(context.Background(), query, deletedAt, id)

	if err != nil {
		return 0, err
	}

	return tags.RowsAffected(), nil
}

func (ur *SqlUserRepo) DelUser(id uuid.UUID) (int64, error) {
	query := sql_builder.BuildDeleteQuery(user_model.TableName, "id = $1")

	tag, err := ur.db.Exec(context.Background(), query, id)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (ur *SqlUserRepo) GetMarkedToDelUsrs() ([]user_model.UserModel, error) {
	cols, err := sql_builder.GetStructCols(user_model.UserModel{})

	if err != nil {
		return nil, err
	}

	where := "deleted_at IS NOT NULL"

	query := sql_builder.BuildSelectQuery(user_model.TableName, cols, &where)

	rows, err := ur.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var selected []user_model.UserModel
	for rows.Next() {
		var usr user_model.UserModel
		err = rows.Scan(
			&usr.ID,
			&usr.Name,
			&usr.Blocked,
			&usr.RoleID,
			&usr.CreatedAt,
			&usr.UpdatedAt,
			&usr.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		selected = append(selected, usr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return selected, nil
}
