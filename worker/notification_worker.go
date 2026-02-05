package worker

import (
	"time"

	"go-demo/database"
	"go-demo/models"
	"go-demo/pkg/logger"
)

// StartNotificationWorker initializes the notification worker that polls the database.
func StartNotificationWorker() {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // Poll every second
		defer ticker.Stop()

		for range ticker.C {
			processNotificationJobs()
		}
	}()
	logger.Log.Info().Msg("notification worker started (db polling)")
}

func processNotificationJobs() {
	if database.GormDB == nil {
		return
	}

	// Fetch up to 100 pending jobs
	var jobs []models.NotificationJob
	// Find jobs where Status is PENDING
	if err := database.GormDB.Where("status = ?", "PENDING").Limit(100).Order("created_at asc").Find(&jobs).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to fetch notification jobs")
		return
	}

	if len(jobs) == 0 {
		return
	}

	for _, job := range jobs {
		processSingleJob(job)
	}
}

func processSingleJob(job models.NotificationJob) {
	// Wrapper to handle panic recovery per job
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error().
				Interface("panic", r).
				Str("notification_type", job.Type).
				Str("recipient", job.Recipient).
				Msg("notification worker panic recovered; marking failed")
			
			// Mark as FAILED
			job.Status = "FAILED"
			job.Error = "Panic recovered"
			database.GormDB.Save(&job)
		}
	}()

	// Simulate processing
	logger.Log.Info().
		Str("notification_type", job.Type).
		Str("recipient", job.Recipient).
		Str("message", job.Message).
		Time("created_at", job.CreatedAt).
		Msg("notification sent")

	// Mark as PROCESSED
	now := time.Now()
	job.Status = "PROCESSED"
	job.ProcessedAt = &now
	if err := database.GormDB.Save(&job).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to update notification job status")
	}
}

// NewNotificationJob is a helper to construct a NotificationJob.
func NewNotificationJob(jobType, recipient, message string) models.NotificationJob {
	return models.NotificationJob{
		Type:      jobType,
		Recipient: recipient,
		Message:   message,
		CreatedAt: time.Now(),
		Status:    "PENDING",
	}
}

// EnqueueNotification queues a notification job in the database.
func EnqueueNotification(job models.NotificationJob) {
	if database.GormDB == nil {
		logger.Log.Warn().Msg("notification enqueue skipped: no DB connection")
		return
	}

	if err := database.GormDB.Create(&job).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to enqueue notification job")
	}
}
