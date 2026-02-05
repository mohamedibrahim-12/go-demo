package main

import (
	"go-demo/config"
	"go-demo/database"
	"go-demo/pkg/logger"
	"go-demo/pkg/validator"
	"go-demo/worker"
)

func main() {
	config.LoadEnv()
	logger.Init()
	// initialize validator in case workers need it for data validation
	validator.Init()

	database.Connect()

	logger.Log.Info().Msg("Starting background workers...")

	// start the background audit worker
	// This starts a goroutine that listens on a channel
	worker.StartWorker()

	// start the background notification worker
	// This also starts a goroutine
	worker.StartNotificationWorker()

	// start the data cleanup worker once; runs periodically via time.Ticker
	worker.StartCleanupWorker()

	logger.Log.Info().Msg("Background workers started successfully")

	// Block forever to keep the workers running
	select {}
}
