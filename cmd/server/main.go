package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/ivan/storage-project-back/pkg/config"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	var reg *Registry = NewRegistry(cfg)

	var routers = gin.Default()

	reg.Controllers.RegisterRoutes(routers)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	var serverInstance = NewServer(routers)

	var options Options = Options{
		Port: ":" + cfg.Port,
	}

	serverInstance.run(options)
}
