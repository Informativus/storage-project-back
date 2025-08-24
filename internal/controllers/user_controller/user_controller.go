package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/user_dto"
	"github.com/ivan/storage-project-back/internal/models/user_model"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/internal/services/user_service"
	"github.com/ivan/storage-project-back/pkg/errsvc"
)

type UserController struct {
	UserService  *user_service.UserService
	ErrorService *errsvc.ErrorService
}

func NewUserController(services *services.Services, err *errsvc.ErrorService) *UserController {
	return &UserController{
		UserService:  services.UserService,
		ErrorService: err,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	dto := c.MustGet("userDTO").(user_dto.CreateUserDto)

	connUsrToFld := false

	if dto.ConnUserToFld != nil {
		connUsrToFld = *dto.ConnUserToFld
	}

	token, err := uc.UserService.CreateUser(dto.UrsName, connUsrToFld)

	if err != nil {
		httpErr := uc.ErrorService.MapError(err)
		c.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (uc *UserController) DltUser(c *gin.Context) {
	usrDTO := c.MustGet("usrDTO").(*user_model.UserModel)
	err := uc.UserService.DelUser(usrDTO.ID)
	if err != nil {

		httpErr := uc.ErrorService.MapError(err)
		c.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
