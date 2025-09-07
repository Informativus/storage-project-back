package file_service

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/file_dto"
	"github.com/ivan/storage-project-back/internal/models/file_model"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/repository/file_repo"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
	"github.com/ivan/storage-project-back/internal/services/security_service"
	"github.com/ivan/storage-project-back/internal/utils/encryption"
	"github.com/ivan/storage-project-back/internal/utils/storage_util"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/errsvc"
)

type FileService struct {
	Cfg             *config.Config
	SecurityService *security_service.SecurityService
	FolderService   *folder_service.FolderService
	FileRepo        *file_repo.FileRepo
}

func NewFileService(cfg *config.Config, securityService *security_service.SecurityService, folderService *folder_service.FolderService, fileRepo *file_repo.FileRepo) *FileService {
	return &FileService{Cfg: cfg, SecurityService: securityService, FolderService: folderService, FileRepo: fileRepo}
}

func (f *FileService) UploadFileToFld(usrModel *user_model.UserModel, uploadFileDto *file_dto.UploadFileDto) (*uuid.UUID, error) {
	err := f.SecurityService.AccessToUploadFileInFld(usrModel, uploadFileDto.FldID)

	if err != nil {
		return nil, err
	}

	fldModel, err := f.FolderService.GetGeneralFolderByFldId(uploadFileDto.FldID)

	if err != nil {
		return nil, err
	}

	fileEncrypted, err := encryption.IsPGPEncrypted(uploadFileDto.File)

	if err != nil {
		return nil, errsvc.FileErr.Internal.New(err)
	}

	ext := ".gpg"
	uuid := uuid.New()
	storageKey := f.generateStorageKey(uuid, ext)

	pathToSave := f.prepareStorage(fldModel.Name, storageKey)

	var dataToSave []byte
	if fileEncrypted {
		fHeaderFile, err := uploadFileDto.File.Open()

		if err != nil {
			return nil, errsvc.FileErr.Internal.New(err)
		}

		defer fHeaderFile.Close()

		dataToSave, err = io.ReadAll(fHeaderFile)

		if err != nil {
			return nil, errsvc.FileErr.Internal.New(err)
		}
	} else {
		encryptedFile, err := encryption.EncryptFile(uploadFileDto.File, f.Cfg.GpgPublicKeyPath)

		if err != nil {
			return nil, errsvc.FileErr.Internal.New(err)
		}

		dataToSave = encryptedFile
	}

	err = storage_util.SaveFileInStorage(dataToSave, pathToSave)

	if err != nil {
		return nil, errsvc.FileErr.SaveFailed.New(err)
	}

	fileModel := &file_model.FileModel{
		ID:         uuid,
		Name:       uploadFileDto.File.Filename,
		FolderID:   uploadFileDto.FldID,
		Size:       uploadFileDto.File.Size,
		MimeType:   uploadFileDto.File.Header.Get("Content-Type"),
		StorageKey: storageKey,
		CreatedAt:  time.Now(),
		DeletedAt:  nil,
	}

	fileModel, err = f.FileRepo.UploadFile(fileModel)

	if err != nil {
		if err := storage_util.DelFileFromStorage(pathToSave); err != nil {
			return nil, errsvc.FileErr.InconsistentState.New(err)
		}
		return nil, errsvc.FileErr.DelFailed.New(err)
	}

	return &fileModel.ID, nil
}

func (f *FileService) DelFile(usrModel *user_model.UserModel, fileID uuid.UUID) error {
	// TODO: Добавить проверку на права удаления
	fileModel, err := f.FileRepo.GetFileByID(fileID)

	if err != nil {
		return errsvc.FileErr.Internal.New(err)
	}

	if fileModel == nil {
		return errsvc.FileErr.NotFound.New(err)
	}

	markedCount, err := f.FileRepo.MarkFileAsDeleted(fileID, time.Now())

	if err != nil || markedCount == 0 {
		return errsvc.FileErr.Internal.New(err)
	}

	return nil
}

func (f *FileService) prepareStorage(fldPathToSave string, storageFileName string) string {
	dst := filepath.Join(f.Cfg.StoragePath, fldPathToSave, storageFileName)

	return dst
}

func (f *FileService) generateStorageKey(uuid uuid.UUID, ext string) string {
	return fmt.Sprintf("%s%s", uuid, ext)

}
