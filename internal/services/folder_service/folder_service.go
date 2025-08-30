package folder_service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/models/folder_model"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/repository/folder_repo"
	"github.com/ivan/storage-project-back/internal/services/security_service"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/rs/zerolog/log"
)

const (
	folderPermissionRWXRXRX os.FileMode = 0755
)

type FolderService struct {
	StoragePath string
	FldRepo     *folder_repo.FldRepo
	Security    *security_service.SecurityService
}

func NewFldService(cfg *config.Config, fldRepo *folder_repo.FldRepo, securityService *security_service.SecurityService) *FolderService {
	return &FolderService{
		StoragePath: cfg.StoragePath,
		FldRepo:     fldRepo,
		Security:    securityService,
	}
}

func (f *FolderService) MainFolderExist(fldName string) (bool, error) {
	fldModel, err := f.FldRepo.GetGeneralFolderByName(fldName)

	if err != nil {
		log.Error().Err(err).Msg("failed to get folder")
		return false, err
	}

	if fldModel == nil {
		return false, nil
	}

	return true, nil
}

func (f *FolderService) CreateFolder(fldName string, usrModel *user_model.UserModel) error {
	usrID := usrModel.ID
	fullPath := filepath.Join(f.StoragePath, fldName)

	err := os.MkdirAll(fullPath, folderPermissionRWXRXRX)

	if err != nil {
		return errsvc.FldErr.CreateFailed.New()
	}

	fldId := uuid.New()

	fldModel, err := f.FldRepo.CreateFld(&folder_model.FolderModel{
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
		return errsvc.FldErr.CreateFailed.New()
	}

	_, err = f.FldRepo.InsertFolderAccess(folder_model.FolderAccessModel{
		FolderID: fldModel.ID,
		UserID:   usrID,
		RoleID:   f.Security.GetUsrRoleForFld(usrModel, fldModel),
	})

	if err != nil {
		os.RemoveAll(fullPath)
		return errsvc.FldErr.CreateFailed.New()
	}

	log.Debug().
		Str("path", fullPath).
		Msg(fmt.Sprintf("folder created with id %s", fldModel.ID.String()))

	return nil
}

func (f *FolderService) DelFld(fldName string, usrID uuid.UUID) error {
	mainFldModel, err := f.FldRepo.GetGeneralFolderByUsrId(usrID)

	if err != nil || mainFldModel == nil {
		return errsvc.FldErr.NotFound.New()
	}

	fldModel, err := f.FldRepo.GetFldByNameAndMainFldId(fldName, mainFldModel.FolderID)

	if err != nil || fldModel == nil {
		return errsvc.FldErr.NotFound.New()
	}

	if fldModel.MainFldId == nil {
		return errsvc.FldErr.CantDelMainFld.New()
	}

	err = f.FldRepo.DelFld(fldModel.ID)

	if err != nil {
		return errsvc.FldErr.DelFailed.New()
	}

	err = os.RemoveAll(filepath.Join(f.StoragePath, fldName))

	if err != nil {
		return errsvc.FldErr.DelFailed.New()
	}

	return nil
}

func (f *FolderService) CreateSubFld(fldName string, parentID uuid.UUID, usrModel *user_model.UserModel) (*uuid.UUID, error) {
	err := f.Security.AccessToCreateFldForUsr(usrModel, parentID)

	if err != nil {
		return nil, err
	}

	isMainFld, mainFldModel, err := f.isMainFld(parentID)

	if err != nil {
		return nil, err
	}

	if !isMainFld {
		mainFldModel, err = f.FldRepo.GetGeneralFolderBySubFldId(parentID)
	}

	if err != nil {
		return nil, errsvc.FldErr.NotFound.New()
	}

	if mainFldModel.Name == fldName {
		return nil, errsvc.FldErr.AlreadyExists.New()
	}

	fldModel, err := f.FldRepo.GetFldByNameAndMainFldId(fldName, mainFldModel.ID)

	if err != nil {
		return nil, errsvc.FldErr.Internal.New()
	}

	if fldModel != nil {
		return nil, errsvc.FldErr.AlreadyExists.New()
	}

	fldModel = &folder_model.FolderModel{
		ID:        uuid.New(),
		Name:      fldName,
		ParentID:  &parentID,
		OwnerID:   usrModel.ID,
		MainFldId: &mainFldModel.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	fldModel, err = f.FldRepo.CreateFld(fldModel)

	if err != nil {
		return nil, errsvc.FldErr.CreateFailed.New()
	}

	_, err = f.FldRepo.InsertFolderAccess(folder_model.FolderAccessModel{
		FolderID: fldModel.ID,
		UserID:   usrModel.ID,
		RoleID:   f.Security.GetUsrRoleForFld(usrModel, mainFldModel),
	})

	if err != nil {
		return nil, f.rollbackFolder(fldName, fldModel.ID, "failed to insert folder access", err)
	}

	return &fldModel.ID, nil
}

func (f *FolderService) isMainFld(fldID uuid.UUID) (bool, *folder_model.FolderModel, error) {
	fldModel, err := f.FldRepo.GetGeneralFolderById(fldID)

	if err != nil {
		log.Error().Err(err).Msg("failed to get folder")
		return false, fldModel, errsvc.FldErr.Internal.New()
	}

	if fldModel != nil {
		return true, fldModel, nil
	}

	return false, fldModel, nil
}

func (f *FolderService) rollbackFolder(fldName string, id uuid.UUID, logMsg string, err error) error {
	log.Error().Err(err).Msg(logMsg)

	if delErr := f.DelFld(fldName, id); delErr != nil {
		log.Error().Err(delErr).Str("folder_id", id.String()).
			Msg("rollback failed, inconsistent state: manual cleanup required")
		return errsvc.UsrErr.InconsistentState.New()
	}
	return err
}
