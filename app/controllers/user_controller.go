package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/platform/database"
	"fmt"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
)

// UpdateUserPassword method to update user password.
func UpdateUserPassword(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("users", "update", true),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create a new user change password struct.
	passwordChange := &models.UserChangePassword{}

	// Checking received data from JSON body.
	if err := c.BodyParser(passwordChange); err != nil {
		return utilities.CheckForError(c, err, 400, "user password", err.Error())
	}

	// Create a new validator for a User model.
	validate := utilities.NewValidator()

	// Validate sign up fields.
	if err := validate.Struct(passwordChange); err != nil {
		return utilities.CheckForError(
			c, err, 400, "task", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
		)
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Get user by given email.
	foundedUser, status, err := db.GetUserByID(userID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "user", err.Error())
	}

	// Compare given user password with stored in found user.
	matchUserPasswords := utilities.ComparePasswords(foundedUser.PasswordHash, passwordChange.OldPassword)
	if !matchUserPasswords {
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "user", "email or password")
	}

	// Set initialized default data for user:
	newPasswordHash := utilities.GeneratePassword(passwordChange.NewPassword)

	// Create a new user with validated data.
	if err := db.UpdateUserPassword(foundedUser.ID, newPasswordHash); err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 400, "user", err.Error())
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateUserAttrs method for update user attributes.
func UpdateUserAttrs(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("users", "update", true),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create a new user auth struct.
	userAttrs := &models.UserAttrs{}

	// Checking received data from JSON body.
	if err := c.BodyParser(userAttrs); err != nil {
		return utilities.CheckForError(c, err, 400, "user attrs", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get user by email.
	foundedUser, status, err := db.GetUserByID(claims.UserID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "user", err.Error())
	}

	// Update user attributes.
	err = db.UpdateUserAttrs(foundedUser.ID, userAttrs)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 400, "user attrs", err.Error())
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}
