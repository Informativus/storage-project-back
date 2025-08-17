package repository

import (
	"github.com/ivan/storage-project-back/internal/repository/folder_repo"
	user_repo "github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/pkg/database/database"
)

type Repositories struct {
	UserRepo *user_repo.UserRepo
	FldRepo  *folder_repo.FldRepo
}

func NewRepositories(db database.DBClient) *Repositories {
	return &Repositories{
		UserRepo: user_repo.NewUserRepo(db),
		FldRepo:  folder_repo.NewFldRepo(db),
	}
}
