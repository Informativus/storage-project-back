package fld_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/fld_dto"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
	"github.com/ivan/storage-project-back/pkg/errsvc"
)

type FldController struct {
	fldService *folder_service.FolderService
	err        *errsvc.ErrorService
}

func NewFldController(services *services.Services, err *errsvc.ErrorService) *FldController {
	return &FldController{
		fldService: services.FolderService,
		err:        err,
	}
}

func (fc *FldController) DeleteFld(c *gin.Context) {
	dto := c.MustGet("deleteDTO").(fld_dto.DeleteFld)

	err := fc.fldService.DeleteMainFld(dto.Name)

	if err != nil {
		httpErr := fc.err.MapError(err)
		c.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder deleted"})
}
