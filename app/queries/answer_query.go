package queries

import (
	"Komentory/api/app/models"
	"database/sql"
	"fmt"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AnswerQueries struct for queries from Answer model.
type AnswerQueries struct {
	*sqlx.DB
}

// GetAnswerByID method for getting one answer by given ID.
func (q *AnswerQueries) GetAnswerByID(id uuid.UUID) (models.Answer, int, error) {
	// Define project variable.
	task := models.Answer{}

	// Define query string.
	query := `
	SELECT * 
	FROM answers 
	WHERE id = $1::uuid
	LIMIT 1
	`

	// Send query to database.
	err := q.Get(&task, query, id)

	// Get quey result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return task, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return task, fiber.StatusNotFound, fmt.Errorf(utilities.GenerateErrorMessage(404, "answer", "id"))
	default:
		// Return empty object and 400 error.
		return task, fiber.StatusBadRequest, err
	}
}

// GetAnswersByProjectID method for getting all answers for given project.
func (q *AnswerQueries) GetAnswersByProjectID(project_id uuid.UUID) ([]models.Answer, int, error) {
	// Define project variable.
	answers := []models.Answer{}

	// Define query string.
	query := `
	SELECT * 
	FROM answers 
	WHERE (project_id = $1::uuid AND answer_status = 1) 
	ORDER BY created_at DESC
	`

	// Send query to database.
	err := q.Select(&answers, query, project_id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return answers, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return answers, fiber.StatusNotFound, fmt.Errorf(utilities.GenerateErrorMessage(404, "answer", "project_id"))
	default:
		// Return empty object and 400 error.
		return answers, fiber.StatusBadRequest, err
	}
}

// GetAnswersByTaskID method for getting all answers for given task.
func (q *AnswerQueries) GetAnswersByTaskID(task_id uuid.UUID) ([]models.AnswerList, int, error) {
	// Define project variable.
	answers := []models.AnswerList{}

	// Define query string.
	query := `
	SELECT id, created_at, updated_at, answer_attrs 
	FROM answers 
	WHERE (task_id = $1::uuid AND answer_status = 1) 
	ORDER BY created_at DESC
	`

	// Send query to database.
	err := q.Select(&answers, query, task_id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return answers, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return answers, fiber.StatusNotFound, fmt.Errorf(utilities.GenerateErrorMessage(404, "answer", "task_id"))
	default:
		// Return empty object and 400 error.
		return answers, fiber.StatusBadRequest, err
	}
}

// CreateAnswer method for creating answer by given Answer object.
func (q *AnswerQueries) CreateAnswer(a *models.Answer) error {
	// Define query string.
	query := `
	INSERT INTO answers 
	VALUES (
		$1::uuid, $2::timestamp, $3::timestamp, 
		$4::uuid, $5::uuid, $6::uuid, 
		$7::int, $8::jsonb
	)
	`

	// Send query to database.
	_, err := q.Exec(
		query,
		a.ID, a.CreatedAt, a.UpdatedAt,
		a.UserID, a.ProjectID, a.TaskID,
		a.AnswerStatus, a.AnswerAttrs,
	)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateAnswer method for updating answer by given Answer object.
func (q *AnswerQueries) UpdateAnswer(id uuid.UUID, a *models.Answer) error {
	// Define query string.
	query := `
	UPDATE answers 
	SET updated_at = $2::timestamp, answer_status = $3::int, answer_attrs = $4::jsonb 
	WHERE id = $1::uuid
	`

	// Send query to database.
	_, err := q.Exec(query, id, a.UpdatedAt, a.AnswerStatus, a.AnswerAttrs)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteAnswer method for delete task by given ID.
func (q *AnswerQueries) DeleteAnswer(id uuid.UUID) error {
	// Define query string.
	query := `
	DELETE FROM answers 
	WHERE id = $1::uuid
	`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
