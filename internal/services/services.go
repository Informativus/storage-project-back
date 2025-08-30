package services

import (
	"github.com/ivan/storage-project-back/internal/repository"
	fileservice "github.com/ivan/storage-project-back/internal/services/file_service"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
	"github.com/ivan/storage-project-back/internal/services/security_service"
	userservice "github.com/ivan/storage-project-back/internal/services/user_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
)

type Services struct {
	UserService     *userservice.UserService
	FileService     *fileservice.FileService
	FolderService   *folder_service.FolderService
	SecurityService *security_service.SecurityService
}

func NewServices(
	cfg *config.Config,
	repos *repository.Repositories,
	jwt *jwt_service.JwtService,
) *Services {
	securityService := security_service.NewSecurityService(repos.SecRepo)
	fldsrvc := folder_service.NewFldService(cfg, repos.FldRepo, securityService)

	return &Services{
		FileService:   fileservice.NewFileService(cfg),
		FolderService: fldsrvc,
		UserService:   userservice.NewUserService(cfg, jwt, repos.UserRepo, fldsrvc),
	}
}
