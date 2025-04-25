package models

import (
	"time"
)

type User struct {
	ID        string     `db:"id" json:"id,omitempty"`
	Username  string     `db:"username" json:"username,omitempty"`
	Email     string     `db:"email" json:"email,omitempty"`
	Password  string     `db:"password" json:"-"`
	IsActive  bool       `db:"is_active" json:"isActive,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
	CreatedBy string     `db:"created_by" json:"createdBy,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy,omitempty"`
	RoleID    string     `db:"role_id" json:"-"`
	Role      Role       `json:"role,omitempty"`
}
