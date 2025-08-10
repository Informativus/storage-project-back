package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/ivan/storage-project-back/internal/services"
	"gitlab.com/ivan/storage-project-back/internal/services/user_service"
	"gitlab.com/ivan/storage-project-back/pkg/errsvc"
)

type UserController struct {
	UserService  *user_service.UserService
	ErrorService *errsvc.ErrorService
}

func NewUserController(services *services.Services, err *errsvc.ErrorService) *UserController {
	return &UserController{UserService: services.UserService, ErrorService: err}
}

func (uc *UserController) GenerateTokenForUser(c *gin.Context) {
	fldName := c.Query("folderName")

	if fldName == "" {
		httpErr := uc.ErrorService.MapError(errsvc.ErrInvalidFolder)
		c.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
		return
	}

	token, err := uc.UserService.GenerateToken(fldName)
	if err != nil {
		httpErr := uc.ErrorService.MapError(err)
		c.JSON(httpErr.Code, gin.H{"error": httpErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": token})
}
