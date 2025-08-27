package jobs

import (
	"time"

	"github.com/ivan/storage-project-back/internal/repository"
	"github.com/rs/zerolog/log"
)

type StartJobs struct {
	repo *repository.Repositories
}

func NewStartJobs(repo *repository.Repositories) *StartJobs {
	return &StartJobs{
		repo: repo,
	}
}

func (j *StartJobs) StartAllJobs() {
	j.StartTokenCleanupJob()
}

func (j *StartJobs) StartTokenCleanupJob() {
	ticker := time.NewTicker(time.Hour * 24)

	go func() {
		for range ticker.C {
			if err := j.repo.UserRepo.DelExpiredTokens(); err != nil {
				log.Fatal().Err(err).Msg("failed to delete expired tokens")
			}
		}
	}()
}
