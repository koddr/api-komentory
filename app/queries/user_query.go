package queries

import (
	"Komentory/api/app/models"
	"Komentory/api/pkg/repository"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	*sqlx.DB
}

// GetUserByID query for getting one User by given ID.
func (q *UserQueries) GetUserByID(id uuid.UUID) (models.User, int, error) {
	// Define User variable.
	user := models.User{}

	// Define query string.
	query := `
	SELECT * 
	FROM users 
	WHERE id = $1::uuid
	`

	// Send query to database.
	err := q.Get(&user, query, id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return user, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return user, fiber.StatusNotFound, fmt.Errorf(repository.NotFoundUserWithID)
	default:
		// Return empty object and 400 error.
		return user, fiber.StatusBadRequest, err
	}
}

// GetUserByEmail query for getting one User by given Email.
func (q *UserQueries) GetUserByEmail(email string) (models.User, int, error) {
	// Define User variable.
	user := models.User{}

	// Define query string.
	query := `
	SELECT * 
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
		return user, fiber.StatusNotFound, fmt.Errorf(repository.NotFoundUserWithEmail)
	default:
		// Return empty object and 400 error.
		return user, fiber.StatusBadRequest, err
	}
}

// GetUserByUsername query for getting one User by given username.
func (q *UserQueries) GetUserByUsername(username string) (models.User, int, error) {
	// Define User variable.
	user := models.User{}

	// Define query string.
	query := `
	SELECT * 
	FROM users 
	WHERE username = $1::varchar
	`

	// Send query to database.
	err := q.Get(&user, query, username)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return user, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return user, fiber.StatusNotFound, fmt.Errorf(repository.NotFoundUserWithUsername)
	default:
		// Return empty object and 400 error.
		return user, fiber.StatusBadRequest, err
	}
}

// CreateUser query for creating a new user by given email and password hash.
func (q *UserQueries) CreateUser(u *models.User) error {
	// Define query string.
	query := `
	INSERT INTO users 
	VALUES (
		$1::uuid, $2::timestamp, $3::timestamp, 
		$4::varchar, $5::varchar, $6::varchar, 
		$7::int, $8::varchar, $9::jsonb
	)
	`

	// Send query to database.
	_, err := q.Exec(
		query,
		u.ID, u.CreatedAt, u.UpdatedAt,
		u.Email, u.PasswordHash, u.Username,
		u.UserStatus, u.UserRole, u.UserAttrs,
	)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateUserPassword method for updating user password by given user ID.
func (q *UserQueries) UpdateUserPassword(id uuid.UUID, p string) error {
	// Define query string.
	query := `
	UPDATE users 
	SET updated_at = $2::timestamp, password_hash = $3::varchar 
	WHERE id = $1::uuid
	`

	// Send query to database.
	_, err := q.Exec(query, id, time.Now(), p)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateUserAttrs method for updating user attrs by given user ID.
func (q *UserQueries) UpdateUserAttrs(id uuid.UUID, u *models.UserAttrs) error {
	// Define query string.
	query := `
	UPDATE users 
	SET updated_at = $2::timestamp, user_attrs = $3::jsonb 
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
