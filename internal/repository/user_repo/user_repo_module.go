package user_repo

import (
	db_user "github.com/ivan/storage-project-back/internal/repository/user_repo/db"
	"github.com/ivan/storage-project-back/pkg/database"
)

type UserRepoModule struct {
	UserRepo *UserRepo
}

func NewUserModule(db *database.DatabaseModule) *UserRepoModule {
	sqlRepo := db_user.NewSqlUserRepo(db)
	cacheRepo := db_user.NewCacheUserRepo(db.NoSQLDB)

	return &UserRepoModule{
		UserRepo: NewUserRepo(cacheRepo, sqlRepo),
	}
}
