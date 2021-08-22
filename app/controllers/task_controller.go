package controllers

import (
	"Komentory/api/app/models"
	"Komentory/api/pkg/repository"
	"Komentory/api/pkg/utils"
	"Komentory/api/platform/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetTaskByID func for get one task by ID.
func GetTaskByID(c *fiber.Ctx) error {
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

	// Get one task.
	task, status, errGetTaskByID := db.GetTaskByID(taskID)
	if errGetTaskByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetTaskByID.Error(),
			"task":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"task":  task,
	})
}

// GetTasksByProjectID func for get all exists tasks by project ID.
func GetTasksByProjectID(c *fiber.Ctx) error {
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

	// Get all tasks.
	tasks, status, errGetTasksByProjectID := db.GetTasksByProjectID(projectID)
	if errGetTasksByProjectID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetTasksByProjectID.Error(),
			"tasks": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(tasks),
		"tasks": tasks,
	})
}

// CreateTask func for create a new task for project.
func CreateTask(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		repository.TaskCreateCredential,
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

	// Create new Task struct
	task := &models.Task{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(task); err != nil {
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
	foundedProject, status, errGetProjectByID := db.GetProjectByID(task.ProjectID)
	if errGetProjectByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetProjectByID.Error(),
			"task":  nil,
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can add a new task for his project.
	if foundedProject.UserID == userID {
		// Create a new validator for a Task model.
		validate := utils.NewValidator()

		// Set initialized default data for task:
		task.ID = uuid.New()
		task.CreatedAt = time.Now()
		task.UserID = userID
		task.ProjectID = foundedProject.ID
		task.TaskStatus = 0 // 0 == draft, 1 == active, 2 == blocked

		// Validate task fields.
		if err := validate.Struct(task); err != nil {
			// Return, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utils.ValidatorErrors(err),
			})
		}

		// Create a new task with given attrs.
		if err := db.CreateTask(task); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 201 created.
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error": false,
			"msg":   nil,
			"task":  task,
		})
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   repository.PermissionDeniedCreateTask,
			"task":  nil,
		})
	}
}

// UpdateTask func for update task by given ID.
func UpdateTask(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		repository.TaskOwnUpdateCredential,
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

	// Create new Task struct
	task := &models.Task{}

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
	foundedTask, status, errGetTaskByID := db.GetTaskByID(task.ID)
	if errGetTaskByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetTaskByID.Error(),
			"task":  nil,
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his task.
	if foundedTask.UserID == userID {
		// Set initialized default data for task:
		task.UpdatedAt = time.Now()
		task.UserID = userID
		task.ProjectID = foundedTask.ProjectID

		// Create a new validator for a Task model.
		validate := utils.NewValidator()

		// Validate project fields.
		if err := validate.Struct(task); err != nil {
			// Return 400, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utils.ValidatorErrors(err),
			})
		}

		// Update task by given ID.
		if err := db.UpdateTask(foundedTask.ID, task); err != nil {
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
			"msg":   repository.PermissionDeniedUpdateTask,
			"task":  nil,
		})
	}
}

// DeleteTask func for delete task by given ID.
func DeleteTask(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		repository.TaskOwnDeleteCredential,
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

	// Create new Task struct
	project := &models.Task{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(project); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Task model.
	validate := utils.NewValidator()

	// Validate project fields.
	if err := validate.StructPartial(project, "id"); err != nil {
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

	// Checking, if project with given ID is exists.
	foundedTask, status, errGetTaskByID := db.GetTaskByID(project.ID)
	if errGetTaskByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetTaskByID.Error(),
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his task.
	if foundedTask.UserID == userID {
		// Delete project by given ID.
		if err := db.DeleteTask(foundedTask.ID); err != nil {
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
			"msg":   repository.PermissionDeniedDeleteTask,
		})
	}
}
