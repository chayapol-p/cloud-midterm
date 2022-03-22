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
func (q *UpdatedMessageQueries) GetUpdates(time string, offset string, limit string) ([]models.OutputUpdatedMessage, error) {
	// Define messages variable.
	updates := []models.OutputUpdatedMessage{}

	// Define query string.
	// query := `SELECT
	// 	updated_messages.uuid,
	// 	updated_messages.author,
	// 	updated_messages.message,
	// 	updated_messages.likes
	// FROM
	// 	( select * from
	// 		messages
	// 		where timestamp > $1
	// 	) as messages
	// RIGHT JOIN updated_messages
	// ON messages.uuid = updated_messages.uuid
	// WHERE updated_messages.timestamp > $1
	// 	and messages.uuid is NULL
	// 	and not updated_messages.is_deleted
	// offset $2 limit $3;`

	query := `SELECT uuid, author, message, likes FROM updated_messages WHERE not is_deleted and timestamp > $1 offset $2 limit $3;`
	// Send query to database.
	err := q.Select(&updates, query, time, offset, limit)
	if err != nil {
		// Return empty object and error.
		return updates, err
	}

	// Return query result.
	return updates, nil
}

// GetMessages method for getting all messages.
func (q *UpdatedMessageQueries) GetDeletes(time string, offset string, limit string) ([]string, error) {
	// Define messages variable.
	deletes := []string{}

	// Define query string.
	// query := `SELECT
	// 	updated_messages.uuid
	// FROM
	// 	( select * from
	// 		messages
	// 		where timestamp > $1
	// 	) as messages
	// RIGHT JOIN updated_messages
	// ON messages.uuid = updated_messages.uuid
	// WHERE updated_messages.timestamp > $1
	// 	and messages.uuid is NULL
	// 	and updated_messages.is_deleted
	// offset $2 limit $3;`

	query := `SELECT uuid FROM updated_messages WHERE is_deleted and timestamp > $1 offset $2 limit $3;`

	// Send query to database.
	err := q.Select(&deletes, query, time, offset, limit)
	if err != nil {
		// Return empty object and error.
		return deletes, err
	}

	// Return query result.
	return deletes, nil
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
