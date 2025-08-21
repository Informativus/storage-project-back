package user_service

import (
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
)

type UserService struct {
	cfg           *config.Config
	jwt           *jwt_service.JwtService
	UserRepo      *user_repo.UserRepo
	FolderService *folder_service.FolderService
}

func NewUserService(
	cfg *config.Config,
	jwt *jwt_service.JwtService,
	userRepo *user_repo.UserRepo,
	fldService *folder_service.FolderService,
) *UserService {
	return &UserService{
		cfg:           cfg,
		jwt:           jwt,
		UserRepo:      userRepo,
		FolderService: fldService,
	}
}

func (u *UserService) GenerateToken(payload jwt_service.JwtPayload) (string, error) {
	token, err := u.jwt.GenerateToken(payload)

	if err != nil {
		log.Error().Err(err).Msg("failed to sign token")
		return "", err
	}

	return token, nil
}

func (u *UserService) CreateUser(fldName string, connUsrToFld bool) (string, error) {
	folderExt := u.FolderService.FolderExist(fldName)

	if folderExt {
		return "", errsvc.ErrFolderExist
	}

	usrID := uuid.New()

	token, err := u.GenerateToken(jwt_service.JwtPayload{
		ID: usrID,
	})

	if err != nil {
		return "", err
	}

	usrModel, err := u.UserRepo.CreateUser(user_model.UserModel{
		ID:      usrID,
		Token:   token,
		Blocked: false,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return "", err
	}

	err = u.FolderService.CreateFolder(fldName, usrModel.ID)

	if err != nil {
		log.Error().Err(err).Msg("failed to create folder")

		if delErr := u.DelUser(usrModel.ID); delErr != nil {
			log.Error().Err(delErr).Str("user_id", usrModel.ID.String()).
				Msg("rollback failed, inconsistent state: user exists without folder. You must delete the user yourself")
			return "", errsvc.ErrInconsistentState
		}
		return "", err
	}

	return usrModel.Token, nil
}

func (u *UserService) DelUser(id uuid.UUID) error {
	return u.UserRepo.DelUser(id)
}
