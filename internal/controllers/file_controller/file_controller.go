package file_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/file_dto"
	"github.com/ivan/storage-project-back/internal/middlewares/file_middleware"
	"github.com/ivan/storage-project-back/internal/middlewares/guard"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/internal/services/file_service"
)

type FileController struct {
	FileService *file_service.FileService
}

func NewFileController(services *services.Services) *FileController {
	return &FileController{FileService: services.FileService}
}

// @Summary Upload file to server
// @Description Upload user file to server with folder id.
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param name formData string true "Name of the file"
// @Param folderID formData string true "Folder UUID"
// @Param publicKey formData file true "Public key for encryption" example(public.key)
// @Security BearerAuth
// @Success 200 {object} file_dto.UploadFileDtoRes "Successful response"
// @Router /file/upload [post]
func (f *FileController) Upload(c *gin.Context) {
	usrDto := c.MustGet(guard.SetUsrDtoKey).(*user_model.UserModel)
	uploadFileDto := c.MustGet(file_middleware.SetUploadFileKey).(*file_dto.UploadFileDto)

	fileID, err := f.FileService.UploadFileToFld(usrDto, uploadFileDto)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"fileID": fileID,
	})
}

// @Summary Delete a file
// @Description Soft delete a file by its UUID. File will be marked as deleted and removed asynchronously by cleanup job.
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param fileID path string true "File ID to delete"
// @Security BearerAuth
// @Success 204 "No Content"
// @Router /file/delete/{fileID} [delete]
func (f *FileController) Del(c *gin.Context) {
	// usrDto := c.MustGet(guard.SetUsrDtoKey).(*user_model.UserModel)
	delFileDto := c.MustGet(file_middleware.SetDelFileKey).(*file_dto.DelFileDto)

	err := f.FileService.DelFile(delFileDto.FileID)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
