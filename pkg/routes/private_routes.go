package routes

import (
	"Komentory/api/app/controllers"
	"Komentory/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1", middleware.JWTProtected())

	// Routes for POST method:
	route.Post("/project", controllers.CreateNewProject) // create a new project
	route.Post("/task", controllers.CreateNewTask)       // create a new task
	route.Post("/answer", controllers.CreateNewAnswer)   // create a new answer

	// Routes for PATCH method:
	route.Patch("/project", controllers.UpdateProject) // update one project
	route.Patch("/task", controllers.UpdateTask)       // update one task
	route.Patch("/answer", controllers.UpdateAnswer)   // update one answer

	// Routes for PUT method:
	route.Put("/cdn/upload", controllers.PutFileToCDN) // upload file object to CDN

	// Routes for DELETE method:
	route.Delete("/project", controllers.DeleteProject)        // delete one project
	route.Delete("/task", controllers.DeleteTask)              // delete one task
	route.Delete("/answer", controllers.DeleteAnswer)          // delete one answer
	route.Delete("/cdn/remove", controllers.RemoveFileFromCDN) // remove one file from CDN
}
