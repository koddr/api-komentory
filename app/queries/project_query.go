package queries

import (
	"Komentory/api/app/models"
	"Komentory/api/pkg/repository"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ProjectQueries struct for queries from Project model.
type ProjectQueries struct {
	*sqlx.DB
}

// GetProjects method for getting all project.
func (q *ProjectQueries) GetProjects() ([]models.Project, error) {
	// Define project variable.
	project := []models.Project{}

	// Define query string.
	query := `
	SELECT * 
	FROM projects
	ORDER BY created_at DESC
	`

	// Send query to database.
	err := q.Select(&project, query)
	if err != nil {
		// Return empty object and error.
		return project, err
	}

	// Return query result.
	return project, nil
}

// GetProjectsByUserID method for getting all project by given user ID.
func (q *ProjectQueries) GetProjectsByUserID(user_id uuid.UUID) ([]models.Project, int, error) {
	// Define project variable.
	project := []models.Project{}

	// Define query string.
	query := `
	SELECT * 
	FROM projects 
	WHERE (user_id = $1::uuid AND project_status = 1) 
	ORDER BY created_at DESC
	`

	// Send query to database.
	err := q.Select(&project, query, user_id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return project, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return project, fiber.StatusNotFound, fmt.Errorf(repository.NotFoundProjectsByUser)
	default:
		// Return empty object and 400 error.
		return project, fiber.StatusBadRequest, err
	}
}

// GetProjectByID method for getting one project by given ID.
func (q *ProjectQueries) GetProjectByID(id uuid.UUID) (models.Project, int, error) {
	// Define project variable.
	project := models.Project{}

	// Define query string.
	query := `
	SELECT * 
	FROM projects 
	WHERE id = $1::uuid
	LIMIT 1
	`

	// Send query to database.
	err := q.Get(&project, query, id)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return project, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return project, fiber.StatusNotFound, fmt.Errorf(repository.NotFoundProjectWithID)
	default:
		// Return empty object and 400 error.
		return project, fiber.StatusBadRequest, err
	}
}

// GetProjectByAlias method for getting one project by given alias.
func (q *ProjectQueries) GetProjectByAlias(alias string) (models.Project, int, error) {
	// Define project variable.
	project := models.Project{}

	// Define query string.
	query := `
	SELECT * 
	FROM projects 
	WHERE alias = $1::varchar
	LIMIT 1
	`

	// Send query to database.
	err := q.Get(&project, query, alias)

	// Get query result.
	switch err {
	case nil:
		// Return object and 200 OK.
		return project, fiber.StatusOK, nil
	case sql.ErrNoRows:
		// Return empty object and 404 error.
		return project, fiber.StatusNotFound, fmt.Errorf(repository.NotFoundProjectWithAlias)
	default:
		// Return empty object and 400 error.
		return project, fiber.StatusBadRequest, err
	}
}

// CreateProject method for creating project by given Project object.
func (q *ProjectQueries) CreateProject(p *models.Project) error {
	// Define query string.
	query := `
	INSERT INTO projects 
	VALUES (
		$1::uuid, $2::timestamp, $3::timestamp, 
		$4::uuid, $5::varchar, $6::int, 
		$7::jsonb
	)
	`

	// Send query to database.
	_, err := q.Exec(
		query,
		p.ID, p.CreatedAt, p.UpdatedAt,
		p.UserID, p.Alias, p.ProjectStatus,
		p.ProjectAttrs,
	)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateProject method for updating project by given Project object.
func (q *ProjectQueries) UpdateProject(id uuid.UUID, p *models.Project) error {
	// Define query string.
	query := `
	UPDATE projects 
	SET updated_at = $2::timestamp, project_status = $3::int, project_attrs = $4::jsonb 
	WHERE id = $1::uuid
	`

	// Send query to database.
	_, err := q.Exec(query, id, p.UpdatedAt, p.ProjectStatus, p.ProjectAttrs)
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
