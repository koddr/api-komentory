package controllers

import (
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

// GetProjectByID func for get project by given project ID.
func GetProjectByID(c *fiber.Ctx) error {
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

	// Get project by ID.
	project, status, err := db.GetProjectByID(projectID)
	if err != nil {
		return utilities.CheckForError(c, err, status, "project", err.Error())
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"project": project,
	})
}

// GetProjectsByUserID func for get all exists projects by given user ID.
func GetProjectsByUserID(c *fiber.Ctx) error {
	// Catch project ID from URL.
	userID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return utilities.CheckForError(c, err, 400, "user id", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Get all projects by username.
	projects, status, err := db.GetProjectsByUserID(userID)
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

// CreateNewProject func for create a new project.
func CreateNewProject(c *fiber.Ctx) error {
	// Set needed credentials.
	credentials := []string{
		utilities.GenerateCredential("projects", "create", false),
	}

	// Validate JWT token.
	claims, err := utilities.TokenValidateExpireTimeAndCredentials(c, credentials)
	if err != nil {
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Create a new struct for JSON body.
	jsonBody := &models.CreateNewProject{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(jsonBody); err != nil {
		return utilities.CheckForError(c, err, 400, "project", err.Error())
	}

	// Create new Project struct.
	project := &models.Project{}

	// Set initial data for project:
	project.ID = uuid.New()
	project.UserID = claims.UserID

	// Set project attributes from JSON body:
	project.ProjectStatus = jsonBody.ProjectStatus // 0 == draft, 1 == active, 2 == unpublished
	project.ProjectAttrs = jsonBody.ProjectAttrs

	// Create a new validator for a Project model.
	validate := utilities.NewValidator()

	// Validate project fields.
	if err := validate.Struct(project); err != nil {
		return utilities.CheckForValidationError(c, err, 400, "project")
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Create a new project with given attrs.
	if err := db.CreateNewProject(project); err != nil {
		return utilities.CheckForError(c, err, 400, "project", err.Error())
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
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Create a new struct for JSON body.
	jsonBody := &models.UpdateProject{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(jsonBody); err != nil {
		return utilities.CheckForError(c, err, 400, "project", err.Error())
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedProject, status, err := db.FindProjectByID(jsonBody.ID)
	if err != nil {
		return utilities.CheckForError(c, err, status, "project", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his project.
	if foundedProject.UserID == userID {
		// Create new Project struct.
		project := &models.Project{}

		// Set project attributes from JSON body:
		project.ProjectStatus = jsonBody.ProjectStatus // 0 == draft, 1 == active, 2 == unpublished
		project.ProjectAttrs = jsonBody.ProjectAttrs

		// Create a new validator for a Project model.
		validate := utilities.NewValidator()

		// Validate project fields.
		if err := validate.Struct(project); err != nil {
			return utilities.CheckForValidationError(c, err, 400, "project")
		}

		// Update project by given ID.
		if err := db.UpdateProject(foundedProject.ID, project); err != nil {
			return utilities.CheckForError(c, err, 400, "project", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONError(c, 403, "project", "you have no permissions")
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
		return utilities.CheckForError(c, err, 401, "jwt", err.Error())
	}

	// Create a new struct for JSON body.
	jsonBody := &models.DeleteProject{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(jsonBody); err != nil {
		return utilities.CheckForError(c, err, 400, "project", err.Error())
	}

	// Create a new validator for a Project model.
	validate := utilities.NewValidator()

	// Validate project fields.
	if err := validate.Struct(jsonBody); err != nil {
		return utilities.CheckForValidationError(c, err, 400, "project")
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		return utilities.CheckForErrorWithStatusCode(c, err, 500, "database", err.Error())
	}

	// Checking, if project with given ID is exists.
	foundedProject, status, err := db.FindProjectByID(jsonBody.ID)
	if err != nil {
		return utilities.CheckForError(c, err, status, "project", err.Error())
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his project.
	if foundedProject.UserID == userID {
		// Delete project by given ID.
		if err := db.DeleteProject(foundedProject.ID); err != nil {
			return utilities.CheckForError(c, err, 400, "project", err.Error())
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return utilities.ThrowJSONError(c, 403, "project", "you have no permissions")
	}
}
