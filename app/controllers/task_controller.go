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

// GetTaskByID func for get one task by ID.
func GetTaskByID(c *fiber.Ctx) error {
	// Catch task ID from URL.
	taskID, err := uuid.Parse(c.Params("task_id"))
	if err != nil {
		return utilities.CheckForError(c, err, 400, "task", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get one task.
	task, status, err := db.GetTaskByID(taskID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "task", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"task":   task,
	})
}

// GetTaskByAlias func for get one task by alias.
func GetTaskByAlias(c *fiber.Ctx) error {
	// Catch task alias from URL.
	alias := c.Params("alias")

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get one task.
	task, status, err := db.GetTaskByAlias(alias)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "task", err.Error())
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
		return utilities.CheckForErrorWithStatusCode(c, err, status, "tasks", err.Error())
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
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new Task struct
	task := &models.Task{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(task); err != nil {
		return utilities.CheckForError(c, err, 400, "task", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedProject, status, err := db.GetProjectByID(task.ProjectID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "project", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can add a new task for his project.
	if foundedProject.UserID == userID {
		// Generate random string for the project's alias.
		randomAlias, err := utilities.GenerateNewNanoID(utilities.LowerCaseWithoutDashesChars, 16)
		if err != nil {
			return utilities.CheckForError(c, err, 400, "task alias", err.Error())
		}

		// Set initialized default data for task:
		task.ID = uuid.New()
		task.CreatedAt = time.Now()
		task.UserID = userID
		task.ProjectID = foundedProject.ID
		task.Alias = randomAlias
		task.TaskStatus = 0 // 0 == draft, 1 == active, 2 == blocked

		// Create a new validator for a Task model.
		validate := utilities.NewValidator()

		// Validate task fields.
		if err := validate.Struct(task); err != nil {
			return utilities.CheckForError(
				c, err, 400, "task", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
			)
		}

		// Create a new task with given attrs.
		if err := db.CreateNewTask(task); err != nil {
			return utilities.CheckForErrorWithStatusCode(c, err, 400, "task", err.Error())
		}

		// Return status 201 created.
		return c.SendStatus(fiber.StatusCreated)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "project", "you have no permissions")
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
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new Task struct
	task := &models.Task{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(task); err != nil {
		return utilities.CheckForError(c, err, 400, "task", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedTask, status, err := db.GetTaskByID(task.ID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "task", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his task.
	if foundedTask.UserID == userID {
		// Set default data for task:
		task.UpdatedAt = time.Now()
		task.UserID = userID
		task.ProjectID = foundedTask.ProjectID

		// Create a new validator for a Task model.
		validate := utilities.NewValidator()

		// Validate task fields.
		if err := validate.Struct(task); err != nil {
			return utilities.CheckForError(
				c, err, 400, "task", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
			)
		}

		// Update task by given ID.
		if err := db.UpdateTask(foundedTask.ID, task); err != nil {
			return utilities.CheckForErrorWithStatusCode(c, err, 400, "task", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "project", "you have no permissions")
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
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new Task struct
	task := &models.Task{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(task); err != nil {
		return utilities.CheckForError(c, err, 400, "task", err.Error())
	}

	// Create a new validator for a Task model.
	validate := utilities.NewValidator()

	// Validate task fields.
	if err := validate.StructPartial(task, "id"); err != nil {
		return utilities.CheckForError(
			c, err, 400, "task", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
		)
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if task with given ID is exists.
	foundedTask, status, err := db.GetTaskByID(task.ID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "task", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his task.
	if foundedTask.UserID == userID {
		// Delete task by given ID.
		if err := db.DeleteTask(foundedTask.ID); err != nil {
			return utilities.CheckForErrorWithStatusCode(c, err, 400, "task", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "task", "you have no permissions")
	}
}
