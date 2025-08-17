package user_service

import (
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	user_repo "github.com/ivan/storage-project-back/internal/repository/user"
	"github.com/ivan/storage-project-back/internal/services/file_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
)

type UserService struct {
	cfg         *config.Config
	jwt         *jwt_service.JwtService
	UserRepo    *user_repo.UserRepo
	FileService *file_service.FileService
}

func NewUserService(
	cfg *config.Config,
	jwt *jwt_service.JwtService,
	userRepo *user_repo.UserRepo,
	fileSrvc *file_service.FileService,
) *UserService {
	return &UserService{
		cfg:         cfg,
		jwt:         jwt,
		UserRepo:    userRepo,
		FileService: fileSrvc,
	}
}

func (u *UserService) GenerateToken(folderName string) (string, error) {
	token, err := u.jwt.GenerateToken(jwt_service.JwtPayload{FolderName: folderName})

	if err != nil {
		log.Error().Err(err).Msg("failed to sign token")
		return "", err
	}

	return token, nil
}

func (u *UserService) CreateUser(fldName string) (string, error) {
	folderExt := u.FileService.FolderExist(fldName)

	if folderExt {
		return "", errsvc.ErrFolderExist
	}

	token, err := u.GenerateToken(fldName)

	if err != nil {
		return "", err
	}

	usrModel, err := u.UserRepo.CreateUser(user_model.UserModel{
		ID:      uuid.New(),
		Token:   token,
		Blocked: false,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return "", err
	}

	err = u.FileService.CreateFolder(fldName)

	if err != nil {
		log.Error().Err(err).Msg("failed to create folder")
		return "", err
	}

	return usrModel.Token, nil
}
