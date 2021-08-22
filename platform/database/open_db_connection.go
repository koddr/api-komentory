package database

import "Komentory/api/app/queries"

// Queries struct for collect all app queries.
type Queries struct {
	*queries.UserQueries    // load queries from User model
	*queries.ProjectQueries // load queries from Project model
	*queries.TaskQueries    // load queries from Task model
	*queries.AnswerQueries  // load queries from Answer model
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
		UserQueries:    &queries.UserQueries{DB: db},    // from User model
		ProjectQueries: &queries.ProjectQueries{DB: db}, // from Project model
		TaskQueries:    &queries.TaskQueries{DB: db},    // from Task model
		AnswerQueries:  &queries.AnswerQueries{DB: db},  // from Answer model
	}, nil
}
