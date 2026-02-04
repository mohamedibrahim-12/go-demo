package worker

import (
	"time"

	"go-demo/pkg/logger"
)

// NotificationJob represents a single notification task to be processed
// asynchronously by the background worker.
type NotificationJob struct {
	Type      string    // e.g. WELCOME_EMAIL, PASSWORD_RESET, etc.
	Recipient string    // email address or identifier of the recipient
	Message   string    // notification message content
	CreatedAt time.Time // timestamp when the job was created
}

// notificationQueueSize controls how many notification jobs can be buffered
// before we start dropping new ones. This keeps the API fast under backpressure.
const notificationQueueSize = 100

// NotificationQueue is the buffered channel used to hand off notification jobs
// to the background worker. It is initialized by StartNotificationWorker.
var NotificationQueue chan NotificationJob

// StartNotificationWorker initializes the global notification queue and starts
// a single background goroutine that processes notification jobs.
//
// The worker continuously listens on NotificationQueue and simulates sending
// emails by logging them. Failures are logged but never crash the application.
//
// This should be called once from main() during application startup.
func StartNotificationWorker() {
	NotificationQueue = make(chan NotificationJob, notificationQueueSize)

	go func() {
		// Worker loop: continuously process jobs from the queue
		for job := range NotificationQueue {
			// Simulate email sending by logging the notification.
			// In production, this would call an email service (SMTP, SendGrid, etc.)
			// Wrapped in a recover to ensure worker never crashes on failure.
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Log.Error().
							Interface("panic", r).
							Str("notification_type", job.Type).
							Str("recipient", job.Recipient).
							Msg("notification worker panic recovered; continuing")
					}
				}()

				// Simulate time-consuming email operation
				// In real implementation, this would be: smtp.Send(...)
				logger.Log.Info().
					Str("notification_type", job.Type).
					Str("recipient", job.Recipient).
					Str("message", job.Message).
					Time("created_at", job.CreatedAt).
					Msg("notification sent")
			}()
		}
	}()
}

// NewNotificationJob is a helper to construct a NotificationJob with the
// current timestamp.
func NewNotificationJob(jobType, recipient, message string) NotificationJob {
	return NotificationJob{
		Type:      jobType,
		Recipient: recipient,
		Message:   message,
		CreatedAt: time.Now(),
	}
}

// EnqueueNotification queues a notification job for asynchronous processing.
// The send is non-blocking; if the queue is full the job is dropped and a
// warning is logged, ensuring API handlers are never slowed by notification
// processing.
func EnqueueNotification(job NotificationJob) {
	if NotificationQueue == nil {
		// Worker not started; avoid panics in tests or misconfiguration.
		logger.Log.Warn().Msg("notification worker not started; dropping job")
		return
	}

	select {
	case NotificationQueue <- job:
		// Successfully enqueued
	default:
		// Queue full; drop the job to avoid blocking the API
		logger.Log.Warn().
			Str("notification_type", job.Type).
			Str("recipient", job.Recipient).
			Msg("notification queue full; dropping job")
	}
}
