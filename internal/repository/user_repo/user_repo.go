package user_repo

import (
	"time"

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

func (ur *UserRepo) GetUserByIdCache(id uuid.UUID) (*user_model.UserDto, error) {
	cachedUser, err := ur.cacheDb.GetUsrById(id)

	if err != nil {
		return nil, err
	}

	return cachedUser, nil
}

func (ur *UserRepo) GetUserByIdDb(id uuid.UUID) (*user_model.UserModel, error) {
	usrModel, err := ur.db.GetUserById(id)

	if err != nil {
		return nil, err
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

	return updatedUsr, nil
}

func (ur *UserRepo) DelExpiredTokens() (int64, error) {
	return ur.db.DelExpiredTokens()
}

func (ur *UserRepo) MarkUsrAsDeleted(id uuid.UUID) (int64, error) {
	row, err := ur.db.MarkUsrAsDeleted(time.Now(), id)

	if err != nil {
		return 0, err
	}

	return row, nil
}

func (ur *UserRepo) GetMarkedToDelUsrs() ([]user_model.UserModel, error) {
	return ur.db.GetMarkedToDelUsrs()
}

func (ur *UserRepo) DelUser(id uuid.UUID) (int64, error) {
	return ur.db.DelUser(id)
}

func (ur *UserRepo) DelUsrFromCache(id uuid.UUID) error {
	return ur.cacheDb.DelUsrById(id)
}

func (ur *UserRepo) SetUsrInCache(usrDto *user_model.UserDto) error {
	return ur.cacheDb.SetUsrById(*usrDto)
}

func (ur *UserRepo) HardDel(id uuid.UUID) (int64, error) {
	err := ur.cacheDb.DelUsrById(id)

	if err != nil {
		return 1, err
	}

	return ur.db.DelUser(id)
}
