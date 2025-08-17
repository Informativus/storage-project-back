package repository

import (
	user_repo "github.com/ivan/storage-project-back/internal/repository/user"
	"github.com/ivan/storage-project-back/pkg/database/database"
)

type Repositories struct {
	UserRepo *user_repo.UserRepo
}

func NewRepositories(conn database.DBClient) *Repositories {
	return &Repositories{
		UserRepo: user_repo.NewUserRepo(conn),
	}
}
