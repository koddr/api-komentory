package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/platform/database"
	"fmt"
	"time"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAnswerByID func for get one answer by ID.
func GetAnswerByID(c *fiber.Ctx) error {
	// Catch answer ID from URL.
	answerID, err := uuid.Parse(c.Params("answer_id"))
	if err != nil {
		return utilities.CheckForError(c, err, 400, "answer id", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get one answer.
	answer, status, err := db.GetAnswerByID(answerID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "answer", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"answer": answer,
	})
}

// GetAnswersByProjectID func for get all exists answers by project ID.
func GetAnswersByProjectID(c *fiber.Ctx) error {
	// Catch project ID from URL.
	projectID, err := uuid.Parse(c.Params("project_id"))
	if err != nil {
		return utilities.CheckForError(c, err, 400, "project id", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get all answers.
	answers, status, err := db.GetAnswersByProjectID(projectID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "answers", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"count":   len(answers),
		"answers": answers,
	})
}

// GetAnswersByTaskID func for get all exists answers by task ID.
func GetAnswersByTaskID(c *fiber.Ctx) error {
	// Catch task ID from URL.
	taskID, err := uuid.Parse(c.Params("task_id"))
	if err != nil {
		return utilities.CheckForError(c, err, 400, "task id", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get all answers.
	answers, status, err := db.GetAnswersByTaskID(taskID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "answers", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"count":   len(answers),
		"answers": answers,
	})
}

// CreateAnswer func for create a new answer for project.
func CreateAnswer(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("answers", "create", false),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new Answer struct
	answer := &models.Answer{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(answer); err != nil {
		return utilities.CheckForError(c, err, 400, "answer", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedProject, status, err := db.GetProjectByID(answer.ProjectID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "project", err.Error())
	}

	// Checking, if answer with given ID is exists.
	foundedTask, status, err := db.GetTaskByID(answer.TaskID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "task", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Set initialized default data for answer:
	answer.ID = uuid.New()
	answer.CreatedAt = time.Now()
	answer.UserID = userID
	answer.ProjectID = foundedProject.ID
	answer.TaskID = foundedTask.ID
	answer.AnswerStatus = 0 // 0 == draft, 1 == active, 2 == blocked

	// Create a new validator for a Answer model.
	validate := utilities.NewValidator()

	// Validate answer fields.
	if err := validate.Struct(answer); err != nil {
		return utilities.CheckForError(
			c, err, 400, "answer", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
		)
	}

	// Create a new answer with given attrs.
	if err := db.CreateAnswer(answer); err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 400, "answer", err.Error())
	}

	// Return status 201 created.
	return c.SendStatus(fiber.StatusCreated)
}

// UpdateAnswer func for update answer by given ID.
func UpdateAnswer(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("answers", "update", true),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new Answer struct
	answer := &models.Answer{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(answer); err != nil {
		return utilities.CheckForError(c, err, 400, "answer", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if answer with given ID is exists.
	foundedAnswer, status, err := db.GetAnswerByID(answer.ID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "answer", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can update his answer.
	if foundedAnswer.UserID == userID {
		// Set initialized default data for answer:
		answer.UpdatedAt = time.Now()
		answer.UserID = userID
		answer.ProjectID = foundedAnswer.ProjectID
		answer.TaskID = foundedAnswer.TaskID

		// Create a new validator for a Answer model.
		validate := utilities.NewValidator()

		// Validate answer fields.
		if err := validate.Struct(answer); err != nil {
			return utilities.CheckForError(
				c, err, 400, "answer", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
			)
		}

		// Update answer by given ID.
		if err := db.UpdateAnswer(foundedAnswer.ID, answer); err != nil {
			return utilities.CheckForErrorWithStatusCode(c, err, 400, "answer", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "answer", "you have no permissions")

	}
}

// DeleteAnswer func for delete answer by given ID.
func DeleteAnswer(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("answers", "delete", true),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new Answer struct
	answer := &models.Answer{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(answer); err != nil {
		return utilities.CheckForError(c, err, 400, "answer", err.Error())
	}

	// Create a new validator for a Answer model.
	validate := utilities.NewValidator()

	// Validate answer fields.
	if err := validate.StructPartial(answer, "id"); err != nil {
		return utilities.CheckForError(
			c, err, 400, "answer", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
		)
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if answer with given ID is exists.
	foundedAnswer, status, err := db.GetAnswerByID(answer.ID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "answer", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his answer.
	if foundedAnswer.UserID == userID {
		// Delete answer by given ID.
		if err := db.DeleteAnswer(foundedAnswer.ID); err != nil {
			return utilities.CheckForError(c, err, 400, "answer", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "answer", "you have no permissions")
	}
}
