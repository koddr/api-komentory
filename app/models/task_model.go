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
	TaskStatus int       `db:"task_status" json:"task_status" validate:"int"`
	TaskAttrs  TaskAttrs `db:"task_attrs" json:"task_attrs" validate:"required,dive"`
}

// TaskList struct to describe task list object.
type TaskList struct {
	ID        uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	TaskAttrs TaskAttrs `db:"task_attrs" json:"task_attrs" validate:"required,dive"`
}

// TaskAttrs struct to describe task attributes.
type TaskAttrs struct {
	Title       string `json:"title" validate:"required,lte=255"`
	Description string `json:"description" validate:"required"`
	Picture     string `json:"picture"`
	URL         string `json:"url"`
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
