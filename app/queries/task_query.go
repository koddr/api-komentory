package queries

import (
	"Komentory/api/app/models"
	"Komentory/api/platform/embed_files"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TaskQueries struct for queries from Task model.
type TaskQueries struct {
	*sqlx.DB
}

// FindTaskByID method for find one task by given ID.
func (q *TaskQueries) FindTaskByID(task_id uuid.UUID) (models.Task, int, error) {
	// Define task variable.
	task := models.Task{}

	// Define query string.
	query := `
	SELECT
		id,
		user_id
	FROM
		tasks
	WHERE
		id = $1::uuid
	LIMIT 1
	`

	// Send query to database.
	err := q.Get(&task, query, task_id)

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

// CreateNewTask method for creating a new task.
func (q *TaskQueries) CreateNewTask(t *models.Task) error {
	// Define query string.
	query := `
	INSERT INTO tasks
	VALUES (
		$1::uuid, $2::timestamp, $3::timestamp, 
		$4::uuid, $5::uuid, $6::int, 
		$7::jsonb
	)
	`

	// Send query to database.
	_, err := q.Exec(
		query,
		t.ID, t.CreatedAt, t.UpdatedAt,
		t.UserID, t.ProjectID, t.TaskStatus,
		t.TaskAttrs,
	)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateTask method for updating task by given Task object.
func (q *TaskQueries) UpdateTask(id uuid.UUID, t *models.UpdateTask) error {
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
	_, err := q.Exec(query, id, time.Now(), t.TaskStatus, t.TaskAttrs)
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

// GetTaskByID method for getting one project by given ID.
func (q *TaskQueries) GetTaskByID(task_id uuid.UUID) (models.GetTask, int, error) {
	// Define project variable.
	task := models.GetTask{}

	// Define query string.
	query := embed_files.SQLQueryGetOneTaskByID

	// Send query to database.
	err := q.Get(&task, query, task_id)

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
func (q *TaskQueries) GetTasksByProjectID(project_id uuid.UUID) ([]models.GetTasks, int, error) {
	// Define project variable.
	tasks := []models.GetTasks{}

	// Define query string.
	query := embed_files.SQLQueryGetManyTasksByProjectID

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
