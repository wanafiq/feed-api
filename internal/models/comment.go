// File: /internal/models/comment.go
// Purpose: Defines the Comment model for blog posts with auditing fields.

package models

import (
	"time" // For times and auditing fields
	// TODO: Additional imports if needed
)

// Comment represents a comment on a blog post.
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`   // Reference to the Post.ID
	AuthorID  int       `json:"author_id"` // Reference to the User.ID
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	// TODO: Add additional fields if necessary
}
