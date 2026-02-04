package config

import (
	"os"
	"path/filepath"

	"go-demo/pkg/logger"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	root, _ := os.Getwd()

	// Move to project root if running from tests/
	if filepath.Base(root) == "tests" {
		root = filepath.Dir(root)
	}

	env := os.Getenv("APP_ENV")

	var envFile string
	if env == "test" {
		envFile = filepath.Join(root, ".env.test")
		logger.Log.Info().Msg("loading test environment")
	} else {
		envFile = filepath.Join(root, ".env")
		logger.Log.Info().Msg("loading production environment")
	}

	logger.Log.Info().Str("envFile", envFile).Msg("loading env file")

	if err := godotenv.Load(envFile); err != nil {
		logger.Log.Fatal().Err(err).Str("envFile", envFile).Msg("failed to load env file")
	}
}
