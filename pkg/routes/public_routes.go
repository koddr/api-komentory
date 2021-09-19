package routes

import (
	"Komentory/api/app/controllers"
	"Komentory/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")

	// Routes for GET method (many):
	route.Get("/projects", middleware.Cached(), controllers.GetProjects)                          // get all projects
	route.Get("/project/:project_id/tasks", controllers.GetTasksByProjectID)                      // get tasks by project ID
	route.Get("/project/:project_id/answers", controllers.GetAnswersByProjectID)                  // get answers by project ID
	route.Get("/task/:task_id/answers", controllers.GetAnswersByTaskID)                           // get answers by task ID
	route.Get("/user/:username/projects", middleware.Cached(), controllers.GetProjectsByUsername) // get projects by username

	// Routes for GET method (single):
	route.Get("/project/:alias", controllers.GetProjectByAlias) // get one project by alias
	route.Get("/task/:alias", controllers.GetTaskByAlias)       // get one task by alias
	route.Get("/answer/:alias", controllers.GetAnswerByAlias)   // get one answer by ID
}
