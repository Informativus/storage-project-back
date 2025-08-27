package fld_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/fld_dto"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/internal/services/folder_service"
)

type FldController struct {
	fldService *folder_service.FolderService
}

func NewFldController(services *services.Services) *FldController {
	return &FldController{
		fldService: services.FolderService,
	}
}

// @Summary Delete a user folder with all his data
// @Description Deletes a folder with all its data (subfolders, files)
// @Tags Folders
// @Accept json
// @Produce json
// @Param fldName path string true "Folder name to delete"
// @Security BearerAuth
// @Success 204 "No Content"
// @Router /fld/delete/{fldName} [delete]
func (fc *FldController) DelFld(c *gin.Context) {
	dto := c.MustGet("dltDTO").(fld_dto.DelFld)

	err := fc.fldService.DelFld(dto.Name)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
