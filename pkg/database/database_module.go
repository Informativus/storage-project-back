package database

import (
	"context"
	"errors"

	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/database/no_sql_database"
	"github.com/ivan/storage-project-back/pkg/database/sql_database"
)

type DatabaseModule struct {
	SQLDB   sql_database.DBClient
	NoSQLDB *no_sql_database.RedisClient
}

func NewDatabaseModule(cfg *config.Config) (*DatabaseModule, error) {
	sqlDB, err := sql_database.ConnectPg(cfg)
	if err != nil {
		return nil, err
	}

	redisDB := no_sql_database.NewRedisClient(cfg)
	if redisDB == nil {
		return nil, errors.New("failed to create redis client")
	}

	if err := redisDB.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &DatabaseModule{
		SQLDB:   sqlDB,
		NoSQLDB: redisDB,
	}, nil
}
