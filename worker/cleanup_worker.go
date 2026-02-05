package worker

import (
	"time"

	"go-demo/database"
	"go-demo/models"
	"go-demo/pkg/logger"
)

// DefaultRetention is the period after which records are considered expired
// and removed by the cleanup worker (e.g. 7 days).
const DefaultRetention = 1 * time.Minute

// CleanupInterval is how often the worker runs (e.g. every 1 minute).
const CleanupInterval = 1 * time.Minute

// StartCleanupWorker starts a background goroutine that periodically deletes
// old records from users and products tables based on created_at.
//
// Scheduler: uses time.Ticker to run every CleanupInterval (1 minute).
// Retention: records with created_at older than DefaultRetention (7 days) are deleted.
// Errors are logged; the worker never panics or crashes the application.
//
// Call once from main() after database.Connect().
func StartCleanupWorker() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error().Interface("panic", r).Msg("cleanup worker panic recovered; worker stopped")
			}
		}()

		ticker := time.NewTicker(CleanupInterval)
		defer ticker.Stop()

		// Run once immediately, then on each tick
		runCleanup(DefaultRetention)
		for range ticker.C {
			runCleanup(DefaultRetention)
		}
	}()
	logger.Log.Info().Dur("interval", CleanupInterval).Dur("retention", DefaultRetention).Msg("cleanup worker started")
}

// runCleanup executes the cleanup logic once: deletes users and products
// where created_at is older than the given retention. Logs deleted counts.
func runCleanup(retention time.Duration) {
	if database.GormDB == nil {
		logger.Log.Warn().Msg("cleanup skipped: no DB connection")
		return
	}

	cutoff := time.Now().Add(-retention)

	// Delete old users
	resultUsers := database.GormDB.Delete(
		&models.User{},
		"created_at < ?",
		cutoff,
	)
	if resultUsers.Error != nil {
		logger.Log.Error().Err(resultUsers.Error).Msg("cleanup users failed")
		return
	}
	usersDeleted := resultUsers.RowsAffected

	// Delete old products
	resultProducts := database.GormDB.Delete(
		&models.Product{},
		"created_at < ?",
		cutoff,
	)
	if resultProducts.Error != nil {
		logger.Log.Error().Err(resultProducts.Error).Msg("cleanup products failed")
		return
	}
	productsDeleted := resultProducts.RowsAffected

	if usersDeleted > 0 || productsDeleted > 0 {
		logger.Log.Info().
			Int64("users_deleted", usersDeleted).
			Int64("products_deleted", productsDeleted).
			Msg("cleanup run completed")
	}
}

// RunCleanupOnce runs the cleanup logic once with the given retention.
// Used by tests to exercise cleanup without waiting for the ticker.
func RunCleanupOnce(retention time.Duration) {
	runCleanup(retention)
}
