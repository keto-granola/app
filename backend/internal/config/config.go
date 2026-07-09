package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
	EnvironmentTest        Environment = "test"
	EnvironmentCI          Environment = "ci"

	APIPrefix   = "api"
	APIVersion  = "v1"
	APIBasePath = "/" + APIPrefix + "/" + APIVersion
)

type App struct {
	Port        string
	ClientURL   string
	DbURL       string
	LogLevel    slog.Level
	Environment Environment
}

type Environment string

var validEnvironments = []Environment{
	EnvironmentDevelopment,
	EnvironmentProduction,
	EnvironmentTest,
	EnvironmentCI,
}

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func ParseEnv() (*App, error) {
	// Ignore error because in production there will be no .env file, env vars will be passed
	// in at runtime via docker run command/docker-compose
	_ = godotenv.Load()

	envVars := map[string]string{
		"SERVER_PORT": "",
		"DB_URL":      "",
		"LOG_LEVEL":   "",
		"CLIENT_URL":  "",
		"ENVIRONMENT": "",
	}

	for key := range envVars {
		value := os.Getenv(key)
		if value == "" {
			return nil, fmt.Errorf("%s environment variable is not set", key)
		}
		envVars[key] = value
	}

	logLevel, ok := logLevelMap[envVars["LOG_LEVEL"]]
	if !ok {
		return nil, errors.New("LOG_LEVEL should be one of debug|info|warning|error")
	}

	environment := Environment(envVars["ENVIRONMENT"])
	if !slices.Contains(validEnvironments, environment) {
		return nil, errors.New("ENVIRONMENT should be one of development|production|test|ci")
	}

	return &App{
		Port:        envVars["SERVER_PORT"],
		DbURL:       envVars["DB_URL"],
		ClientURL:   envVars["CLIENT_URL"],
		LogLevel:    logLevel,
		Environment: environment,
	}, nil
}
