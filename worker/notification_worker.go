package worker

import (
	"encoding/json"
	"time"

	"go-demo/database"
	"go-demo/models"
	"go-demo/pkg/logger"
)

// StartNotificationWorker initializes the notification worker that polls the notification outbox.
func StartNotificationWorker() {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // Poll every second
		defer ticker.Stop()

		for range ticker.C {
			processNotificationOutbox()
		}
	}()
	logger.Log.Info().Msg("notification worker started (outbox polling)")
}

func processNotificationOutbox() {
	if database.GormDB == nil {
		return
	}

	// Fetch up to 100 pending messages from outbox
	var messages []models.NotificationOutbox
	// Find messages where Status is PENDING
	if err := database.GormDB.Where("status = ?", "PENDING").Limit(100).Order("created_at asc").Find(&messages).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to fetch notification outbox messages")
		return
	}

	if len(messages) == 0 {
		return
	}

	for _, msg := range messages {
		processSingleMessage(msg)
	}
}

func processSingleMessage(msg models.NotificationOutbox) {
	// Log: PICKED
	logger.Log.Info().
		Uint("job_id", msg.ID).
		Str("event_type", msg.EventType).
		Str("lifecycle_action", "picked").
		Msg("notification job picked from outbox")

	// Wrapper to handle panic recovery per message
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error().
				Interface("panic", r).
				Uint("job_id", msg.ID).
				Str("event_type", msg.EventType).
				Str("lifecycle_action", "failed").
				Msg("notification worker panic recovered; marking failed")
			
			// Mark as FAILED
			msg.Status = "FAILED"
			msg.Error = "Panic recovered"
			database.GormDB.Save(&msg)
		}
	}()

	// Mark as PROCESSING
	msg.Status = "PROCESSING"
	if err := database.GormDB.Save(&msg).Error; err != nil {
		logger.Log.Error().
			Err(err).
			Uint("job_id", msg.ID).
			Str("event_type", msg.EventType).
			Str("lifecycle_action", "failed").
			Msg("failed to mark outbox message as processing")
		return // Skip processing if we can't lock/status update
	}

	// Log: PROCESSING
	logger.Log.Info().
		Uint("job_id", msg.ID).
		Str("event_type", msg.EventType).
		Str("lifecycle_action", "processing").
		Msg("notification job status updated to processing")

	// Payload struct for generic notifications (simplified for this demo)
	type NotificationPayload struct {
		Recipient string `json:"recipient"`
		Message   string `json:"message"`
	}

	var payload NotificationPayload
	if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
		logger.Log.Error().
			Err(err).
			Uint("job_id", msg.ID).
			Str("event_type", msg.EventType).
			Str("lifecycle_action", "failed").
			Str("payload", msg.Payload).
			Msg("failed to unmarshal payload")
		
		msg.Status = "FAILED"
		msg.Error = err.Error()
		database.GormDB.Save(&msg)
		return
	}

	// Simulate processing / Sending Email
	// In a real system, you'd call an external service here.
	logger.Log.Info().
		Uint("job_id", msg.ID).
		Str("event_type", msg.EventType).
		Str("lifecycle_action", "sent").
		Str("recipient", payload.Recipient).
		Str("message", payload.Message).
		Time("created_at", msg.CreatedAt).
		Msg("notification logic executed successfully")

	// Mark as PROCESSED/DONE
	now := time.Now()
	msg.Status = "DONE"
	msg.ProcessedAt = &now
	if err := database.GormDB.Save(&msg).Error; err != nil {
		logger.Log.Error().
			Err(err).
			Uint("job_id", msg.ID).
			Str("event_type", msg.EventType).
			Str("lifecycle_action", "failed").
			Msg("failed to update notification outbox status to DONE")
		return
	}

	// Log: DONE
	logger.Log.Info().
		Uint("job_id", msg.ID).
		Str("event_type", msg.EventType).
		Str("lifecycle_action", "done").
		Msg("notification job completed successfully")
}
