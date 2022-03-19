package queries

import (
	"github.com/chayapol-p/cloud-midterm-server/models"
	"github.com/jmoiron/sqlx"
)

// MessageQueries struct for queries from Message model.
type UpdatedMessageQueries struct {
	*sqlx.DB
}

// GetMessages method for getting all messages.
func (q *UpdatedMessageQueries) GetUpdates(time string) ([]models.OutputUpdatedMessage, error) {
	// Define messages variable.
	updates := []models.OutputUpdatedMessage{}

	// Define query string.
	query := `SELECT uuid, author, message, likes, is_deleted FROM updated_messages Where timestamp > $1`

	// Send query to database.
	err := q.Select(&updates, query, time)
	if err != nil {
		// Return empty object and error.
		return updates, err
	}

	// Return query result.
	return updates, nil
}

// // GetMessage method for getting one message by given ID.
func (q *UpdatedMessageQueries) GetUpdate(uuid string) (models.UpdatedMessage, error) {
	// Define message variable.
	update := models.UpdatedMessage{}

	// Define query string.
	query := `SELECT * FROM updated_messages WHERE uuid = $1`

	// Send query to database.
	err := q.Get(&update, query, uuid)
	if err != nil {
		// Return empty object and error.
		return update, err
	}

	// Return query result.
	return update, nil
}

// CreateMessage method for creating message by given Message object.
func (q *UpdatedMessageQueries) CreateUpdate(m *models.UpdatedMessage) error {
	// Define query string.
	query := `INSERT INTO updated_messages VALUES ($1, $2, $3, $4, $5, $6)`

	// Send query to database.
	_, err := q.Exec(query, m.UUID, m.Timestamp, m.Author, m.Message, m.Likes, m.IsDeleted)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateMessage method for updating message by given Message object.
func (q *UpdatedMessageQueries) UpdateUpdate(id string, m *models.UpdatedMessage) error {
	// Define query string.
	query := `UPDATE updated_messages SET timestamp = $2, message = $3, author = $4, likes = $5, is_deleted = $6 WHERE uuid = $1`

	// Send query to database.
	_, err := q.Exec(query, id, m.Timestamp, m.Message, m.Author, m.Likes, m.IsDeleted)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteMessage method for delete message by given ID.
func (q *UpdatedMessageQueries) DeleteUpdate(id string) error {
	// Define query string.
	query := `DELETE FROM updated_messages WHERE uuid = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
