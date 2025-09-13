package jobs

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ivan/storage-project-back/internal/repository"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/rs/zerolog/log"
)

type StartJobs struct {
	repo *repository.Repositories
	cfg  *config.Config
}

func NewStartJobs(repo *repository.Repositories, cfg *config.Config) *StartJobs {
	return &StartJobs{
		repo: repo,
		cfg:  cfg,
	}
}

func (j *StartJobs) StartAllJobs() {
	j.StartTokenCleanupJob()
	j.StartFileCleanUpJob()
	j.StartCleanUpUsrs()
}

func (j *StartJobs) StartTokenCleanupJob() {
	ticker := time.NewTicker(time.Hour * 24)

	go func() {
		for range ticker.C {
			if _, err := j.repo.UserRepo.DelExpiredTokens(); err != nil {
				log.Error().Err(err).Msg("failed to delete expired tokens")
			}
		}
	}()
}

func (j *StartJobs) StartFileCleanUpJob() {
	ticker := time.NewTicker(time.Hour)

	go func() {
		for range ticker.C {
			files, err := j.repo.FileRepo.GetMarkedToDelFiles()

			if err != nil {
				log.Error().Err(err).Msg("failed to get marked to del files")
				continue
			}

			for _, file := range files {
				fldModel, err := j.repo.FldRepo.GetGeneralFolderBySubFldId(file.FolderID)
				if err != nil {
					log.Error().Err(err).Msg("failed to get fld model")
					continue
				}

				if fldModel == nil {
					fldModel, err = j.repo.FldRepo.GetGeneralFolderById(file.FolderID)
					if err != nil {
						log.Error().Err(err).Msg("failed to get fld model")
						continue
					}
				}

				storagePath := filepath.Join(j.cfg.StoragePath, fldModel.Name, file.StorageKey)
				if err := os.Remove(storagePath); err != nil {
					log.Error().Err(err).Msg("failed to delete file from disk")
					continue
				}

				if tag, err := j.repo.FileRepo.HardDelFile(file.ID); err != nil || tag == 0 {
					log.Error().Err(err).Msg("failed to delete file from db")
					continue
				}
			}
		}
	}()
}

func (j *StartJobs) StartCleanUpUsrs() {
	ticker := time.NewTicker(time.Hour)

	go func() {
		for range ticker.C {
			delUsrs, err := j.repo.UserRepo.GetMarkedToDelUsrs()

			log.Info().Msgf("found %d users to delete", len(delUsrs))
			if err != nil {
				log.Error().Err(err).Msg("failed to get marked to del users")
				continue
			}

			for _, usr := range delUsrs {
				if tag, err := j.repo.UserRepo.HardDel(usr.ID); err != nil || tag == 0 {
					log.Error().Err(err).Msg("failed to delete user from db")
					continue
				}
			}
		}
	}()
}
