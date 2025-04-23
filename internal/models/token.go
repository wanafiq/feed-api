package models

import (
	"time"
)

type Token struct {
	ID        string    `db:"id" json:"id"`
	Type      string    `db:"type" json:"type"`
	Value     string    `db:"value" json:"value"`
	ExpiredAt time.Time `db:"expired_at" json:"expiredAt"`
	UserID    string    `db:"user_id" json:"userId"`
}
