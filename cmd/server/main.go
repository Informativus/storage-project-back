package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title Storage Project API
// @version 1.0
// @description REST API для проекта Storage Project. Позволяет работать с пользователями, папками и файлами.
// @termsOfService http://example.com/terms/

// @contact.name Ivan Popov

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	var routers = gin.Default()

	var reg *Registry = NewRegistry(cfg, routers)

	reg.Controllers.RegisterRoutes(routers)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	var serverInstance = NewServer(routers)

	var options Options = Options{
		Port: ":" + cfg.Port,
	}

	serverInstance.run(options)
}
