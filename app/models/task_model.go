package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ---
// Structures to describing task model.
// ---

// Task struct to describe task object.
type Task struct {
	ID         uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	UserID     uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
	ProjectID  uuid.UUID `db:"project_id" json:"project_id" validate:"required,uuid"`
	TaskStatus int       `db:"task_status" json:"task_status" validate:"int"`
	TaskAttrs  TaskAttrs `db:"task_attrs" json:"task_attrs" validate:"required,dive"`
}

// TaskAttrs struct to describe task attributes.
type TaskAttrs struct {
	Name        string     `json:"name" validate:"required,lte=255"`
	Description string     `json:"description" validate:"required"`
	Steps       []taskStep `json:"steps" validate:"required,dive"`
	Documents   []string   `json:"documents"`
	Images      []string   `json:"images"`
	Links       []string   `json:"links"`
}

// ---
// Structures to creating a new task.
// ---

// CreateNewTask struct to describe create a new task process.
type CreateNewTask struct {
	ProjectID  uuid.UUID `json:"project_id" validate:"required,uuid"`
	TaskStatus int       `json:"task_status" validate:"int"`
	TaskAttrs  TaskAttrs `json:"task_attrs" validate:"required,dive"`
}

// ---
// Structures to updating one task.
// ---

// UpdateTask struct to describe update process of the given task.
type UpdateTask struct {
	ID         uuid.UUID `json:"id" validate:"required,uuid"`
	TaskStatus int       `json:"task_status" validate:"int"`
	TaskAttrs  TaskAttrs `json:"task_attrs" validate:"required,dive"`
}

// ---
// Structures to deleting one task.
// ---

// DeleteTask struct to describe delete process of the given task.
type DeleteTask struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

// ---
// Structures to getting only one task.
// ---

// GetTasks struct to describe getting one task.
type GetTask struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`
	Status    int       `db:"task_status" json:"status"`
	Attrs     TaskAttrs `db:"task_attrs" json:"attrs"`

	// Fields for JOIN tables:
	AnswersCount int `db:"answers_count" json:"answers_count"`
}

// ---
// Structures to getting many tasks.
// ---

// GetTasks struct to describe getting tasks list.
type GetTasks struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Attrs     TaskAttrs `db:"task_attrs" json:"attrs"`

	// Fields for JOIN tables:
	AnswersCount int `db:"answers_count" json:"answers_count"`
}

// ---
// Private structures to building better model JSON output.
// ---

// taskStep (private) struct to describe step of given task.
type taskStep struct {
	Position    int    `json:"position" validate:"required,int"`
	Description string `json:"description" validate:"required"`
}

// ---
// This methods simply returns the JSON-encoded representation of the struct.
// ---

// Value make the TaskAttrs struct implement the driver.Valuer interface.
func (t *TaskAttrs) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// ---
// This methods simply decodes a JSON-encoded value into the struct fields.
// ---

// Scan make the TaskAttrs struct implement the sql.Scanner interface.
func (t *TaskAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(j, &t)
}
