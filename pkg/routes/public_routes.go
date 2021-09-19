package routes

import (
	"Komentory/api/app/controllers"
	"Komentory/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	r := a.Group("/v1")

	// Routes for GET method (many, cached):
	r.Get("/projects", middleware.Cached(), controllers.GetProjects)                          // get all projects
	r.Get("/user/:username/projects", middleware.Cached(), controllers.GetProjectsByUsername) // get projects by username

	// Routes for GET method (many, non-cached):
	r.Get("/project/:project_id/tasks", controllers.GetTasksByProjectID)     // get tasks by project ID
	r.Get("/project/:project_id/answers", controllers.GetAnswersByProjectID) // get answers by project ID
	r.Get("/task/:task_id/answers", controllers.GetAnswersByTaskID)          // get answers by task ID

	// Routes for GET method (single, non-cached):
	r.Get("/project/:alias", controllers.GetProjectByAlias) // get one project by alias
	r.Get("/task/:alias", controllers.GetTaskByAlias)       // get one task by alias
	r.Get("/answer/:alias", controllers.GetAnswerByAlias)   // get one answer by ID
}
