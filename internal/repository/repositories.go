package repository

import (
	"github.com/ivan/storage-project-back/internal/repository/folder_repo"
	"github.com/ivan/storage-project-back/internal/repository/security_repo"
	user_repo "github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/pkg/database/database"
)

type Repositories struct {
	UserRepo *user_repo.UserRepo
	FldRepo  *folder_repo.FldRepo
	SecRepo  *security_repo.SecurityRepo
}

func NewRepositories(db database.DBClient) *Repositories {
	return &Repositories{
		UserRepo: user_repo.NewUserRepo(db),
		FldRepo:  folder_repo.NewFldRepo(db),
		SecRepo:  security_repo.NewSecurityRepo(db),
	}
}
