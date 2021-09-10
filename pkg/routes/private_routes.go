package routes

import (
	"Komentory/api/app/controllers"
	"Komentory/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")

	// Routes for GET method:
	route.Get("/cdn/list", middleware.JWTProtected(), controllers.GetFileListFromCDN) // get file list from CDN

	// Routes for POST method:
	route.Post("/project", middleware.JWTProtected(), controllers.CreateProject) // create a new project
	route.Post("/task", middleware.JWTProtected(), controllers.CreateTask)       // create a new task
	route.Post("/answer", middleware.JWTProtected(), controllers.CreateAnswer)   // create a new answer

	// Routes for PUT method:
	route.Put("/project", middleware.JWTProtected(), controllers.UpdateProject)   // update one project
	route.Put("/task", middleware.JWTProtected(), controllers.UpdateTask)         // update one task
	route.Put("/answer", middleware.JWTProtected(), controllers.UpdateAnswer)     // update one answer
	route.Put("/cdn/upload", middleware.JWTProtected(), controllers.PutFileToCDN) // upload file object to CDN

	// Routes for PATCH method:
	route.Patch("/user/edit/password", middleware.JWTProtected(), controllers.UpdateUserPassword) // update user password
	route.Patch("/user/edit/attrs", middleware.JWTProtected(), controllers.UpdateUserAttrs)       // update user attrs

	// Routes for DELETE method:
	route.Delete("/project", middleware.JWTProtected(), controllers.DeleteProject)        // delete one project
	route.Delete("/task", middleware.JWTProtected(), controllers.DeleteTask)              // delete one task
	route.Delete("/answer", middleware.JWTProtected(), controllers.DeleteAnswer)          // delete one answer
	route.Delete("/cdn/remove", middleware.JWTProtected(), controllers.RemoveFileFromCDN) // delete one file from CDN
}
