package models

type Tag struct {
	ID   string `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}
