package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ivan/storage-project-back/docs"
	"github.com/ivan/storage-project-back/internal/controllers"
	"github.com/ivan/storage-project-back/internal/repository"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/database"
	"github.com/ivan/storage-project-back/pkg/jobs"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Registry struct {
	Controllers *controllers.Controllers
	Services    *services.Services
}

func NewRegistry(cfg *config.Config, routers *gin.Engine) *Registry {
	conn, err := database.NewDatabaseModule(cfg)
	routers.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	jwt := jwt_service.NewJwtService(cfg)
	repos := repository.NewRepositories(conn)
	services := services.NewServices(cfg, repos, jwt)
	controllers := controllers.NewControllers(services, jwt)

	jbs := jobs.NewStartJobs(repos, cfg)

	jbs.StartAllJobs()

	return &Registry{
		Controllers: controllers,
		Services:    services,
	}
}
