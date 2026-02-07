package tests

import (
	"encoding/json"
	"testing"
	"time"

	"go-demo/database"
	"go-demo/models"
	"go-demo/repositories"
	"go-demo/worker"

	. "github.com/onsi/gomega"
)

// TestNotificationWorkerEnqueue verifies that a notification job is persisted to the database.
func TestNotificationWorkerEnqueue(t *testing.T) {
	RegisterTestingT(t)

	// Create and enqueue a test notification job via Outbox
	payloadMap := map[string]string{
		"recipient": "test_enqueue@example.com",
		"message":   "Test enqueue message",
	}
	payloadBytes, _ := json.Marshal(payloadMap)

	// Use repository to create outbox entry
	repositories.CreateNotificationOutbox("WELCOME_EMAIL", string(payloadBytes))

	// Verify it exists in DB
	var savedJob models.NotificationOutbox
	// We need to parse the JSON payload in a real app, but for this test checking persistence is enough
	// or we can use LIKE query if we don't want to parse
	err := database.GormDB.Where("payload LIKE ?", "%test_enqueue@example.com%").Last(&savedJob).Error
	Expect(err).To(BeNil())
	Expect(savedJob.EventType).To(Equal("WELCOME_EMAIL"))
	Expect(savedJob.Status).To(Equal("PENDING"))
}

// TestNotificationWorkerProcessesJob verifies that if we run the worker logic, it processes the job.
func TestNotificationWorkerProcessesJob(t *testing.T) {
	RegisterTestingT(t)

	// Start the worker (it will poll every second)
	worker.StartNotificationWorker()

	// Enqueue a job via Outbox
	payloadMap := map[string]string{
		"recipient": "process_test@example.com",
		"message":   "Reset your password",
	}
	payloadBytes, _ := json.Marshal(payloadMap)
	repositories.CreateNotificationOutbox("RESET_PASSWORD", string(payloadBytes))

	// Wait for worker to pick it up (poll interval is 1s)
	Eventually(func() string {
		var j models.NotificationOutbox
		database.GormDB.Where("payload LIKE ?", "%process_test@example.com%").Last(&j)
		return j.Status
	}, 3*time.Second, 500*time.Millisecond).Should(Equal("DONE"))

	// Verify ProcessedAt is set
	var j models.NotificationOutbox
	database.GormDB.Where("payload LIKE ?", "%process_test@example.com%").Last(&j)
	Expect(j.ProcessedAt).NotTo(BeNil())
}
