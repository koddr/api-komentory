package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/platform/database"
	"time"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetTaskByID func for get one task by ID.
func GetTaskByID(c *fiber.Ctx) error {
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

	// Get one task.
	task, status, err := db.GetTaskByID(taskID)
	if err != nil {
		return utilities.CheckForError(c, err, status, "task", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"task":   task,
	})
}

// GetTasksByProjectID func for get all exists tasks by project ID.
func GetTasksByProjectID(c *fiber.Ctx) error {
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

	// Get all tasks.
	tasks, status, err := db.GetTasksByProjectID(projectID)
	if err != nil {
		return utilities.CheckForError(c, err, status, "tasks", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"count":  len(tasks),
		"tasks":  tasks,
	})
}

// CreateNewTask func for create a new task for project.
func CreateNewTask(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("tasks", "create", false),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Create a new struct for JSON body.
	jsonBody := &models.CreateNewTask{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(jsonBody); err != nil {
		return utilities.CheckForError(c, err, 400, "task", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedProject, status, err := db.FindProjectByID(jsonBody.ProjectID)
	if err != nil {
		return utilities.CheckForError(c, err, status, "project", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can add a new task for his project.
	if foundedProject.UserID == userID {
		// Create new Task struct.
		task := &models.Task{}

		// Set initialized default data for task:
		task.ID = uuid.New()
		task.CreatedAt = time.Now()
		task.UserID = userID

		// Set project attributes from JSON body:
		task.ProjectID = jsonBody.ProjectID
		task.TaskStatus = jsonBody.TaskStatus // 0 == draft, 1 == active, 2 == unpublished
		task.TaskAttrs = jsonBody.TaskAttrs

		// Create a new validator for a Task model.
		validate := utilities.NewValidator()

		// Validate task fields.
		if err := validate.Struct(task); err != nil {
			return utilities.CheckForValidationError(c, err, 400, "task")
		}

		// Create a new task with given attrs.
		if err := db.CreateNewTask(task); err != nil {
			return utilities.CheckForError(c, err, 400, "task", err.Error())
		}

		// Return status 201 created.
		return c.SendStatus(fiber.StatusCreated)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONError(c, 403, "project", "you have no permissions")
	}
}

// UpdateTask func for update task by given ID.
func UpdateTask(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("tasks", "update", true),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Create a new struct for JSON body.
	jsonBody := &models.UpdateTask{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(jsonBody); err != nil {
		return utilities.CheckForError(c, err, 400, "task", err.Error())
	}

	// Create a new validator.
	validate := utilities.NewValidator()

	// Validate task fields.
	if err := validate.Struct(jsonBody); err != nil {
		return utilities.CheckForValidationError(c, err, 400, "task")
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedTask, status, err := db.FindTaskByID(jsonBody.ID)
	if err != nil {
		return utilities.CheckForError(c, err, status, "task", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his task.
	if foundedTask.UserID == userID {
		// Update task by given ID.
		if err := db.UpdateTask(foundedTask.ID, jsonBody); err != nil {
			return utilities.CheckForError(c, err, 400, "task", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONError(c, 403, "task", "you have no permissions")
	}
}

// DeleteTask func for delete task by given ID.
func DeleteTask(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("tasks", "delete", true),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Create a new struct for JSON body.
	jsonBody := &models.DeleteTask{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(jsonBody); err != nil {
		return utilities.CheckForError(c, err, 400, "task", err.Error())
	}

	// Create a new validator.
	validate := utilities.NewValidator()

	// Validate task fields.
	if err := validate.Struct(jsonBody); err != nil {
		return utilities.CheckForValidationError(c, err, 400, "task")
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if task with given ID is exists.
	foundedTask, status, err := db.GetTaskByID(jsonBody.ID)
	if err != nil {
		return utilities.CheckForError(c, err, status, "task", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his task.
	if foundedTask.UserID == userID {
		// Delete task by given ID.
		if err := db.DeleteTask(jsonBody.ID); err != nil {
			return utilities.CheckForError(c, err, 400, "task", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONError(c, 403, "task", "you have no permissions")
	}
}
