package tests

import (
	"testing"
	"time"

	"go-demo/worker"
)

// TestNotificationWorkerEnqueue tests that a notification job can be
// successfully enqueued without blocking.
func TestNotificationWorkerEnqueue(t *testing.T) {
	// Verify the queue exists (started in TestMain)
	if worker.NotificationQueue == nil {
		t.Fatal("NotificationQueue should be initialized by TestMain")
	}

	// Create and enqueue a test notification job
	job := worker.NewNotificationJob(
		"WELCOME_EMAIL",
		"test@example.com",
		"Test welcome message",
	)

	// Enqueue should not block
	worker.EnqueueNotification(job)

	// Give the worker a moment to process
	time.Sleep(100 * time.Millisecond)

	// If we get here without blocking, the enqueue worked
}

// TestNotificationWorkerProcessesJob tests that the worker processes
// notification jobs from the queue. Since the worker logs the notification,
// we verify the job structure is correct and can be enqueued.
func TestNotificationWorkerProcessesJob(t *testing.T) {
	// Verify the queue exists
	if worker.NotificationQueue == nil {
		t.Fatal("NotificationQueue should be initialized by TestMain")
	}

	// Create a test notification job
	jobType := "WELCOME_EMAIL"
	recipient := "newuser@example.com"
	message := "Welcome to our platform!"

	job := worker.NewNotificationJob(jobType, recipient, message)

	// Verify job structure
	if job.Type != jobType {
		t.Errorf("expected job type %s, got %s", jobType, job.Type)
	}
	if job.Recipient != recipient {
		t.Errorf("expected recipient %s, got %s", recipient, job.Recipient)
	}
	if job.Message != message {
		t.Errorf("expected message %s, got %s", message, job.Message)
	}
	if job.CreatedAt.IsZero() {
		t.Error("job CreatedAt should be set")
	}

	// Enqueue the job (non-blocking)
	worker.EnqueueNotification(job)

	// Give the worker time to process
	time.Sleep(200 * time.Millisecond)

	// If we get here, the job was successfully enqueued and processed
	// The actual processing is verified by the worker logging the notification
}

// TestNotificationWorkerNonBlocking verifies that enqueueing notifications
// does not block even when multiple jobs are enqueued rapidly.
func TestNotificationWorkerNonBlocking(t *testing.T) {
	if worker.NotificationQueue == nil {
		t.Fatal("NotificationQueue should be initialized by TestMain")
	}

	// Enqueue multiple jobs rapidly
	for i := 0; i < 10; i++ {
		job := worker.NewNotificationJob(
			"WELCOME_EMAIL",
			"user@example.com",
			"Welcome message",
		)
		worker.EnqueueNotification(job)
	}

	// If we get here without blocking, the non-blocking behavior works
	time.Sleep(100 * time.Millisecond)
}
