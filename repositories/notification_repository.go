package repositories

import (
	"go-demo/database"
	"go-demo/models"
	"go-demo/pkg/logger"
	"time"
)

// CreateNotificationOutbox inserts a new job into the notification outbox.
func CreateNotificationOutbox(eventType, payload string) {
	if database.GormDB == nil {
		logger.Log.Warn().Msg("notification outbox insert skipped: no DB connection")
		return
	}

	outboxMsg := models.NotificationOutbox{
		EventType: eventType,
		Payload:   payload,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	if err := database.GormDB.Create(&outboxMsg).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to insert into notification outbox")
	}
}
