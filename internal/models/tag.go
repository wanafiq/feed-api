package models

type Tag struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
