package folder_service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
	"github.com/ivan/storage-project-back/internal/repository/folder_repo"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/rs/zerolog/log"
)

type FolderService struct {
	StoragePath string
	FldRepo     *folder_repo.FldRepo
}

func NewFldService(cfg *config.Config, fldRepo *folder_repo.FldRepo) *FolderService {
	return &FolderService{
		StoragePath: cfg.StoragePath,
		FldRepo:     fldRepo,
	}
}

// TODO: Проанализировать работу метода. То что он не возвращает ошибок, странно...
func (f *FolderService) FolderExist(fldName string) bool {
	fldModel, err := f.FldRepo.GetGeneralFolderByName(fldName)

	if err != nil {
		log.Error().Err(err).Msg("failed to get folder")
		return false
	}

	if fldModel == nil {
		return false
	}

	return true
}

func (f *FolderService) CreateFolder(fldName string, usrID uuid.UUID) error {
	fullPath := filepath.Join(f.StoragePath, fldName)

	err := os.MkdirAll(fullPath, 0755) // TODO: Убрать магические числа

	if err != nil {
		return errsvc.FldErr.CreateFailed
	}

	fldId := uuid.New()

	fldModel, err := f.FldRepo.CreateFld(folder_model.FolderModel{
		ID:        fldId,
		Name:      fldName,
		ParentID:  nil,
		OwnerID:   usrID,
		MainFldId: nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		os.RemoveAll(fullPath)
		return errsvc.FldErr.CreateFailed
	}

	_, err = f.FldRepo.InsertFolderAccess(folder_model.FolderAccessModel{
		FolderID: fldModel.ID,
		UserID:   usrID,
		RoleID:   roles_model.Admin,
	})

	if err != nil {
		os.RemoveAll(fullPath)
		return errsvc.FldErr.CreateFailed
	}

	log.Debug().
		Str("path", fullPath).
		Msg(fmt.Sprintf("folder created with id %s", fldModel.ID.String()))

	return nil
}

func (f *FolderService) DelFld(fldName string) error {
	fldModel, err := f.FldRepo.GetFldByName(fldName)

	if err != nil || fldModel == nil {
		return errsvc.FldErr.NotFound
	}

	if fldModel.MainFldId == nil {
		return errsvc.FldErr.CantDelMainFld
	}

	err = f.FldRepo.DelFld(fldModel.ID)

	if err != nil {
		return errsvc.FldErr.DelFailed
	}

	err = os.RemoveAll(filepath.Join(f.StoragePath, fldName))

	if err != nil {
		return errsvc.FldErr.DelFailed
	}

	return nil
}
