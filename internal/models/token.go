// File: /internal/models/token.go
// Purpose: Defines the Token model for managing confirmation and password reset tokens.

package models

import (
	"time" // For expiration fields and auditing
	// TODO: Add other necessary imports if needed
)

// Token represents an authentication token (confirmation, password reset, etc.).
type Token struct {
	ID         int       `json:"id"`
	TokenValue string    `json:"token_value"` // Secure/hashed token value
	TokenType  string    `json:"token_type"`  // E.g., "confirmation", "password_reset"
	UserID     int       `json:"user_id"`     // Reference to the User.ID
	Expiration time.Time `json:"expiration"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  string    `json:"updated_by"`
	// TODO: Add additional fields if necessary
}
