package middleware

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

// Cached func for specify routes group with cached the response.
// See: https://docs.gofiber.io/api/middleware/cache
func Cached() func(*fiber.Ctx) error {
	// Check environment variable.
	cacheExpirationMinutesCount, err := strconv.Atoi(os.Getenv("SERVER_CACHE_EXPIRATION_MINUTES_COUNT"))
	if err != nil {
		return cache.New(cache.ConfigDefault)
	}

	// Create config for Cache middleware.
	config := cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("no-cache") == "true" // if route has query ?no-cache=true, skip caching
		},
		Expiration:   time.Minute * time.Duration(cacheExpirationMinutesCount),
		CacheControl: true,
	}

	return cache.New(config)
}
