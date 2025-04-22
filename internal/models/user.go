// File: /internal/models/user.go
// Purpose: Defines the User model, including auditing fields.

package models

import (
	"time" // For time-related auditing fields
	// TODO: Add other necessary imports if needed
)

// User represents a user in the blogging platform.
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password_hash"`
	Role      string    `json:"role"` // "User", "Moderator", "Admin"
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	// TODO: Add additional fields if necessary
}
