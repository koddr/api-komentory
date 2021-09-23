package controllers

import (
	"fmt"
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
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get all projects.
	projects, status, err := db.GetProjects()
	if err != nil {
		return utilities.CheckForError(c, err, status, "projects", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status":   fiber.StatusOK,
		"count":    len(projects),
		"projects": projects,
	})
}

// GetProjectsByUsername func for get all exists projects by given username.
func GetProjectsByUsername(c *fiber.Ctx) error {
	// Catch project ID from URL.
	username := c.Params("username")

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get all projects by username.
	projects, status, err := db.GetProjectsByUsername(username)
	if err != nil {
		return utilities.CheckForError(c, err, status, "projects", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status":   fiber.StatusOK,
		"count":    len(projects),
		"projects": projects,
	})
}

// GetProject func for get project by given project alias.
func GetProjectByAlias(c *fiber.Ctx) error {
	// Catch project ID from URL.
	alias := c.Params("alias")

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get project by ID.
	project, status, err := db.GetProjectByAlias(alias)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "project", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"project": project,
	})
}

// CreateNewProject func for create a new project.
func CreateNewProject(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("projects", "create", false),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new ProjectAttrs struct
	projectAttrs := &models.ProjectAttrs{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(projectAttrs); err != nil {
		return utilities.CheckForError(c, err, 400, "project", err.Error())
	}

	// Generate random string for the project's alias.
	randomAlias, err := utilities.GenerateNewNanoID("", 24)
	if err != nil {
		return utilities.CheckForError(c, err, 400, "project alias", err.Error())
	}

	// Create new ProjectAttrs struct
	project := &models.Project{}

	// Set initial data for project:
	project.ID = uuid.New()
	project.CreatedAt = time.Now()
	project.UserID = claims.UserID
	project.Alias = randomAlias
	project.ProjectStatus = 0 // 0 == draft, 1 == active, 2 == unpublished

	// Set project attributes from request body:
	project.ProjectAttrs = *projectAttrs

	// Create a new validator for a Project model.
	validate := utilities.NewValidator()

	// Validate project fields.
	if err := validate.Struct(project); err != nil {
		return utilities.CheckForError(
			c, err, 400, "project", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
		)
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Create a new project with given attrs.
	if err := db.CreateNewProject(project); err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 400, "project", err.Error())
	}

	// Return status 201 created.
	return c.SendStatus(fiber.StatusCreated)
}

// UpdateProject func for update project by given ID.
func UpdateProject(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("projects", "update", true),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new Project struct
	project := &models.Project{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(project); err != nil {
		return utilities.CheckForError(c, err, 400, "project", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedProject, status, err := db.GetProjectByID(project.ID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "project", err.Error())
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
			return utilities.CheckForError(
				c, err, 400, "project", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
			)
		}

		// Update project by given ID.
		if err := db.UpdateProject(foundedProject.ID, project); err != nil {
			return utilities.CheckForErrorWithStatusCode(c, err, 400, "project", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "project", "you have no permissions")
	}
}

// DeleteProject func for delete project by given ID.
func DeleteProject(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("projects", "delete", true),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 401, "jwt", err.Error())
	}

	// Create new Project struct
	project := &models.Project{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(project); err != nil {
		return utilities.CheckForError(c, err, 400, "project", err.Error())
	}

	// Create a new validator for a Project model.
	validate := utilities.NewValidator()

	// Validate project fields.
	if err := validate.StructPartial(project, "id"); err != nil {
		return utilities.CheckForError(
			c, err, 400, "project", fmt.Sprintf("validation error, %v", utilities.ValidatorErrors(err)),
		)
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedProject, status, err := db.GetProjectByID(project.ID)
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, status, "project", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his project.
	if foundedProject.UserID == userID {
		// Delete project by given ID.
		if err := db.DeleteProject(foundedProject.ID); err != nil {
			return utilities.CheckForErrorWithStatusCode(c, err, 400, "project", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONErrorWithStatusCode(c, 403, "project", "you have no permissions")
	}
}
