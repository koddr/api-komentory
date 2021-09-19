package routes

import (
	"Komentory/api/app/controllers"
	"Komentory/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// WebhookRoutes func for describe group of webhook routes.
func WebhookRoutes(a *fiber.App) {
	// Create routes group.
	r := a.Group("/v1")

	// Routes for POST method (with BasicAuth):
	r.Post("/webhook/postmark/subscriptions", middleware.BasicAuthProtected(), controllers.UpdateUserSubscriptions) // update email subscriptions
}
