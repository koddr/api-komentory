package queries

import (
	"Komentory/api/app/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AnswerQueries struct for queries from Answer model.
type AnswerQueries struct {
	*sqlx.DB
}

// GetAnswerByID method for getting one answer by given ID.
// NOTE: This method is using ONLY for cheking if given object is exists.
func (q *AnswerQueries) GetAnswerByID(id uuid.UUID) (models.Answer, int, error) {
	// Define project variable.
	task := models.Answer{}

	// Define query string.
	query := `
	SELECT *
	FROM
		answers
	WHERE
		id = $1::uuid
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
		return task, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return task, fiber.StatusBadRequest, err
	}
}

// GetAnswerByAlias method for getting one answer by given alias.
func (q *AnswerQueries) GetAnswerByAlias(alias string) (models.GetAnswer, int, error) {
	// Define project variable.
	task := models.GetAnswer{}

	// Define query string.
	query := `
	SELECT
		a.*,
		u.username,
		u.user_attrs
	FROM
		answers AS a
		LEFT JOIN users AS u ON u.id = a.user_id
	WHERE 
		a.alias = $1::varchar
	LIMIT 1
	`

	// Send query to database.
	err := q.Get(&task, query, alias)

	// Get quey result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return task, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return task, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return task, fiber.StatusBadRequest, err
	}
}

// GetAnswersByProjectID method for getting all answers for given project.
func (q *AnswerQueries) GetAnswersByProjectID(project_id uuid.UUID) ([]models.GetAnswers, int, error) {
	// Define project variable.
	answers := []models.GetAnswers{}

	// Define query string.
	query := `
	SELECT
		a.id,
		a.created_at,
		a.updated_at,
		a.alias,
		a.answer_attrs,
		u.username,
		u.user_attrs
	FROM
		answers AS a
		LEFT JOIN users AS u ON u.id = a.user_id
	WHERE
		a.project_id = $1::uuid
		AND a.answer_status = 1
	ORDER BY
		a.created_at DESC
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
		return answers, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return answers, fiber.StatusBadRequest, err
	}
}

// GetAnswersByTaskID method for getting all answers for given task.
func (q *AnswerQueries) GetAnswersByTaskID(task_id uuid.UUID) ([]models.GetAnswers, int, error) {
	// Define answer variable.
	answers := []models.GetAnswers{}

	// Define query string.
	query := `
	SELECT
		a.id,
		a.created_at,
		a.updated_at,
		a.alias,
		a.answer_attrs,
		u.username,
		u.user_attrs
	FROM
		answers AS a
		LEFT JOIN users AS u ON a.user_id = u.id
	WHERE
		a.task_id = $1::uuid
		AND a.answer_status = 1
	ORDER BY
		a.created_at DESC
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
		return answers, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return answers, fiber.StatusBadRequest, err
	}
}

// CreateAnswer method for creating answer by given Answer object.
func (q *AnswerQueries) CreateNewAnswer(a *models.Answer) error {
	// Define query string.
	query := `
	INSERT INTO answers 
	VALUES (
		$1::uuid, $2::timestamp, $3::timestamp, 
		$4::uuid, $5::uuid, $6::uuid, 
		$7::varchar, $8::int, $9::jsonb
	)
	`

	// Send query to database.
	_, err := q.Exec(
		query,
		a.ID, a.CreatedAt, a.UpdatedAt,
		a.UserID, a.ProjectID, a.TaskID,
		a.Alias, a.AnswerStatus, a.AnswerAttrs,
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
	UPDATE
		answers
	SET
		updated_at = $2::timestamp,
		answer_status = $3::int,
		answer_attrs = $4::jsonb
	WHERE
		id = $1::uuid
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
