package tests

import (
	"testing"
	"time"

	"go-demo/database"
	"go-demo/models"
	"go-demo/worker"

	. "github.com/onsi/gomega"
)

// TestNotificationWorkerEnqueue verifies that a notification job is persisted to the database.
func TestNotificationWorkerEnqueue(t *testing.T) {
	RegisterTestingT(t)

	// Create and enqueue a test notification job
	job := worker.NewNotificationJob(
		"WELCOME_EMAIL",
		"test_enqueue@example.com",
		"Test enqueue message",
	)

	// Enqueue should insert into DB
	worker.EnqueueNotification(job)

	// Verify it exists in DB
	var savedJob models.NotificationJob
	err := database.GormDB.Where("recipient = ?", "test_enqueue@example.com").Last(&savedJob).Error
	Expect(err).To(BeNil())
	Expect(savedJob.Type).To(Equal("WELCOME_EMAIL"))
	Expect(savedJob.Status).To(Equal("PENDING"))
}

// TestNotificationWorkerProcessesJob verifies that if we run the worker logic, it processes the job.
// Note: We cannot easily call the unexported process loop, but we can start the worker and wait.
func TestNotificationWorkerProcessesJob(t *testing.T) {
	RegisterTestingT(t)

	// Start the worker (it will poll every second)
	worker.StartNotificationWorker()

	// Enqueue a job
	recipient := "process_test@example.com"
	job := worker.NewNotificationJob("RESET_PASSWORD", recipient, "Reset your password")
	worker.EnqueueNotification(job)

	// Wait for worker to pick it up (poll interval is 1s)
	Eventually(func() string {
		var j models.NotificationJob
		database.GormDB.Where("recipient = ?", recipient).Last(&j)
		return j.Status
	}, 3*time.Second, 500*time.Millisecond).Should(Equal("PROCESSED"))

	// Verify ProcessedAt is set
	var j models.NotificationJob
	database.GormDB.Where("recipient = ?", recipient).Last(&j)
	Expect(j.ProcessedAt).NotTo(BeNil())
}
