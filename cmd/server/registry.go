package main

import (
	"github.com/ivan/storage-project-back/internal/controllers"
	"github.com/ivan/storage-project-back/internal/repository"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/database/database"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
	"github.com/rs/zerolog/log"
)

type Registry struct {
	Controllers *controllers.Controllers
	Services    *services.Services
}

func NewRegistry(cfg *config.Config) *Registry {
	conn, err := database.ConnectPg(cfg)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	errsvc := errsvc.NewErrorService()
	jwt := jwt_service.NewJwtService(cfg)
	repos := repository.NewRepositories(conn)
	services := services.NewServices(cfg, repos, jwt)
	controllers := controllers.NewControllers(services, errsvc, jwt)

	return &Registry{
		Controllers: controllers,
		Services:    services,
	}
}
