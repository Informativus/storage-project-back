package file_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gitlab.com/ivan/storage-project-back/internal/services"
	"gitlab.com/ivan/storage-project-back/internal/services/file_service"
	"gitlab.com/ivan/storage-project-back/pkg/errsvc"
)

type FileController struct {
	FileService  *file_service.FileService
	ErrorService *errsvc.ErrorService
}

func NewFileController(services *services.Services, err *errsvc.ErrorService) *FileController {
	return &FileController{FileService: services.FileService, ErrorService: err}
}

func (f *FileController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dst, err := f.FileService.PrepareStorage(file)

	if err != nil {
		httpErr := f.ErrorService.MapError(err)
		c.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
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
