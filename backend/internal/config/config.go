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
	EnvironmentTest        Environment = "test"
	EnvironmentStaging     Environment = "staging"
	EnvironmentProduction  Environment = "production"
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
	Auth        *Auth
}

type Config struct {
	AuthOverride     *map[string]interface{} `json:"databaseAuthVariableOverride"`
	DatabaseURL      string                  `json:"databaseURL"`
	ProjectID        string                  `json:"projectId"`
	ServiceAccountID string                  `json:"serviceAccountId"`
	StorageBucket    string                  `json:"storageBucket"`
}

type Auth struct {
	Credentials string
}

type Environment string

var validEnvironments = []Environment{
	EnvironmentDevelopment,
	EnvironmentTest,
	EnvironmentStaging,
	EnvironmentProduction,
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
		"SERVER_PORT":    "",
		"DB_URL":         "",
		"LOG_LEVEL":      "",
		"CLIENT_URL":     "",
		"ENVIRONMENT":    "",
		"FIREBASE_CREDS": "",
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
		return nil, errors.New("ENVIRONMENT should be one of development|test|staging|production|ci")
	}

	return &App{
		Port:        envVars["SERVER_PORT"],
		DbURL:       envVars["DB_URL"],
		ClientURL:   envVars["CLIENT_URL"],
		LogLevel:    logLevel,
		Environment: environment,
		Auth: &Auth{
			Credentials: envVars["FIREBASE_CREDS"],
		},
	}, nil
}
