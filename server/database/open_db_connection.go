package database

import "github.com/chayapol-p/cloud-midterm-server/queries"

// Queries struct for collect all app queries.
type Queries struct {
	*queries.MessageQueries // load queries from Book model
	*queries.UpdatedMessageQueries
}

// OpenDBConnection func for opening database connection.
func OpenDBConnection() (*Queries, error) {
	// Define a new PostgreSQL connection.
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		// Set queries from models:
		MessageQueries:        &queries.MessageQueries{DB: db}, // from Book model
		UpdatedMessageQueries: &queries.UpdatedMessageQueries{DB: db},
	}, nil
}
