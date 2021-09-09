package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Project struct to describe project object.
type Project struct {
	ID            uuid.UUID    `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
	UserID        uuid.UUID    `db:"user_id" json:"user_id" validate:"required,uuid"`
	Alias         string       `db:"alias" json:"alias" validate:"required,lte=24"`
	ProjectStatus int          `db:"project_status" json:"project_status" validate:"int"`
	ProjectAttrs  ProjectAttrs `db:"project_attrs" json:"project_attrs" validate:"required,dive"`
}

// ProjectAttrs struct to describe project attributes.
type ProjectAttrs struct {
	Title       string `json:"title" validate:"required,lte=255"`
	Description string `json:"description" validate:"required"`
	Picture     string `json:"picture"`
	WebsiteURL  string `json:"website_url"`
}

// Value make the ProjectAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (p *ProjectAttrs) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan make the ProjectAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (p *ProjectAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(j, &p)
}
