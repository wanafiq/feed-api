// File: /internal/models/post.go
// Purpose: Defines the Post model with content, tags, slugs, and auditing fields.

package models

import (
	"time" // For time and auditing fields
	// TODO: Additional imports if necessary
)

// Post represents a blog post.
type Post struct {
	ID            int       `json:"id"`
	AuthorID      int       `json:"author_id"` // Reference to the User.ID
	Content       string    `json:"content"`   // Markdown content
	ImageLinks    []string  `json:"image_links"`
	Tags          []string  `json:"tags"`
	Slug          string    `json:"slug"` // User-friendly URL
	PublishedDate time.Time `json:"published_date"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `json:"updated_by"`
	// TODO: Add additional fields if necessary
}
