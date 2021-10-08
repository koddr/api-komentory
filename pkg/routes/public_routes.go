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
	r.Get("/projects", middleware.Cached(), controllers.GetProjects)                       // get all projects
	r.Get("/user/:user_id/projects", middleware.Cached(), controllers.GetProjectsByUserID) // get projects by user ID

	// Routes for GET method (single, cached):
	r.Get("/project/:project_id", middleware.Cached(), controllers.GetProjectByID) // get one project by ID

	// Routes for GET method (many, non-cached):
	r.Get("/project/:project_id/tasks", controllers.GetTasksByProjectID)     // get tasks by project ID
	r.Get("/project/:project_id/answers", controllers.GetAnswersByProjectID) // get answers by project ID
	r.Get("/task/:task_id/answers", controllers.GetAnswersByTaskID)          // get answers by task ID

	// Routes for GET method (single, non-cached):
	r.Get("/task/:task_id", controllers.GetTaskByID)          // get one task by ID
	r.Get("/answer/:answer_id", controllers.GetAnswerByAlias) // get one answer by ID
}
