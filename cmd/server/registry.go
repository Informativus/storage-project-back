package main

import (
	"gitlab.com/ivan/storage-project-back/internal/controllers"
	"gitlab.com/ivan/storage-project-back/internal/services"
	"gitlab.com/ivan/storage-project-back/pkg/config"
	"gitlab.com/ivan/storage-project-back/pkg/errsvc"
	"gitlab.com/ivan/storage-project-back/pkg/jwt_service"
)

type Registry struct {
	Controllers *controllers.Controllers
	Services    *services.Services
}

func NewRegistry(cfg *config.Config) *Registry {
	errsvc := errsvc.NewErrorService()
	jwt := jwt_service.NewJwtService(cfg)
	services := services.NewServices(cfg, jwt)
	controllers := controllers.NewControllers(services, errsvc)

	return &Registry{
		Controllers: controllers,
		Services:    services,
	}
}
