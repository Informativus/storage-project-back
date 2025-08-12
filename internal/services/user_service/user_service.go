package user_service

import (
	"errors"

	"github.com/ivan/storage-project-back/internal/services/file_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
)

type UserService struct {
	cfg         *config.Config
	jwt         *jwt_service.JwtService
	FileService *file_service.FileService
}

func NewUserService(cfg *config.Config, jwt *jwt_service.JwtService, filesrvc *file_service.FileService) *UserService {
	return &UserService{cfg: cfg, jwt: jwt, FileService: filesrvc}
}

func (u *UserService) GenerateToken(folderName string) (string, error) {
	token, err := u.jwt.GenerateToken(jwt_service.JwtPayload{FolderName: folderName})

	if err != nil {
		log.Error().Err(err).Msg("failed to sign token")
		return "", errors.New("internal_error")
	}

	return token, nil
}

func (u *UserService) CreateUser(folderName string) error {
	folderExt := u.FileService.FolderExist(folderName)

	if folderExt {
		return errors.New("folder_exist")
	}

	err := u.FileService.CreateFolder(folderName)

	if err != nil {
		log.Error().Err(err).Msg("failed to create folder")
		return err
	}

	return nil
}
