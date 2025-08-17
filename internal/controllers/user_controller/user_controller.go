package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/dtos/user_dto"
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

	token, err := uc.UserService.CreateUser(dto.FldName, dto.ConnUserToFld)

	if err != nil {
		httpErr := uc.ErrorService.MapError(err)
		c.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
