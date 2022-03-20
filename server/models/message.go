package models

import (
	"database/sql"
	"time"
)

// Message struct to describe message object.
type Message struct {
	UUID      string    `db:"uuid" json:"uuid"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Author    string    `db:"author" json:"author"`
	Message   string    `db:"message" json:"message"`
	Likes     int       `db:"likes" json:"likes"`
}

type UpdatedMessage struct {
	UUID      string    `db:"uuid" json:"uuid"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Author    string    `db:"author" json:"author"`
	Message   string    `db:"message" json:"message"`
	Likes     int       `db:"likes" json:"likes"`
	IsDeleted bool      `db:"is_deleted" json:"is_deleted"`
}

type OutputMessage struct {
	UUID    string `db:"uuid" json:"uuid"`
	Author  string `db:"author" json:"author"`
	Message string `db:"message" json:"message"`
	Likes   int    `db:"likes" json:"likes"`
}

type MessageWithUpdate struct {
	UUID           string         `db:"uuid" json:"uuid"`
	Author         string         `db:"author" json:"author"`
	Message        string         `db:"message" json:"message"`
	Likes          int            `db:"likes" json:"likes"`
	UpdatedAuthor  sql.NullString `db:"updated_author" json:"updated_author"`
	UpdatedMessage sql.NullString `db:"updated_message" json:"updated_message"`
	UpdatedLikes   sql.NullInt32  `db:"updated_likes" json:"updated_likes"`
}

type OutputUpdatedMessage struct {
	UUID      string `db:"uuid" json:"uuid"`
	Author    string `db:"author" json:"author"`
	Message   string `db:"message" json:"message"`
	Likes     int    `db:"likes" json:"likes"`
	IsDeleted bool   `db:"is_deleted" json:"is_deleted"`
}
