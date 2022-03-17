package models

import (
	"time"
)

// Message struct to describe message object.
type Message struct {
	UUID      string    `db:"uuid" json:"uuid"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Message   string    `db:"message" json:"message"`
	Author    string    `db:"author" json:"author"`
	Likes     int       `db:"likes" json:"likes"`
}

type UpdatedMessage struct {
	UUID      string    `db:"uuid" json:"uuid"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Message   string    `db:"message" json:"message"`
	Author    string    `db:"author" json:"author"`
	Likes     int       `db:"likes" json:"likes"`
	IsDeleted bool      `db:"is_deleted" json:"is_deleted"`
}
