package db_user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/pkg/database/no_sql_database"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type CacheUserRepo struct {
	cacheDb *no_sql_database.RedisClient
}

func NewCacheUserRepo(cacheDb *no_sql_database.RedisClient) *CacheUserRepo {
	return &CacheUserRepo{
		cacheDb: cacheDb,
	}
}

func (r *CacheUserRepo) GetUsrById(id uuid.UUID) (*user_model.UserDto, error) {
	row, err := r.cacheDb.Get(context.Background(), id.String())

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	var usrDto user_model.UserDto

	if err := json.Unmarshal([]byte(row), &usrDto); err != nil {
		return nil, err
	}

	return &usrDto, nil
}

func (r *CacheUserRepo) SetUsrById(usrDto user_model.UserDto) error {
	row, err := json.Marshal(usrDto)

	if err != nil {
		return err
	}

	return r.cacheDb.Set(context.Background(), usrDto.ID.String(), row)
}

func (r *CacheUserRepo) DelUsrById(id uuid.UUID) error {
	log.Debug().Str("id", id.String()).Msg("delete user from cache")
	return r.cacheDb.Del(context.Background(), id.String())
}
