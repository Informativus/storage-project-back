package fld_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/fld_dto"
	"github.com/ivan/storage-project-back/internal/middlewares/fld_middleware"
	"github.com/ivan/storage-project-back/internal/middlewares/guard"
	"github.com/ivan/storage-project-back/internal/models/user_model"
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
// @Param fldID path string true "Folder name to delete"
// @Security BearerAuth
// @Success 204 "No Content"
// @Router /fld/delete/{fldID} [delete]
func (fc *FldController) DelFld(c *gin.Context) {
	dto := c.MustGet(fld_middleware.SetDelFldDtoKey).(fld_dto.DelFld)
	usrDto := c.MustGet(guard.SetUsrDtoKey).(*user_model.UserModel)

	err := fc.fldService.DelFld(dto.FldID, usrDto.ID)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

// @Summary Create a user subfolder
// @Description Creates a subfolder for main folder
// @Tags Folders
// @Accept json
// @Produce json
// @Param folder body fld_dto.CreateFldReq true "Fld info"
// @Security BearerAuth
// @Success 200 {object} fld_dto.CreateFldRes "Successful response"
// @Router /fld/create [post]
func (fc *FldController) CreateFld(c *gin.Context) {
	dto := c.MustGet(fld_middleware.SetCreateFldDtoKey).(fld_dto.CreateFldReq)
	usrDto := c.MustGet(guard.SetUsrDtoKey).(*user_model.UserModel)

	fldID, err := fc.fldService.CreateSubFld(dto.Name, dto.ParentID, usrDto)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"fldID": fldID,
	})
}
