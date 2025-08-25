package file_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/internal/services/file_service"
	"github.com/rs/zerolog/log"
)

type FileController struct {
	FileService *file_service.FileService
}

func NewFileController(services *services.Services) *FileController {
	return &FileController{FileService: services.FileService}
}

func (f *FileController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dst, err := f.FileService.PrepareStorage(file)

	if err != nil {
		c.Error(err)
		return
	}

	log.Info().Msg(dst)
	err = c.SaveUploadedFile(file, dst)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "upload file"})
}
