package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Answer struct to describe answer object.
type Answer struct {
	ID           uuid.UUID   `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
	UserID       uuid.UUID   `db:"user_id" json:"user_id" validate:"required,uuid"`
	ProjectID    uuid.UUID   `db:"project_id" json:"project_id" validate:"required,uuid"`
	TaskID       uuid.UUID   `db:"task_id" json:"task_id" validate:"required,uuid"`
	AnswerStatus int         `db:"answer_status" json:"answer_status" validate:"int"`
	AnswerAttrs  AnswerAttrs `db:"answer_attrs" json:"answer_attrs" validate:"required,dive"`
}

// AnswerList struct to describe answer list object.
type AnswerList struct {
	ID          uuid.UUID   `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt   time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updated_at"`
	AnswerAttrs AnswerAttrs `db:"answer_attrs" json:"answer_attrs" validate:"required,dive"`
}

// AnswerAttrs struct to describe answer attributes.
type AnswerAttrs struct {
	Title       string `json:"title" validate:"required,lte=255"`
	Description string `json:"description" validate:"required"`
	Picture     string `json:"picture"`
	URL         string `json:"url"`
}

// Value make the AnswerAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (a AnswerAttrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan make the AnswerAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (a *AnswerAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &a)
}
