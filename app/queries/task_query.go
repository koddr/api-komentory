package queries

import (
	"Komentory/api/app/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TaskQueries struct for queries from Task model.
type TaskQueries struct {
	*sqlx.DB
}

// GetTaskByID method for getting one project by given ID.
func (q *TaskQueries) GetTaskByID(id uuid.UUID) (models.Task, int, error) {
	// Define project variable.
	task := models.Task{}

	// Define query string.
	query := `
	SELECT
		t.*,
		COUNT(a.id) AS answers_count
	FROM
		tasks AS t
		LEFT JOIN answers AS a ON t.id = a.task_id
	WHERE
		t.id = $1::uuid
	GROUP BY
		t.id
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

// GetTaskByAlias method for getting one task by given alias.
func (q *TaskQueries) GetTaskByAlias(alias string) (models.Task, int, error) {
	// Define project variable.
	task := models.Task{}

	// Define query string.
	query := `
	SELECT
		t.*,
		COUNT(a.id) AS answers_count
	FROM
		tasks AS t
		LEFT JOIN answers AS a ON t.id = a.task_id
	WHERE
		t.alias = $1::varchar
	GROUP BY
		t.id
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

// GetTasksByProjectID method for getting all tasks for given project.
func (q *TaskQueries) GetTasksByProjectID(project_id uuid.UUID) ([]models.TaskList, int, error) {
	// Define project variable.
	tasks := []models.TaskList{}

	// Define query string.
	query := `
	SELECT
		t.id,
		t.created_at,
		t.updated_at,
		t.alias,
		t.task_attrs,
		COUNT(a.id) AS answers_count
	FROM
		tasks AS t
		LEFT JOIN answers AS a ON t.id = a.task_id
	WHERE
		t.project_id = $1::uuid
		AND t.task_status = 1
	GROUP BY
		t.id
	ORDER BY
		t.created_at DESC
	`

	// Send query to database.
	err := q.Select(&tasks, query, project_id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return tasks, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return tasks, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return tasks, fiber.StatusBadRequest, err
	}
}

// CreateNewTask method for creating a new task.
func (q *TaskQueries) CreateNewTask(t *models.Task) error {
	// Define query string.
	query := `
	INSERT INTO tasks
	VALUES (
		$1::uuid, $2::timestamp, $3::timestamp, 
		$4::uuid, $5::uuid, $6::varchar, 
		$7::int, $7::jsonb
	)
	`

	// Send query to database.
	_, err := q.Exec(
		query,
		t.ID, t.CreatedAt, t.UpdatedAt,
		t.UserID, t.ProjectID, t.Alias,
		t.TaskStatus, t.TaskAttrs,
	)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateTask method for updating task by given Task object.
func (q *TaskQueries) UpdateTask(id uuid.UUID, t *models.Task) error {
	// Define query string.
	query := `
	UPDATE
		tasks
	SET
		updated_at = $2::timestamp,
		task_status = $3::int,
		task_attrs = $4::jsonb
	WHERE
		id = $1::uuid
	`

	// Send query to database.
	_, err := q.Exec(query, id, t.UpdatedAt, t.TaskStatus, t.TaskAttrs)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteTask method for delete task by given ID.
func (q *TaskQueries) DeleteTask(id uuid.UUID) error {
	// Define query string.
	query := `
	DELETE FROM tasks
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
