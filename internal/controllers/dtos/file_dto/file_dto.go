package file_dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type UploadFileDto struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Name      string                `form:"name" validate:"required,file_name_min,file_name_max,file_name_valid"`
	FldIDStr  string                `form:"folderID" binding:"required"`
	PublicKey *multipart.FileHeader `form:"publicKey" validate:"required"`
	FldID     uuid.UUID             `json:"-"`
}

type UploadFileDtoRes struct {
	FileId uuid.UUID `json:"fileID"`
}
