package models

import "github.com/google/uuid"

// User struct to describe User object.
type User struct {
	ID         uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Email      string    `db:"email" json:"email" validate:"required,email,lte=255"`
	Username   string    `db:"username" json:"username" validate:"required,lte=18"`
	UserStatus int       `db:"user_status" json:"user_status" validate:"int"`
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
