package services

import (
	fileservice "gitlab.com/ivan/storage-project-back/internal/services/file_service"
	userservice "gitlab.com/ivan/storage-project-back/internal/services/user_service"
	"gitlab.com/ivan/storage-project-back/pkg/config"
	"gitlab.com/ivan/storage-project-back/pkg/jwt_service"
)

type Services struct {
	UserService *userservice.UserService
	FileService *fileservice.FileService
}

func NewServices(cfg *config.Config, jwt *jwt_service.JwtService) *Services {
	filesrvc := fileservice.NewFileService(cfg)

	return &Services{
		FileService: filesrvc,
		UserService: userservice.NewUserService(cfg, jwt, filesrvc),
	}
}
