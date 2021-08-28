package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/platform/database"
	"time"

	"github.com/Komentory/repository"
	"github.com/Komentory/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAnswerByID func for get one answer by ID.
func GetAnswerByID(c *fiber.Ctx) error {
	// Catch task ID from URL.
	taskID, errParse := uuid.Parse(c.Params("answer_id"))
	if errParse != nil {
		// Return status 400 and bad request error.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errParse.Error(),
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

	// Get one answer.
	answer, status, errGetAnswerByID := db.GetAnswerByID(taskID)
	if errGetAnswerByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error":  true,
			"msg":    errGetAnswerByID.Error(),
			"answer": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":  false,
		"msg":    nil,
		"answer": answer,
	})
}

// GetAnswersByProjectID func for get all exists answers by project ID.
func GetAnswersByProjectID(c *fiber.Ctx) error {
	// Catch project ID from URL.
	projectID, errParse := uuid.Parse(c.Params("project_id"))
	if errParse != nil {
		// Return status 400 and bad request error.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errParse.Error(),
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

	// Get all answers.
	answers, status, errGetAnswersByProjectID := db.GetAnswersByProjectID(projectID)
	if errGetAnswersByProjectID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error":   true,
			"msg":     errGetAnswersByProjectID.Error(),
			"answers": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"count":   len(answers),
		"answers": answers,
	})
}

// GetAnswersByTaskID func for get all exists answers by task ID.
func GetAnswersByTaskID(c *fiber.Ctx) error {
	// Catch task ID from URL.
	taskID, errParse := uuid.Parse(c.Params("task_id"))
	if errParse != nil {
		// Return status 400 and bad request error.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errParse.Error(),
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

	// Get all answers.
	answers, status, errGetAnswersByTaskID := db.GetAnswersByTaskID(taskID)
	if errGetAnswersByTaskID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error":   true,
			"msg":     errGetAnswersByTaskID.Error(),
			"answers": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"count":   len(answers),
		"answers": answers,
	})
}

// CreateAnswer func for create a new task for project.
func CreateAnswer(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		repository.GenerateCredential("answers", "create", false),
	}

	// Validate JWT token.
	claims, errTokenValidate := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if errTokenValidate != nil {
		// Return status 401 and error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   errTokenValidate.Error(),
		})
	}

	// Create new Answer struct
	answer := &models.Answer{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(answer); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if project with given ID is exists.
	foundedProject, status, errGetProjectByID := db.GetProjectByID(answer.ProjectID)
	if errGetProjectByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error":  true,
			"msg":    errGetProjectByID.Error(),
			"answer": nil,
		})
	}

	// Checking, if task with given ID is exists.
	foundedTask, status, errGetTaskByID := db.GetTaskByID(answer.TaskID)
	if errGetTaskByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error":  true,
			"msg":    errGetTaskByID.Error(),
			"answer": nil,
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Create a new validator for a Answer model.
	validate := utilities.NewValidator()

	// Set initialized default data for answer:
	answer.ID = uuid.New()
	answer.CreatedAt = time.Now()
	answer.UserID = userID
	answer.ProjectID = foundedProject.ID
	answer.TaskID = foundedTask.ID
	answer.AnswerStatus = 0 // 0 == draft, 1 == active, 2 == blocked

	// Validate answer fields.
	if err := validate.Struct(answer); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utilities.ValidatorErrors(err),
		})
	}

	// Create a new answer with given attrs.
	if err := db.CreateAnswer(answer); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 201 created.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":  false,
		"msg":    nil,
		"answer": answer,
	})
}

// UpdateAnswer func for update task by given ID.
func UpdateAnswer(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		repository.GenerateCredential("answers", "update", true),
	}

	// Validate JWT token.
	claims, errTokenValidate := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if errTokenValidate != nil {
		// Return status 401 and error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   errTokenValidate.Error(),
		})
	}

	// Create new Answer struct
	task := &models.Answer{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(task); err != nil {
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

	// Checking, if project with given ID is exists.
	foundedAnswer, status, errGetAnswerByID := db.GetAnswerByID(task.ID)
	if errGetAnswerByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetAnswerByID.Error(),
			"task":  nil,
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can update his task.
	if foundedAnswer.UserID == userID {
		// Set initialized default data for task:
		task.UpdatedAt = time.Now()
		task.UserID = userID
		task.ProjectID = foundedAnswer.ProjectID
		task.TaskID = foundedAnswer.TaskID

		// Create a new validator for a Answer model.
		validate := utilities.NewValidator()

		// Validate project fields.
		if err := validate.Struct(task); err != nil {
			// Return 400, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utilities.ValidatorErrors(err),
			})
		}

		// Update task by given ID.
		if err := db.UpdateAnswer(foundedAnswer.ID, task); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 200 OK.
		return c.JSON(fiber.Map{
			"error": false,
			"msg":   nil,
			"task":  task,
		})
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   repository.GenerateErrorMessage(403, "user", "it's not your task"),
			"task":  nil,
		})
	}
}

// DeleteAnswer func for delete task by given ID.
func DeleteAnswer(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		repository.GenerateCredential("answers", "delete", true),
	}

	// Validate JWT token.
	claims, errTokenValidate := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if errTokenValidate != nil {
		// Return status 401 and error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   errTokenValidate.Error(),
		})
	}

	// Create new Answer struct
	project := &models.Answer{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(project); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Answer model.
	validate := utilities.NewValidator()

	// Validate project fields.
	if err := validate.StructPartial(project, "id"); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utilities.ValidatorErrors(err),
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

	// Checking, if project with given ID is exists.
	foundedAnswer, status, errGetAnswerByID := db.GetAnswerByID(project.ID)
	if errGetAnswerByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetAnswerByID.Error(),
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his task.
	if foundedAnswer.UserID == userID {
		// Delete project by given ID.
		if err := db.DeleteAnswer(foundedAnswer.ID); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   repository.GenerateErrorMessage(403, "user", "it's not your task"),
		})
	}
}
