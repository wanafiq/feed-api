package models

import (
	"time"
)

type Post struct {
	ID          string     `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Slug        string     `db:"slug" json:"slug"`
	Content     string     `db:"content" json:"content"`
	IsPublished bool       `db:"is_published" json:"is_published"`
	PublishedAt *time.Time `db:"published_at" json:"published_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	CreatedBy   string     `db:"created_by" json:"created_by"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitzero"`
	UpdatedBy   *string    `db:"updated_by" json:"updated_by,omitzero"`
	AuthorID    string     `db:"author_id" json:"author_id"`
	Tags        []Tag      `json:"tags,omitempty,omitzero"`
	Author      User       `json:"author,omitempty,omitzero"`
}
