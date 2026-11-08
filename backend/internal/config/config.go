package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
	MaxBotToken   string
	RedisDB       int
}

func Load() (Config, error) {
	appEnv := getEnv("APP_ENV", "development")
	envFile, err := envFileFor(appEnv)
	if err != nil {
		return Config{}, err
	}
	if err := godotenv.Load(envFile); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load environment file %s: %w", envFile, err)
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))

	if err != nil {
		return Config{}, fmt.Errorf("Invalid REDIS_DB: %w", err)
	}

	cfg := Config{
		Port:          getEnv("APP_PORT", "8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		MaxBotToken:   getEnv("MAX_BOT_TOKEN", ""),
		RedisDB:       redisDB,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is not set")
	}

	return cfg, nil
}

func envFileFor(appEnv string) (string, error) {
	switch appEnv {
	case "development", "dev":
		return ".env.dev", nil
	case "production", "prod":
		return ".env.prod", nil
	default:
		return "", fmt.Errorf("unsupported APP_ENV %q: expected development or production", appEnv)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
