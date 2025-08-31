package file_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/file_dto"
	"github.com/ivan/storage-project-back/internal/utils/validation"
	"github.com/rs/zerolog/log"
)

const SetUploadFileKey = "uploadFileDTO"

func UploadFileMidd(c *gin.Context) {
	var dto file_dto.UploadFileDto

	if err := c.ShouldBindWith(&dto, binding.FormMultipart); err != nil {
		log.Error().Err(err).Msg("failed to bind upload file dto")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request should be multipart/form-data"})
		c.Abort()
		return
	}

	fldUUID, err := uuid.Parse(dto.FldIDStr)
	if err != nil {
		log.Error().Err(err).Msg("invalid folderID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "folderID is not valid UUID"})
		c.Abort()
		return
	}

	dto.FldID = fldUUID

	if err := validation.Validate.Struct(dto); err != nil {
		log.Error().Err(err).Msg("failed to validate upload file dto")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate request"})
		c.Abort()
		return
	}

	c.Set(SetUploadFileKey, &dto)

	c.Next()
}
