package file_middleware

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
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

	validatePublicKey(dto.PublicKey, c)

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

func validatePublicKey(pubKey *multipart.FileHeader, c *gin.Context) {
	if pubKey == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "publicKey is required"})
		c.Abort()
		return
	}

	f, err := pubKey.Open()

	if err != nil {
		log.Error().Err(err).Msg("failed to open public key file")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open public key file"})
		c.Abort()
		return
	}

	defer f.Close()

	data, err := io.ReadAll(f)

	if err != nil {
		log.Error().Err(err).Msg("failed to read public key file")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read public key file"})
		c.Abort()
		return
	}

	if _, err := crypto.NewKeyFromArmored(string(data)); err != nil {
		log.Error().Err(err).Msg("failed to parse public key")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse public key"})
		c.Abort()
		return
	}
}
