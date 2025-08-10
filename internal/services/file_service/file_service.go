package file_service

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"

	"gitlab.com/ivan/storage-project-back/pkg/config"
)

type FileService struct {
	StoragePath string
}

func NewFileService(cfg *config.Config) *FileService {
	return &FileService{StoragePath: cfg.StoragePath}
}

func (f *FileService) PrepareStorage(file *multipart.FileHeader) (string, error) {
	err := os.MkdirAll(f.StoragePath, os.ModePerm)

	if err != nil {
		return "", errors.New("gen_folder_failed")
	}

	dst := filepath.Join(f.StoragePath, file.Filename)

	return dst, nil
}
