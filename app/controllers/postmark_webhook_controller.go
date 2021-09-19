package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/platform/database"
	"fmt"
	"os"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
)

// UpdateUserSubscription method to update user email subscriptions.
func UpdateUserSubscription(c *fiber.Ctx) error {
	// Define User-Agent Header.
	postmarkUserAgentHeader := c.Get("User-Agent")

	// Check, if User-Agent Header is set.
	if postmarkUserAgentHeader != os.Getenv("POSTMARK_USER_AGENT_HEADER") {
		return utilities.ThrowJSONErrorWithStatusCode(c, 400, "postmark webhook", "bad User-Agent header")
	}

	// Create a new user change email subscription struct.
	subscriptionChange := &models.PostmarkSuppressSendingWebhook{}

	// Checking received data from JSON body.
	if err := c.BodyParser(subscriptionChange); err != nil {
		return utilities.CheckForError(c, err, 400, "postmark webhook", err.Error())
	}

	// Create a new validator.
	validate := utilities.NewValidator()

	// Validate webhook fields.
	if err := validate.Struct(subscriptionChange); err != nil {
		return utilities.CheckForError(
			c, err, 400, "postmark webhook",
			fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
		)
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get user by given email.
	foundedUser, status, err := db.GetUserByEmail(subscriptionChange.Recipient)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "user", err.Error())
	}

	// Create a new user settings struct.
	userSettings := &models.UserSettings{}

	// If Postmark pushed SuppressSending attribute with false,
	// it means reactivation (user was subscribed again).
	if !subscriptionChange.SuppressSending {
		userSettings.EmailSubscriptions.Transactional = true
		userSettings.EmailSubscriptions.Marketing = true
	}

	// Validate user settings fields.
	if err := validate.Struct(userSettings); err != nil {
		return utilities.CheckForError(
			c, err, 400, "user settings", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
		)
	}

	// Change user settings with validated data.
	if err := db.UpdateUserSettings(foundedUser.ID, userSettings); err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 400, "user", err.Error())
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
