package tests

import (
	"os"
	"testing"

	"go-demo/config"
	"go-demo/database"
	"go-demo/pkg/logger"
	"go-demo/pkg/validator"
	"go-demo/worker"
)

func TestMain(m *testing.M) {
	// Ensure test env is loaded before connecting to DB
	os.Setenv("APP_ENV", "test")
	config.LoadEnv()
	logger.Init()
	validator.Init()
	database.Connect()

	// start audit worker in test binary so audit events are processed here too
	worker.StartWorker()

	// start notification worker in test binary so notifications are processed here too
	worker.StartNotificationWorker()

	// Run tests
	code := m.Run()

	logger.Log.Info().Msg("Tests finished")

	os.Exit(code)
}
