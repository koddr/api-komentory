package middleware

import (
	"os"

	"github.com/Komentory/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// BasicAuthProtected func for specify routes group with BasicAuth authentication.
// See: https://docs.gofiber.io/api/middleware/basicauth
func BasicAuthProtected() func(*fiber.Ctx) error {
	// Create config for BasicAuth authentication middleware.
	config := basicauth.Config{
		Users: map[string]string{
			os.Getenv("POSTMARK_BASICAUTH_USER"): os.Getenv("POSTMARK_BASICAUTH_PASSWORD"), // for Postmark
		},
		Realm: "Forbidden",
		Unauthorized: func(c *fiber.Ctx) error {
			return utilities.ThrowJSONErrorWithStatusCode(c, 403, "basic auth", "you have no permissions")
		},
	}

	return basicauth.New(config)
}
