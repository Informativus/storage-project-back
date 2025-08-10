package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDb       string
	PostgresPort     string
	Port             string
	StoragePath      string
	SecretKey        string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	cfg := &Config{
		PostgresUser:     getStrFromEnv("POSTGRES_USER", true),
		PostgresPassword: getStrFromEnv("POSTGRES_PASSWORD", true),
		PostgresDb:       getStrFromEnv("POSTGRES_DB", true),
		PostgresPort:     getStrFromEnv("POSTGRES_PORT", true),
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
