package routes

import (
	"Komentory/api/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")

	// Routes for GET method (arrays):
	route.Get("/projects", controllers.GetProjects)                              // get list of all projects
	route.Get("/project/:project_id/tasks", controllers.GetTasksByProjectID)     // get list of all tasks by project ID
	route.Get("/project/:project_id/answers", controllers.GetAnswersByProjectID) // get list of all answers by project ID
	route.Get("/task/:task_id/answers", controllers.GetAnswersByTaskID)          // get list of all answers by task ID
	route.Get("/user/:user_id/projects", controllers.GetProjectsByUserID)        // get list of all projects by user ID

	// Routes for GET method (single):
	route.Get("/project/:alias", controllers.GetProjectByAlias) // get one project by alias
	route.Get("/task/:task_id", controllers.GetTaskByID)        // get one task by ID
	route.Get("/answer/:answer_id", controllers.GetAnswerByID)  // get one answer by ID
}
