package models

import (
	"time"
)

type Comment struct {
	ID        string    `db:"id" json:"id,omitempty"`
	AuthorID  string    `db:"author_id" json:"authorId,omitempty"`
	Content   string    `db:"content" json:"content,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt,omitempty"`
	CreatedBy string    `db:"created_by" json:"createdBy,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	UpdatedBy string    `db:"updated_by" json:"updatedBy,omitempty"`
}
