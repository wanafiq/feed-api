package models

import "time"

type Role struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"` // "user", "moderator", "admin"
	Level       int       `db:"level" json:"level"`
	Description string    `db:"description" json:"description"`
	IsActive    bool      `db:"is_active" json:"isActive"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at,omitzero"`
	UpdatedBy   string    `json:"updated_by,omitzero"`
}
