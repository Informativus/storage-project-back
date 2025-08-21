package folder_service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/repository/folder_repo"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/rs/zerolog/log"
)

type FolderService struct {
	StoragePath string
	FolderRepo  *folder_repo.FldRepo
}

func NewFldService(cfg *config.Config, fldRepo *folder_repo.FldRepo) *FolderService {
	return &FolderService{
		StoragePath: cfg.StoragePath,
		FolderRepo:  fldRepo,
	}
}

func (f *FolderService) FolderExist(fldName string) bool {
	fldModel, err := f.FolderRepo.GetGeneralFolderByName(fldName)

	if err != nil {
		log.Error().Err(err).Msg("failed to get folder")
		return false
	}

	if fldModel == (folder_model.FolderModel{}) {
		return false
	}

	return true
}

func (f *FolderService) CreateFolder(folderName string, usrID uuid.UUID) error {
	fullPath := filepath.Join(f.StoragePath, folderName)

	err := os.MkdirAll(fullPath, 0755) // TODO: Убрать магические числа

	if err != nil {
		log.Error().Err(err).Msg("failed to create folder")
		return errsvc.ErrGenFolderFailed
	}

	fldModel, err := f.FolderRepo.CreateFld(folder_model.FolderModel{
		ID:         uuid.New(),
		Name:       folderName,
		UserID:     usrID,
		ParentID:   nil,
		Path:       folderName,
		CreatedAt:  time.Now(),
		LastUpdate: time.Now(),
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to create folder model")
		return err
	}

	log.Debug().
		Str("path", fullPath).
		Msg(fmt.Sprintf("folder created with id %s", fldModel.ID.String()))

	return nil
}

func (f *FolderService) DelMainFld(fldName string) error {
	fldModel, err := f.FolderRepo.GetGeneralFolderByName(fldName)

	if err != nil {
		log.Error().Err(err).Msg("failed to get folder")
		return errsvc.ErrFldNotFound
	}

	err = os.RemoveAll(filepath.Join(f.StoragePath, fldName))

	if err != nil {
		log.Error().Err(err).Msg("failed to delete folder")
		return errsvc.ErrFldDeleteFailed
	}

	err = f.FolderRepo.DelMainFld(fldModel.ID)

	if err != nil {
		log.Error().Err(err).Msg("failed to delete folder")
		return errsvc.ErrFldDeleteFailed
	}

	return nil
}
