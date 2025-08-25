package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/internal/controllers/file_controller"
	"github.com/ivan/storage-project-back/internal/controllers/fld_controller"
	"github.com/ivan/storage-project-back/internal/controllers/user_controller"
	"github.com/ivan/storage-project-back/internal/middlewares/error_middleware"
	"github.com/ivan/storage-project-back/internal/middlewares/fld_middleware"
	"github.com/ivan/storage-project-back/internal/middlewares/guard"
	"github.com/ivan/storage-project-back/internal/middlewares/users_middleware"
	"github.com/ivan/storage-project-back/internal/models/roles_model"
	"github.com/ivan/storage-project-back/internal/repository/user_repo"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
)

type Controllers struct {
	FileController *file_controller.FileController
	FldController  *fld_controller.FldController
	UserController *user_controller.UserController
	UserRepo       *user_repo.UserRepo
	jwt            *jwt_service.JwtService
}

func NewControllers(services *services.Services, jwt *jwt_service.JwtService, userRepo *user_repo.UserRepo) *Controllers {
	return &Controllers{
		FileController: file_controller.NewFileController(services),
		FldController:  fld_controller.NewFldController(services),
		UserController: user_controller.NewUserController(services),
		UserRepo:       userRepo,
		jwt:            jwt,
	}
}

func (c *Controllers) RegisterRoutes(router *gin.Engine) {
	router.Use(error_middleware.ErrorHandler)

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/create", guard.AuthGuard(c.jwt, c.UserRepo, []roles_model.Role{roles_model.Admin}), users_middleware.CreateUserMidd, c.UserController.CreateUser)
			user.DELETE("/delete", guard.AuthGuard(c.jwt, c.UserRepo, []roles_model.Role{roles_model.User, roles_model.Admin}), c.UserController.DltUser)
		}

		fld := api.Group("/fld")
		{
			fld.DELETE("/delete/:fldName", guard.AuthGuard(c.jwt, c.UserRepo, []roles_model.Role{roles_model.User, roles_model.Admin}), fld_middleware.DelFld, c.FldController.DelFld)
		}

		file := api.Group("/file")
		{
			file.POST("/upload", c.FileController.Upload)
		}
	}
}
