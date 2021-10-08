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

// ProjectQueries struct for queries from Project model.
type ProjectQueries struct {
	*sqlx.DB
}

// FindProjectByID method for find one project by given ID.
func (q *ProjectQueries) FindProjectByID(project_id uuid.UUID) (models.Project, int, error) {
	// Define project variable.
	project := models.Project{}

	// Define query string.
	query := `
	SELECT 
		id,
		user_id
	FROM
		projects
	WHERE
		id = $1::uuid
	LIMIT 1
	`

	// Send query to database.
	err := q.Get(&project, query, project_id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return project, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return project, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return project, fiber.StatusBadRequest, err
	}
}

// CreateProject method for creating project by given Project object.
func (q *ProjectQueries) CreateNewProject(p *models.Project) error {
	// Define query string.
	query := `
	INSERT INTO projects
	VALUES (
		$1::uuid, $2::timestamp, $3::timestamp, 
		$4::uuid, $5::int, $6::jsonb
	)
	`

	// Send query to database.
	_, err := q.Exec(
		query,
		p.ID, time.Now(), p.UpdatedAt,
		p.UserID, p.ProjectStatus, p.ProjectAttrs,
	)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateProject method for updating project by given Project object.
func (q *ProjectQueries) UpdateProject(id uuid.UUID, p *models.UpdateProject) error {
	// Define query string.
	query := `
	UPDATE
		projects
	SET
		updated_at = $2::timestamp,
		project_status = $3::int,
		project_attrs = $4::jsonb
	WHERE
		id = $1::uuid
	`

	// Send query to database.
	_, err := q.Exec(query, id, time.Now(), p.ProjectStatus, p.ProjectAttrs)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteProject method for delete project by given ID.
func (q *ProjectQueries) DeleteProject(id uuid.UUID) error {
	// Define query string.
	query := `
	DELETE FROM projects
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

// GetProjects method for getting all projects.
func (q *ProjectQueries) GetProjects() ([]models.GetProjects, int, error) {
	// Define project variable.
	projects := []models.GetProjects{}

	// Define query string.
	query := embed_files.SQLQueryGetManyProjects

	// Send query to database.
	err := q.Select(&projects, query)

	// Return query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return projects, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return projects, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return projects, fiber.StatusBadRequest, err
	}
}

// GetProjectsByUserID method for getting all project by given user ID.
func (q *ProjectQueries) GetProjectsByUserID(user_id uuid.UUID) ([]models.GetProjects, int, error) {
	// Define project variable.
	projects := []models.GetProjects{}

	// Define query string.
	query := embed_files.SQLQueryGetManyProjectsByUserID

	// Send query to database.
	err := q.Select(&projects, query, user_id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return projects, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return projects, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return projects, fiber.StatusBadRequest, err
	}
}

// GetProjectByAlias method for getting one project by given alias.
func (q *ProjectQueries) GetProjectByID(project_id uuid.UUID) (models.GetProject, int, error) {
	// Define project variable.
	project := models.GetProject{}

	// Define query string.
	query := embed_files.SQLQueryGetOneProjectByID

	// Send query to database.
	err := q.Get(&project, query, project_id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return project, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return project, fiber.StatusNotFound, err
	default:
		// Return empty object and 400 error.
		return project, fiber.StatusBadRequest, err
	}
}
