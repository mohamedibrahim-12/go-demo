package models

import (
	"time"
)

// AuditLog acts as a persistent queue for audit events.
type AuditLog struct {
	ID        uint      `gorm:"primaryKey"`
	Action    string    `gorm:"not null"`
	Entity    string    `gorm:"not null"`
	EntityID  int       `gorm:"not null"`
	Message   string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	
	// ProcessedAt is null when pending, set when worker handles it.
	// In a real queue, we might delete it, but keeping it is good for audit trail anyway.
	ProcessedAt *time.Time `gorm:"index"` 
	
	CreatedAt time.Time
}

// NotificationOutbox acts as a persistent queue (Outbox) for notifications.
// It uses a generic payload to allow different types of notifications.
type NotificationOutbox struct {
	ID        uint      `gorm:"primaryKey"`
	EventType string    `gorm:"not null"`          // e.g. WELCOME_EMAIL, PASSWORD_RESET
	Payload   string    `gorm:"not null;type:text"` // JSON payload
	Status    string    `gorm:"default:'PENDING';index"` // PENDING, PROCESSING, DONE, FAILED
	
	ProcessedAt *time.Time
	Error       string
	
	CreatedAt time.Time
}

// Deprecated: Scan NotificationJob from NotificationOutbox instead
type NotificationJob struct {
	ID        uint      `gorm:"primaryKey"`
	Type      string    `gorm:"not null"`
	Recipient string    `gorm:"not null"`
	Message   string    `gorm:"not null"`
	
	// Status tracking
	Status      string     `gorm:"default:'PENDING';index"` // PENDING, PROCESSED, FAILED
	ProcessedAt *time.Time
	Error       string
	
	CreatedAt time.Time
}
