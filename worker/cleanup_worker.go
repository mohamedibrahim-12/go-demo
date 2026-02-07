package worker

import (
	"time"

	"go-demo/database"
	"go-demo/models"
	"go-demo/pkg/logger"

	"github.com/robfig/cron/v3"
)

// DefaultRetention is the period after which records are considered expired
// and removed by the cleanup worker (e.g. 7 days).
const DefaultRetention = 1 * time.Minute

// CleanupSchedule is the cron expression for how often the worker runs.
const CleanupSchedule = "@every 1m"

// RegisterCleanupWorker registers the cleanup job with the provided cron scheduler.
func RegisterCleanupWorker(c *cron.Cron) {
	_, err := c.AddFunc(CleanupSchedule, func() {
		// Wrapper to handle panic recovery per job run
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error().Interface("panic", r).Msg("cleanup worker panic recovered")
			}
		}()
		runCleanup(DefaultRetention)
	})

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to register cleanup worker")
	}

	logger.Log.Info().
		Str("schedule", CleanupSchedule).
		Dur("retention", DefaultRetention).
		Msg("cleanup worker registered")
}

// runCleanup executes the cleanup logic once: deletes users and products
// where created_at is older than the given retention. Logs deleted counts.
func runCleanup(retention time.Duration) {
	logger.Log.Info().Msg("cleanup job executing") // Log when job starts

	if database.GormDB == nil {
		logger.Log.Warn().Msg("cleanup skipped: no DB connection")
		return
	}

	cutoff := time.Now().Add(-retention)

	// Delete old users
	resultUsers := database.GormDB.Unscoped().Delete(
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
	resultProducts := database.GormDB.Unscoped().Delete(
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
	} else {
		// Optional: log that nothing was deleted to verify it ran
		logger.Log.Info().Msg("cleanup run completed; no records deleted")
	}
}

// RunCleanupOnce runs the cleanup logic once with the given retention.
// Used by tests to exercise cleanup without waiting for the cron.
func RunCleanupOnce(retention time.Duration) {
	runCleanup(retention)
}
