package routes

import (
	"Komentory/api/app/controllers"
	"Komentory/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	r := a.Group("/v1")

	// Routes for POST method:
	r.Post("/project", middleware.JWTProtected(), controllers.CreateNewProject) // create a new project
	r.Post("/task", middleware.JWTProtected(), controllers.CreateNewTask)       // create a new task
	r.Post("/answer", middleware.JWTProtected(), controllers.CreateNewAnswer)   // create a new answer

	// Routes for PATCH method:
	r.Patch("/project", middleware.JWTProtected(), controllers.UpdateProject) // update one project
	r.Patch("/task", middleware.JWTProtected(), controllers.UpdateTask)       // update one task
	r.Patch("/answer", middleware.JWTProtected(), controllers.UpdateAnswer)   // update one answer

	// Routes for PUT method:
	r.Put("/cdn/upload", middleware.JWTProtected(), controllers.PutFileToCDN) // upload file object to CDN

	// Routes for DELETE method:
	r.Delete("/project", middleware.JWTProtected(), controllers.DeleteProject)        // delete one project
	r.Delete("/task", middleware.JWTProtected(), controllers.DeleteTask)              // delete one task
	r.Delete("/answer", middleware.JWTProtected(), controllers.DeleteAnswer)          // delete one answer
	r.Delete("/cdn/remove", middleware.JWTProtected(), controllers.RemoveFileFromCDN) // remove one file from CDN
}
