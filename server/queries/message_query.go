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
func (q *MessageQueries) GetMessages() ([]models.Message, error) {
	// Define messages variable.
	messages := []models.Message{}

	// Define query string.
	query := `SELECT * FROM messages`

	// Send query to database.
	err := q.Get(&messages, query)
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
	query := `INSERT INTO messages VALUES ($1, $2, $3, $4, $5, $6)`

	// Send query to database.
	_, err := q.Exec(query, m.UUID, m.Timestamp, m.Message, m.Author, m.Likes)
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
