package routes

import (
	"Komentory/api/app/controllers"
	"Komentory/api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// WebhookRoutes func for describe group of webhook routes.
func WebhookRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")

	// Routes for POST method:
	route.Post("/webhook/postmark/subscription", middleware.BasicAuthProtected(), controllers.UpdateUserSubscription) // update Postmark subscription
}
