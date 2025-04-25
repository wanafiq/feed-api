package models

import (
	"time"
)

type Token struct {
	ID        string    `db:"id" json:"id,omitempty"`
	Type      string    `db:"type" json:"type,omitempty"`
	Value     string    `db:"value" json:"value,omitempty"`
	ExpiredAt time.Time `db:"expired_at" json:"expiredAt,omitempty"`
	UserID    string    `db:"user_id" json:"userId,omitempty"`
}
