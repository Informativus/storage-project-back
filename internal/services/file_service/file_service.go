package file_service

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/rs/zerolog/log"
)

type FileService struct {
	StoragePath string
}

func NewFileService(cfg *config.Config) *FileService {
	return &FileService{StoragePath: cfg.StoragePath}
}

func (f *FileService) PrepareStorage(file *multipart.FileHeader) (string, error) {
	dst := filepath.Join(f.StoragePath, file.Filename)

	return dst, nil
}

func (f *FileService) FolderExist(folderName string) bool {
	fullPath := filepath.Join(f.StoragePath, folderName)

	fileInfo, err := os.Stat(fullPath)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		log.Error().Err(err).Msg("failed to get file info")

		return false
	}

	return fileInfo.IsDir()

}

func (f *FileService) CreateFolder(folderName string) error {
	if folderName == "" || strings.ContainsAny(folderName, `/\:*?"<>|`) {
		return errors.New("invalid_folder_name")
	}

	fullPath := filepath.Join(f.StoragePath, folderName)

	err := os.MkdirAll(fullPath, 0755) // TODO: Убрать магические числа

	if err != nil {
		log.Error().Err(err).Msg("failed to create folder")
		return errors.New("gen_folder_failed")
	}

	log.Debug().Str("path", fullPath).Msg("folder created")

	return nil
}
