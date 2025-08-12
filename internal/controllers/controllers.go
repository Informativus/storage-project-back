package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/file_controller"
	"github.com/ivan/storage-project-back/internal/controllers/user_controller"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/pkg/errsvc"
)

type Controllers struct {
	FileController *file_controller.FileController
	UserController *user_controller.UserController
}

func NewControllers(services *services.Services, err *errsvc.ErrorService) *Controllers {
	return &Controllers{
		FileController: file_controller.NewFileController(services, err),
		UserController: user_controller.NewUserController(services, err),
	}
}

func (c *Controllers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		file := api.Group("/file")
		{
			file.POST("/upload", c.FileController.Upload)
		}

		user := api.Group("/user")
		{
			user.GET("/token", c.UserController.GenerateTokenForUser)
			user.POST("/create", c.UserController.CreateUser)
		}
	}
}
