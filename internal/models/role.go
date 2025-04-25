package models

import "time"

type Role struct {
	ID          string    `db:"id" json:"id,omitempty"`
	Name        string    `db:"name" json:"name,omitempty"`
	Level       int       `db:"level" json:"level,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
	IsActive    bool      `db:"is_active" json:"isActive,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitzero,omitempty"`
	CreatedBy   string    `json:"created_by,omitempty,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitzero,omitempty"`
	UpdatedBy   string    `json:"updated_by,omitzero,omitempty"`
}
