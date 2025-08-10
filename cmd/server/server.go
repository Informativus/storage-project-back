package main

import (
	"github.com/rs/zerolog/log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Routers *gin.Engine
}

type Options struct {
	Port string
}

func NewServer(routers *gin.Engine) *Server {
	routers.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	return &Server{Routers: routers}
}

func (s *Server) run(options Options) {
	err := s.Routers.Run(options.Port)

	if err != nil {
		log.Fatal().
			Err(err).
			Str("port", options.Port).
			Msg("Failed to start HTTP server")
	}

}
