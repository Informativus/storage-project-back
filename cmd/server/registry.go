package main

import (
	"github.com/ivan/storage-project-back/internal/controllers"
	"github.com/ivan/storage-project-back/internal/services"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/ivan/storage-project-back/pkg/errsvc"
	"github.com/ivan/storage-project-back/pkg/jwt_service"
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
