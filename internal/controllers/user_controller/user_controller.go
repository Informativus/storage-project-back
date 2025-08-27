package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/user_dto"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/internal/services/user_service"
)

type UserController struct {
	UserService *user_service.UserService
}

func NewUserController(services *services.Services) *UserController {
	return &UserController{
		UserService: services.UserService,
	}
}

// @Summary Create a new user
// @Description Creates a user with the folder name and connects the user to the folder
// @Tags User
// @Accept json
// @Produce json
// @Param user body user_dto.CreateUserDto true "User info"
// @Security BearerAuth
// @Success 200 {object} user_dto.CreateUserResponse "Successful response"
// @Router /user/create [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	dto := c.MustGet("userDTO").(user_dto.CreateUserDto)

	connUsrToFld := false

	if dto.ConnUserToFld != nil {
		connUsrToFld = *dto.ConnUserToFld
	}

	token, err := uc.UserService.CreateUser(dto.UrsName, connUsrToFld)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Delete a user with all his data
// @Description Deletes a user with all his data (folders, files)
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 204 "No Content"
// @Router /user/delete [delete]
func (uc *UserController) DltUser(c *gin.Context) {
	usrDTO := c.MustGet("usrDTO").(*user_model.UserModel)
	err := uc.UserService.DelUser(usrDTO.ID)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
