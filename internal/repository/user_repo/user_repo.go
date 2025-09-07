package user_repo

import (
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	db_user "github.com/ivan/storage-project-back/internal/repository/user_repo/db"
	"github.com/rs/zerolog/log"
)

type UserRepo struct {
	db      *db_user.SqlUserRepo
	cacheDb *db_user.CacheUserRepo
}

func NewUserRepo(cacheDb *db_user.CacheUserRepo, sqlDB *db_user.SqlUserRepo) *UserRepo {
	return &UserRepo{
		db:      sqlDB,
		cacheDb: cacheDb,
	}
}

func (ur *UserRepo) GetUserById(id uuid.UUID) (*user_model.UserModel, error) {
	cachedUser, err := ur.cacheDb.GetUsrById(id)

	if err != nil {
		return nil, err
	}

	if cachedUser != nil {
		return cachedUser, nil
	}

	usrModel, err := ur.db.GetUserById(id)

	if err != nil {
		return nil, err
	}

	if usrModel != nil {
		if err := ur.cacheDb.SetUsrById(*usrModel); err != nil {
			log.Error().Err(err).Msg("failed to cache user")
		}
	}

	return usrModel, nil
}

func (ur *UserRepo) GetUserByName(name string) (*user_model.UserModel, error) {
	return ur.db.GetUserByName(name)
}

func (ur *UserRepo) GetUserAccessByToken(token string) (*user_model.UserTokensModel, error) {
	return ur.db.GetUserAccessByToken(token)
}

func (ur *UserRepo) CreateUser(user user_model.UserModel) (*user_model.UserModel, error) {
	usrModel, err := ur.db.CreateUser(user)

	if err != nil {
		return nil, err
	}

	if usrModel != nil {
		if err := ur.cacheDb.SetUsrById(*usrModel); err != nil {
			log.Error().Err(err).Msg("failed to cache user")
		}
	}

	return usrModel, nil
}

func (ur *UserRepo) InsertUserToken(user user_model.UserTokensModel) (user_model.UserTokensModel, error) {
	return ur.db.InsertUserToken(user)
}

func (ur *UserRepo) UpdateBlockUserInf(blocked bool, id uuid.UUID) (*user_model.UserModel, error) {
	err := ur.cacheDb.DelUsrById(id)

	if err != nil {
		log.Error().Err(err).Msg("failed to delete user from cache")
		return nil, err
	}

	updatedUsr, err := ur.db.UpdateBlockUserInf(blocked, id)

	if err != nil {
		return nil, err
	}

	if updatedUsr != nil {
		if err := ur.cacheDb.SetUsrById(*updatedUsr); err != nil {
			log.Error().Err(err).Msg("failed to cache user")
		}
	}

	return updatedUsr, nil
}

func (ur *UserRepo) DelExpiredTokens() (int64, error) {
	return ur.db.DelExpiredTokens()
}

func (ur *UserRepo) DelUser(id uuid.UUID) (int64, error) {
	row, err := ur.db.DelUser(id)

	if err != nil {
		return 0, err
	}

	if err := ur.cacheDb.DelUsrById(id); err != nil {
		log.Error().Err(err).Msg("failed to delete user from cache")
	}
	return row, nil
}
