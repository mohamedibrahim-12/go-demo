package worker

import (
	"time"

	"go-demo/database"
	"go-demo/models"
	"go-demo/pkg/logger"
)

// StartWorker initializes the audit worker that polls the database for new audit logs.
func StartWorker() {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // Poll every second
		defer ticker.Stop()

		for range ticker.C {
			processAuditLogs()
		}
	}()
	logger.Log.Info().Msg("audit worker started (db polling)")
}

func processAuditLogs() {
	if database.GormDB == nil {
		return
	}

	// Fetch up to 100 pending logs
	var logs []models.AuditLog
	// Find logs where ProcessedAt is NULL
	if err := database.GormDB.Where("processed_at IS NULL").Limit(100).Order("created_at asc").Find(&logs).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to fetch audit logs")
		return
	}

	if len(logs) == 0 {
		return
	}

	for _, logEntry := range logs {
		// "Process" the log by actual logging
		logger.Log.Info().
			Str("audit_action", logEntry.Action).
			Str("audit_entity", logEntry.Entity).
			Int("audit_entity_id", logEntry.EntityID).
			Str("audit_message", logEntry.Message).
			Time("audit_timestamp", logEntry.Timestamp).
			Msg("audit event processed")

		// Mark as processed
		now := time.Now()
		logEntry.ProcessedAt = &now
		if err := database.GormDB.Save(&logEntry).Error; err != nil {
			logger.Log.Error().Err(err).Msg("failed to mark audit log as processed")
		}
	}
}

// NewEvent helper is no longer strictly needed but kept for compatibility if referenced elsewhere.
// In this refactor, we usually just create the struct directly or use this helper to create the struct
// before passing to Publish.
func NewEvent(action, entity string, entityID int, message string) models.AuditLog {
	return models.AuditLog{
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		Message:   message,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}
}

// Publish writes an audit event to the database queue.
func Publish(ev models.AuditLog) {
	if database.GormDB == nil {
		logger.Log.Warn().Msg("audit publish skipped: no DB connection")
		return
	}

	if err := database.GormDB.Create(&ev).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to publish audit event")
	}
}
