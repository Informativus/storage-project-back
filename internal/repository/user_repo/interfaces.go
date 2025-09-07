package user_repo

import (
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/user_model"
)

type IUserRepo interface {
	GetUserById(id uuid.UUID) (*user_model.UserModel, error)
	CreateUser(user user_model.UserModel) (*user_model.UserModel, error)
	GetUserAccessByToken(token string) (*user_model.UserTokensModel, error)
}
