package main

import (
	"go-demo/config"
	"go-demo/database"
	"go-demo/pkg/logger"
	"go-demo/pkg/validator"
	"go-demo/worker"

	"github.com/robfig/cron/v3"
)

func main() {
	config.LoadEnv()
	logger.Init()
	// initialize validator in case workers need it for data validation
	validator.Init()

	database.Connect()

	logger.Log.Info().Msg("Starting background workers...")

	// start the background audit worker
	// This starts a goroutine that listens on a channel (or DB polling in v2)
	worker.StartWorker()

	// start the background notification worker
	// This also starts a goroutine
	worker.StartNotificationWorker()

	// Initialize Cron Scheduler
	c := cron.New()

	// Register the cleanup worker
	worker.RegisterCleanupWorker(c)

	// Start the cron scheduler (runs in its own goroutine)
	c.Start()
	logger.Log.Info().Msg("cron scheduler started")

	logger.Log.Info().Msg("Background workers started successfully")

	// Block forever to keep the workers running
	select {}
}
