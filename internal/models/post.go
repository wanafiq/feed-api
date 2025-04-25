package models

import (
	"time"
)

type Post struct {
	ID          string     `db:"id" json:"id,omitempty"`
	Title       string     `db:"title" json:"title,omitempty"`
	Slug        string     `db:"slug" json:"slug,omitempty"`
	Content     string     `db:"content" json:"content,omitempty"`
	IsPublished bool       `db:"is_published" json:"isPublished,omitempty"`
	PublishedAt *time.Time `db:"published_at" json:"publishedAt,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt,omitempty"`
	CreatedBy   string     `db:"created_by" json:"createdBy,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	UpdatedBy   *string    `db:"updated_by" json:"updatedBy,omitempty"`
	AuthorID    string     `db:"author_id" json:"authorId,omitempty"`

	Tags   []Tag `json:"tags,omitempty"`
	Author User  `json:"author,omitempty"`
}

type PostFilter struct {
	Offset   int        `json:"offset,omitempty"`
	Limit    int        `json:"limit,omitempty"`
	Search   string     `json:"search,omitempty"`
	Sort     string     `json:"sort,omitempty"` // "asc" or "desc"
	DateFrom *time.Time `json:"date_from,omitempty"`
	DateTo   *time.Time `json:"date_to,omitempty"`
	Tags     []string   `json:"tags,omitempty"`
}
