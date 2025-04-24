// File: /internal/models/comment.go
// Purpose: Defines the Comment model for blog posts with auditing fields.

package models

import (
	"time" // For times and auditing fields
	// TODO: Additional imports if needed
)

// Comment represents a comment on a blog post.
type Comment struct {
	ID        string    `db:"id" json:"id"`
	AuthorID  string    `db:"author_id" json:"author_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitzero"`
	UpdatedBy string    `db:"updated_by" json:"updated_by,omitzero"`
}
