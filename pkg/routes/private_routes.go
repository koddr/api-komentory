package routes

import (
	"Komentory/api/app/controllers"
	"Komentory/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	r := a.Group("/v1", middleware.JWTProtected())

	// Routes for POST method:
	r.Post("/project", controllers.CreateNewProject) // create a new project
	r.Post("/task", controllers.CreateNewTask)       // create a new task
	r.Post("/answer", controllers.CreateNewAnswer)   // create a new answer

	// Routes for PATCH method:
	r.Patch("/project", controllers.UpdateProject) // update one project
	r.Patch("/task", controllers.UpdateTask)       // update one task
	r.Patch("/answer", controllers.UpdateAnswer)   // update one answer

	// Routes for PUT method:
	r.Put("/cdn/upload", controllers.PutFileToCDN) // upload file object to CDN

	// Routes for DELETE method:
	r.Delete("/project", controllers.DeleteProject)        // delete one project
	r.Delete("/task", controllers.DeleteTask)              // delete one task
	r.Delete("/answer", controllers.DeleteAnswer)          // delete one answer
	r.Delete("/cdn/remove", controllers.RemoveFileFromCDN) // remove one file from CDN
}
