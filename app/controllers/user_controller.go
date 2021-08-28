package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/pkg/utils"
	"Komentory/api/platform/database"

	"github.com/Komentory/repository"

	"github.com/gofiber/fiber/v2"
)

// UpdateUserPassword method to update user password.
func UpdateUserPassword(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		repository.GenerateCredential("users", "update", true),
	}

	// Validate JWT token.
	claims, errTokenValidate := utils.TokenValidateExpireTimeAndCredentials(c, credentials)
	if errTokenValidate != nil {
		// Return status 401 and error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   errTokenValidate.Error(),
		})
	}

	// Create a new user change password struct.
	passwordChange := &models.UserChangePassword{}

	// Checking received data from JSON body.
	if err := c.BodyParser(passwordChange); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a User model.
	validate := utils.NewValidator()

	// Validate sign up fields.
	if err := validate.Struct(passwordChange); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create database connection.
	db, errOpenDBConnection := database.OpenDBConnection()
	if errOpenDBConnection != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errOpenDBConnection.Error(),
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Get user by given email.
	foundedUser, status, errGetUserByID := db.GetUserByID(userID)
	if errGetUserByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetUserByID.Error(),
			"user":  nil,
		})
	}

	// Compare given user password with stored in found user.
	matchUserPasswords := utils.ComparePasswords(foundedUser.PasswordHash, passwordChange.OldPassword)
	if !matchUserPasswords {
		// Return status 403, if password is not compare to stored in database.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   repository.GenerateErrorMessage(403, "user", "email or password"),
		})
	}

	// Set initialized default data for user:
	newPasswordHash := utils.GeneratePassword(passwordChange.NewPassword)

	// Create a new user with validated data.
	if err := db.UpdateUserPassword(foundedUser.ID, newPasswordHash); err != nil {
		// Return status 400 and bad request error.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateUserAttrs method for update user attributes.
func UpdateUserAttrs(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		repository.GenerateCredential("users", "update", true),
	}

	// Validate JWT token.
	claims, errTokenValidate := utils.TokenValidateExpireTimeAndCredentials(c, credentials)
	if errTokenValidate != nil {
		// Return status 401 and error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   errTokenValidate.Error(),
		})
	}

	// Create a new user auth struct.
	userAttrs := &models.UserAttrs{}

	// Checking received data from JSON body.
	if err := c.BodyParser(userAttrs); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, errOpenDBConnection := database.OpenDBConnection()
	if errOpenDBConnection != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errOpenDBConnection.Error(),
		})
	}

	// Get user by email.
	foundedUser, status, errGetUserByID := db.GetUserByID(claims.UserID)
	if errGetUserByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetUserByID.Error(),
			"user":  nil,
		})
	}

	// Update user attributes.
	err := db.UpdateUserAttrs(foundedUser.ID, userAttrs)
	if err != nil {
		// Return, status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user": fiber.Map{
			"id":         foundedUser.ID,
			"user_attrs": userAttrs,
		},
	})
}
