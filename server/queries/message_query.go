package queries

import (
	"github.com/chayapol-p/cloud-midterm-server/models"
	"github.com/jmoiron/sqlx"
)

// MessageQueries struct for queries from Message model.
type MessageQueries struct {
	*sqlx.DB
}

// GetMessages method for getting all messages.
func (q *MessageQueries) GetMessages(time string, offset string, limit string) ([]models.OutputMessage, error) {
	// Define messages variable.
	messages := []models.OutputMessage{}

	// Define query string.
	// query := `SELECT
	// 	messages.uuid,
	// 	messages.author,
	// 	messages.message,
	// 	messages.likes,
	// 	updated_messages.author as updated_author,
	// 	updated_messages.message as updated_message,
	// 	updated_messages.likes as updated_likes
	// FROM
	// 	messages
	// LEFT JOIN updated_messages 
	// ON messages.uuid = updated_messages.uuid
	// WHERE messages.timestamp > $1 offset $2 limit $3;`

	query := `SELECT uuid, author, message, likes FROM messages WHERE timestamp > $1 offset $2 limit $3;`
	// Send query to database.
	err := q.Select(&messages, query, time, offset, limit)
	if err != nil {
		// Return empty object and error.
		return messages, err
	}

	// Return query result.
	return messages, nil
}

// GetMessage method for getting one message by given ID.
func (q *MessageQueries) GetMessage(uuid string) (models.Message, error) {
	// Define message variable.
	message := models.Message{}

	// Define query string.
	query := `SELECT * FROM messages WHERE uuid = $1`

	// Send query to database.
	err := q.Get(&message, query, uuid)
	if err != nil {
		// Return empty object and error.
		return message, err
	}

	// Return query result.
	return message, nil
}

// CreateMessage method for creating message by given Message object.
func (q *MessageQueries) CreateMessage(m *models.Message) error {
	// Define query string.
	query := `INSERT INTO messages VALUES ($1, $2, $3, $4, $5)`

	// Send query to database.
	_, err := q.Exec(query, m.UUID, m.Timestamp, m.Author, m.Message, m.Likes)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateMessage method for updating message by given Message object.
func (q *MessageQueries) UpdateMessage(id string, m *models.Message) error {
	// Define query string.
	query := `UPDATE messages SET timestamp = $2, message = $3, author = $4, likes = $5 WHERE uuid = $1`

	// Send query to database.
	_, err := q.Exec(query, id, m.Timestamp, m.Message, m.Author, m.Likes)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteMessage method for delete message by given ID.
func (q *MessageQueries) DeleteMessage(id string) error {
	// Define query string.
	query := `DELETE FROM messages WHERE uuid = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
