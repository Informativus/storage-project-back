package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/file_controller"
	"github.com/ivan/storage-project-back/internal/controllers/fld_controller"
	"github.com/ivan/storage-project-back/internal/controllers/user_controller"
	"github.com/ivan/storage-project-back/internal/middlewares/fld_middleware"
	"github.com/ivan/storage-project-back/internal/middlewares/users_middleware"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/pkg/errsvc"
)

type Controllers struct {
	FileController *file_controller.FileController
	FldController  *fld_controller.FldController
	UserController *user_controller.UserController
}

func NewControllers(services *services.Services, err *errsvc.ErrorService) *Controllers {
	return &Controllers{
		FileController: file_controller.NewFileController(services, err),
		FldController:  fld_controller.NewFldController(services, err),
		UserController: user_controller.NewUserController(services, err),
	}
}

func (c *Controllers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/create", users_middleware.CreateUserMidd, c.UserController.CreateUser)
		}

		fld := api.Group("/fld")
		{
			fld.DELETE("/delete/:fldName", fld_middleware.DelFld, c.FldController.DelFld)
		}

		file := api.Group("/file")
		{
			file.POST("/upload", c.FileController.Upload)
		}
	}
}
