package controllers

import (
	"time"

	"Komentory/api/app/models"
	"Komentory/api/platform/database"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetProjects func for get all exists projects.
func GetProjects(c *fiber.Ctx) error {
	// Create database connection.
	db, errOpenDBConnection := database.OpenDBConnection()
	if errOpenDBConnection != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errOpenDBConnection.Error(),
		})
	}

	// Get all projects.
	projects, errGetProjects := db.GetProjects()
	if errGetProjects != nil {
		// Return status 400 and bad request error.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   errGetProjects.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"count":    len(projects),
		"projects": projects,
	})
}

// GetProjectsByUserID func for get all exists projects by user ID.
func GetProjectsByUserID(c *fiber.Ctx) error {
	// Catch project ID from URL.
	userID, errParse := uuid.Parse(c.Params("user_id"))
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

	// Get all projects.
	projects, status, errGetProjectsByUserID := db.GetProjectsByUserID(userID)
	if errGetProjectsByUserID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetProjectsByUserID.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"count":    len(projects),
		"projects": projects,
	})
}

// GetProject func for get project by given project alias.
func GetProjectByAlias(c *fiber.Ctx) error {
	// Catch project ID from URL.
	alias := c.Params("alias")

	// Create database connection.
	db, errOpenDBConnection := database.OpenDBConnection()
	if errOpenDBConnection != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errOpenDBConnection.Error(),
		})
	}

	// Get project by ID.
	project, status, errGetProjectByAlias := db.GetProjectByAlias(alias)
	if errGetProjectByAlias != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetProjectByAlias.Error(),
		})
	}

	// Get all tasks for this project ID.
	tasks, status, errGetTasksByProjectID := db.GetTasksByProjectID(project.ID)
	if errGetTasksByProjectID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetTasksByProjectID.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":       false,
		"project":     project,
		"tasks_count": len(tasks),
		"tasks":       tasks,
	})
}

// CreateProject func for create a new project.
func CreateProject(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("projects", "create", false),
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

	// Create new Project struct
	project := &models.Project{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(project); err != nil {
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

	// Create a new validator for a Project model.
	validate := utilities.NewValidator()

	// Set initialized default data for project:
	project.ID = uuid.New()
	project.CreatedAt = time.Now()
	project.UserID = claims.UserID
	project.Alias = project.ID.String()[:4] + project.ID.String()[24:]
	project.ProjectStatus = 0 // 0 == draft, 1 == active, 2 == blocked

	// Validate project fields.
	if err := validate.Struct(project); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utilities.ValidatorErrors(err),
		})
	}

	// Create a new project with given attrs.
	if err := db.CreateProject(project); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 201 created.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"project": project,
	})
}

// UpdateProject func for update project by given ID.
func UpdateProject(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("projects", "update", true),
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

	// Create new Project struct
	project := &models.Project{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(project); err != nil {
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
	foundedProject, status, errGetProjectByID := db.GetProjectByID(project.ID)
	if errGetProjectByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetProjectByID.Error(),
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his project.
	if foundedProject.UserID == userID {
		// Set initialized default data for project:
		project.UserID = userID
		project.UpdatedAt = time.Now()

		// Create a new validator for a Project model.
		validate := utilities.NewValidator()

		// Validate project fields.
		if err := validate.Struct(project); err != nil {
			// Return 400, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utilities.ValidatorErrors(err),
			})
		}

		// Update project by given ID.
		if err := db.UpdateProject(foundedProject.ID, project); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 200 OK.
		return c.JSON(fiber.Map{
			"error":   false,
			"project": project,
		})
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   utilities.GenerateErrorMessage(403, "user", "it's not your project"),
		})
	}
}

// DeleteProject func for delete project by given ID.
func DeleteProject(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("projects", "delete", true),
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

	// Create new Project struct
	project := &models.Project{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(project); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Project model.
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
	foundedProject, status, errGetProjectByID := db.GetProjectByID(project.ID)
	if errGetProjectByID != nil {
		// Return status and error message.
		return c.Status(status).JSON(fiber.Map{
			"error": true,
			"msg":   errGetProjectByID.Error(),
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his project.
	if foundedProject.UserID == userID {
		// Delete project by given ID.
		if err := db.DeleteProject(foundedProject.ID); err != nil {
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
			"msg":   utilities.GenerateErrorMessage(403, "user", "it's not your project"),
		})
	}
}
