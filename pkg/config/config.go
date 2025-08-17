package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DatabaseUser     string
	DatabasePassword string
	DatabaseDb       string
	DatabasePort     string
	DatabaseHost     string
	Port             string
	StoragePath      string
	SecretKey        string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	cfg := &Config{
		DatabaseUser:     getStrFromEnv("DATABASE_USER", true),
		DatabasePassword: getStrFromEnv("DATABASE_PASSWORD", true),
		DatabaseDb:       getStrFromEnv("DATABASE_DB", true),
		DatabasePort:     getStrFromEnv("DATABASE_PORT", true),
		DatabaseHost:     getStrFromEnv("DATABASE_HOST", true),
		Port:             getStrFromEnv("PORT", true),
		StoragePath:      getStrFromEnv("STORAGE_PATH", true),
		SecretKey:        getStrFromEnv("SECRET_KEY", true),
	}

	return cfg, nil
}

func getStrFromEnv(key string, req bool) string {
	val := os.Getenv(key)

	if req && val == "" {
		log.Fatal().Str("env_err", key).Msg("environment variable is required")
	}

	return val
}
