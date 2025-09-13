package services

import (
	"github.com/ivan/storage-project-back/internal/repository"
	"github.com/ivan/storage-project-back/internal/services/file_service"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
	"github.com/ivan/storage-project-back/internal/services/security_service"
	userservice "github.com/ivan/storage-project-back/internal/services/user_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
)

type Services struct {
	UserService     *userservice.UserService
	FolderService   *folder_service.FolderService
	FileService     *file_service.FileService
	SecurityService *security_service.SecurityService
}

func NewServices(
	cfg *config.Config,
	repos *repository.Repositories,
	jwt *jwt_service.JwtService,
) *Services {
	securityService := security_service.NewSecurityService(repos.SecRepo, repos.FldRepo)
	fldService := folder_service.NewFldService(cfg, repos.FldRepo, securityService)
	fileService := file_service.NewFileService(cfg, securityService, fldService, repos.FileRepo)

	return &Services{
		FileService:     fileService,
		FolderService:   fldService,
		UserService:     userservice.NewUserService(cfg, jwt, repos.UserRepo, fldService, securityService),
		SecurityService: securityService,
	}
}
