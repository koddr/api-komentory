package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ---
// Structures to describing project model.
// ---

// Project struct to describe project object.
type Project struct {
	ID            uuid.UUID    `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
	UserID        uuid.UUID    `db:"user_id" json:"user_id" validate:"required,uuid"`
	ProjectStatus int          `db:"project_status" json:"project_status" validate:"int"`
	ProjectAttrs  ProjectAttrs `db:"project_attrs" json:"project_attrs" validate:"required,dive"`
}

// ProjectAttrs struct to describe project attributes.
type ProjectAttrs struct {
	Title       string   `json:"title" validate:"required,lte=255"`
	Description string   `json:"description" validate:"required"`
	Category    string   `json:"category" validate:"required"`
	WebsiteURL  string   `json:"website_url"`
	Picture     string   `json:"picture"`
	Tags        []string `json:"tags"`
}

// ---
// Structures to creating a new project.
// ---

// CreateNewProject struct to describe create a new project process.
type CreateNewProject struct {
	ProjectStatus int          `db:"project_status" json:"project_status" validate:"int"`
	ProjectAttrs  ProjectAttrs `db:"project_attrs" json:"project_attrs" validate:"required,dive"`
}

// ---
// Structures to updating one project.
// ---

// UpdateProject struct to describe update process of the given project.
type UpdateProject struct {
	ID            uuid.UUID    `db:"id" json:"id" validate:"required,uuid"`
	ProjectStatus int          `db:"project_status" json:"project_status" validate:"int"`
	ProjectAttrs  ProjectAttrs `db:"project_attrs" json:"project_attrs" validate:"required,dive"`
}

// ---
// Structures to deleting one project.
// ---

// DeleteProject struct to describe delete process of the given project.
type DeleteProject struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

// ---
// Structures to getting only one project.
// ---

// GetProject struct to describe getting one project.
type GetProject struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	Status    int          `db:"project_status" json:"status"`
	Attrs     ProjectAttrs `db:"project_attrs" json:"attrs"`

	// Fields for JOIN tables:
	Author     authorAttrs  `db:"author" json:"author"`
	TasksCount int          `db:"tasks_count" json:"tasks_count"`
	Tasks      projectTasks `db:"tasks" json:"tasks"`
}

// GetProjects struct to describe getting list of projects.
type GetProjects struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	Attrs     ProjectAttrs `db:"project_attrs" json:"attrs"`

	// Fields for JOIN tables:
	Author     authorAttrs `db:"author" json:"author"`
	TasksCount int         `db:"tasks_count" json:"tasks_count"`
}

// ---
// Private structures to building better model JSON output.
// ---

// authorAttrs (private) struct to describe author of given project.
type authorAttrs struct {
	ID        uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Picture   string    `json:"picture"`
}

// projectTasks (private) struct to describe getting list of tasks for a project.
type projectTasks []*getProjectTasks

// getProjectTasks (private) struct to describe getting tasks list for given project.
type getProjectTasks struct {
	ID          uuid.UUID `json:"id"`
	Status      int       `json:"status"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StepsCount  int       `json:"steps_count"`
}

// ---
// This methods simply returns the JSON-encoded representation of the struct.
// ---

// Value make the ProjectAttrs struct implement the driver.Valuer interface.
func (p ProjectAttrs) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// ---
// This methods simply decodes a JSON-encoded value into the struct fields.
// ---

// Scan make the ProjectAttrs struct implement the sql.Scanner interface.
func (p *ProjectAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(j, &p)
}

// Scan make the authorAttrs (private) struct implement the sql.Scanner interface.
func (t *authorAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(j, &t)
}

// Scan make the projectTasks (private) struct implement the sql.Scanner interface.
func (t *projectTasks) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(j, &t)
}
