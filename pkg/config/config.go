package config

import (
	"os"
	"strconv"
	"strings"
	"time"

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
	ExpiresIn        int64
	GpgPublicKeyPath string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	CacheLifetime    time.Duration
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
		ExpiresIn:        getTimeFromEnv("EXPIRESIN", true, "seconds"),
		GpgPublicKeyPath: getStrFromEnv("GPG_PUBLIC_KEY_PATH", true),
		RedisHost:        getStrFromEnv("REDIS_HOST", true),
		RedisPort:        getStrFromEnv("REDIS_PORT", true),
		RedisPassword:    getStrFromEnv("REDIS_PASSWORD", false),
		CacheLifetime:    time.Duration(getTimeFromEnv("CACHE_LIFETIME", true, "seconds")) * time.Second,
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

func getTimeFromEnv(key string, req bool, unit string) int64 {
	val := os.Getenv(key)
	if req && val == "" {
		log.Fatal().Str("env_err", key).Msg("environment variable is required")
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal().Str("env_err", key).Str("val", val).Err(err).Msg("cannot convert env to int")
	}

	multiplier := int64(1)
	switch strings.ToLower(unit) {
	case "seconds":
		multiplier = 1
	case "minutes":
		multiplier = 60
	case "hours":
		multiplier = 3600
	case "days":
		multiplier = 86400
	default:
		log.Fatal().Str("unit", unit).Msg("invalid time unit")
	}

	return int64(valInt) * multiplier
}
