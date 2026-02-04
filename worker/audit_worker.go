package worker

import (
	"time"

	"go-demo/pkg/logger"
)

// Event represents a single audit entry for business actions in the API.
type Event struct {
	Action    string    // e.g. CREATE, UPDATE, DELETE, READ
	Entity    string    // e.g. user, product
	EntityID  int       // database identifier when available
	Message   string    // human‑readable description
	Timestamp time.Time // time when the action occurred
}

// queueSize controls how many audit events can be buffered before we start
// dropping new ones. This keeps the API fast under backpressure.
const queueSize = 100

// EventQueue is the buffered channel used to hand off audit events to the
// background worker. It is initialized by StartWorker.
var EventQueue chan Event

// StartWorker initializes the global audit queue and starts a single
// background goroutine that processes audit events.
//
// This should be called once from main() during application startup.
func StartWorker() {
	EventQueue = make(chan Event, queueSize)

	go func() {
		for ev := range EventQueue {
			logger.Log.Info().
				Str("audit_action", ev.Action).
				Str("audit_entity", ev.Entity).
				Int("audit_entity_id", ev.EntityID).
				Str("audit_message", ev.Message).
				Time("audit_timestamp", ev.Timestamp).
				Msg("audit event")
		}
	}()
}

// NewEvent is a small helper to construct an Event with the current timestamp.
func NewEvent(action, entity string, entityID int, message string) Event {
	return Event{
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// Publish queues an audit event for asynchronous processing. The send is
// non‑blocking; if the queue is full the event is dropped and a warning is
// logged, ensuring API handlers are never slowed by audit logging.
func Publish(ev Event) {
	if EventQueue == nil {
		// Worker not started; avoid panics in tests or misconfiguration.
		logger.Log.Warn().Msg("audit worker not started; dropping event")
		return
	}

	select {
	case EventQueue <- ev:
	default:
		logger.Log.Warn().Msg("audit queue full; dropping event")
	}
}
