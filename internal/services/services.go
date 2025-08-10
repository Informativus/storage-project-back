package services

import (
	fileservice "gitlab.com/ivan/storage-project-back/internal/services/file_service"
	userservice "gitlab.com/ivan/storage-project-back/internal/services/user_service"
	"gitlab.com/ivan/storage-project-back/pkg/config"
	"gitlab.com/ivan/storage-project-back/pkg/jwt_service"
)

type Services struct {
	FileService *fileservice.FileService
	UserService *userservice.UserService
	JwtService  *jwt_service.JwtService
}

func NewServices(cfg *config.Config, jwt *jwt_service.JwtService) *Services {
	return &Services{
		FileService: fileservice.NewFileService(cfg),
		UserService: userservice.NewUserService(cfg, jwt),
	}
}
