package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

// User struct to describe User object.
type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

// UserAttrs struct to describe user attributes.
type UserAttrs struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Picture   string   `json:"picture"`
	Abilities []string `json:"abilities"`
}

// UserSettings struct to describe user settings.
type UserSettings struct {
	EmailSubscriptions EmailSubscriptions `json:"email_subscriptions"`
}

// EmailSubscriptions struct to describe user email subscriptions.
type EmailSubscriptions struct {
	Transactional bool `json:"transactional"` // like "forgot password"
	Marketing     bool `json:"marketing"`     // like "invite friends and get X"
}

// Value make the UserAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (u UserAttrs) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Scan make the UserAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (u *UserAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &u)
}
