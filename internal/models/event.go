// File: /internal/models/event.go
// Purpose: Defines the Event model for logging asynchronous events (e.g., email events).

package models

import (
	"time" // For auditing and event timestamps
	// TODO: Add other necessary imports if needed
)

// Event represents an asynchronous event logged by the system.
type Event struct {
	ID        int       `json:"id"`
	EventType string    `json:"event_type"` // E.g., "email_confirmation", "password_reset"
	Payload   string    `json:"payload"`    // JSON or string representation of the event data
	Status    string    `json:"status"`     // E.g., "pending", "completed", "failed"
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	// TODO: Add additional fields if necessary
}
