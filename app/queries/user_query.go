package queries

import (
	"Komentory/api/app/models"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

// GetUserByEmail query for getting one User by given Email.
func (q *UserQueries) GetUserByEmail(email string) (models.User, int, error) {
	// Define User variable.
	user := models.User{}

	// Define query string.
	query := `
	SELECT id, email 
	FROM users 
	WHERE email = $1::varchar
	`

	// Send query to database.
	err := q.Get(&user, query, email)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return user, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return user, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return user, fiber.StatusBadRequest, err
	}
}

// UpdateUserSettings method for updating user settings by given user ID.
func (q *UserQueries) UpdateUserSettings(id uuid.UUID, u *models.UserSettings) error {
	// Define query string.
	query := `
	UPDATE users 
	SET updated_at = $2::timestamp, user_settings = $3::jsonb 
	WHERE id = $1::uuid
	`

	// Send query to database.
	_, err := q.Exec(query, id, time.Now(), u)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
