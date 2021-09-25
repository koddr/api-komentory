package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Task struct to describe task object.
type Task struct {
	ID         uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	UserID     uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
	ProjectID  uuid.UUID `db:"project_id" json:"project_id" validate:"required,uuid"`
	Alias      string    `db:"alias" json:"alias" validate:"required,lte=16"`
	TaskStatus int       `db:"task_status" json:"task_status" validate:"int"`
	TaskAttrs  TaskAttrs `db:"task_attrs" json:"task_attrs" validate:"required,dive"`
}

// TaskAttrs struct to describe task attributes.
type TaskAttrs struct {
	Name        string     `json:"name" validate:"required,lte=255"`
	Description string     `json:"description" validate:"required"`
	Steps       []TaskStep `json:"steps" validate:"required,dive"`
	Documents   []string   `json:"documents"`
	Images      []string   `json:"images"`
	Links       []string   `json:"links"`
}

// TaskStep struct to describe task step object.
type TaskStep struct {
	Position    int    `json:"position" validate:"required,int"`
	Description string `json:"description" validate:"required"`
}

// GetTasks struct to describe getting tasks list.
type GetTask struct {
	ID         uuid.UUID `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	ProjectID  uuid.UUID `db:"project_id" json:"project_id"`
	Alias      string    `db:"alias" json:"alias"`
	TaskStatus int       `db:"task_status" json:"task_status"`
	TaskAttrs  TaskAttrs `db:"task_attrs" json:"task_attrs"`

	// Fields for JOIN tables:
	AnswersCount int `db:"answers_count" json:"answers_count"`
}

// GetTasks struct to describe getting tasks list.
type GetTasks struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Alias     string    `db:"alias" json:"alias"`
	TaskAttrs TaskAttrs `db:"task_attrs" json:"task_attrs"`

	// Fields for JOIN tables:
	AnswersCount int `db:"answers_count" json:"answers_count"`
}

// Value make the TaskAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (t TaskAttrs) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan make the TaskAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (t *TaskAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(j, &t)
}
