package repository

import (
	"github.com/ivan/storage-project-back/internal/repository/file_repo"
	"github.com/ivan/storage-project-back/internal/repository/folder_repo"
	"github.com/ivan/storage-project-back/internal/repository/security_repo"
	"github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/pkg/database"
)

type Repositories struct {
	UserRepo *user_repo.UserRepo
	FldRepo  *folder_repo.FldRepo
	SecRepo  *security_repo.SecurityRepo
	FileRepo *file_repo.FileRepo
}

func NewRepositories(db *database.DatabaseModule) *Repositories {
	usrModule := user_repo.NewUserModule(db)

	return &Repositories{
		UserRepo: usrModule.UserRepo,
		FldRepo:  folder_repo.NewFldRepo(db),
		SecRepo:  security_repo.NewSecurityRepo(db),
		FileRepo: file_repo.NewFileRepo(db),
	}
}
