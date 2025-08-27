package services

import (
	"github.com/ivan/storage-project-back/internal/repository"
	fileservice "github.com/ivan/storage-project-back/internal/services/file_service"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
	userservice "github.com/ivan/storage-project-back/internal/services/user_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
)

type Services struct {
	UserService   *userservice.UserService
	FileService   *fileservice.FileService
	FolderService *folder_service.FolderService
}

func NewServices(
	cfg *config.Config,
	repos *repository.Repositories,
	jwt *jwt_service.JwtService,
) *Services {
	fldsrvc := folder_service.NewFldService(cfg, repos.FldRepo)

	return &Services{
		FileService:   fileservice.NewFileService(cfg),
		FolderService: fldsrvc,
		UserService:   userservice.NewUserService(cfg, jwt, repos.UserRepo, fldsrvc),
	}
}
