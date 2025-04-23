package models

import (
	"time"
)

type User struct {
	ID        string     `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password_hash" json:"-"`
	IsActive  bool       `db:"is_active" json:"is_active"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitzero"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by,omitzero"`
	RoleID    string     `db:"role_id" json:"-"`
	Role      Role       `json:"role,omitempty"`
}
