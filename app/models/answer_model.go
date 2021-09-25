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
	Alias        string      `db:"alias" json:"alias" validate:"required,lte=16"`
	AnswerStatus int         `db:"answer_status" json:"answer_status" validate:"int"`
	AnswerAttrs  AnswerAttrs `db:"answer_attrs" json:"answer_attrs" validate:"required,dive"`
}

// AnswerAttrs struct to describe answer attributes.
type AnswerAttrs struct {
	Description string   `json:"description" validate:"required"`
	Documents   []string `json:"documents"`
	Images      []string `json:"images"`
	Links       []string `json:"links"`
}

// GetAnswer struct to describe answer object.
type GetAnswer struct {
	ID           uuid.UUID   `db:"id" json:"id"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
	UserID       uuid.UUID   `db:"user_id" json:"user_id"`
	ProjectID    uuid.UUID   `db:"project_id" json:"project_id"`
	TaskID       uuid.UUID   `db:"task_id" json:"task_id"`
	Alias        string      `db:"alias" json:"alias"`
	AnswerStatus int         `db:"answer_status" json:"answer_status"`
	AnswerAttrs  AnswerAttrs `db:"answer_attrs" json:"answer_attrs"`

	// Fields for JOIN tables:
	Username  string    `db:"username" json:"username"`
	UserAttrs UserAttrs `db:"user_attrs" json:"user_attrs"`
}

// GetAnswers struct to describe answers list object.
type GetAnswers struct {
	ID          uuid.UUID   `db:"id" json:"id"`
	CreatedAt   time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updated_at"`
	Alias       string      `db:"alias" json:"alias"`
	AnswerAttrs AnswerAttrs `db:"answer_attrs" json:"answer_attrs"`

	// Fields for JOIN tables:
	Username  string    `db:"username" json:"username"`
	UserAttrs UserAttrs `db:"user_attrs" json:"user_attrs"`
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
