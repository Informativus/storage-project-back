package services

import (
	"github.com/ivan/storage-project-back/internal/repository"
	fileservice "github.com/ivan/storage-project-back/internal/services/file_service"
	userservice "github.com/ivan/storage-project-back/internal/services/user_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
)

type Services struct {
	UserService *userservice.UserService
	FileService *fileservice.FileService
}

func NewServices(
	cfg *config.Config,
	repos *repository.Repositories,
	jwt *jwt_service.JwtService,
) *Services {
	filesrvc := fileservice.NewFileService(cfg)

	return &Services{
		FileService: filesrvc,
		UserService: userservice.NewUserService(cfg, jwt, repos.UserRepo, filesrvc),
	}
}
