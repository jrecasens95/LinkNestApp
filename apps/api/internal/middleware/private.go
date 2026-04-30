package middleware

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
)

func RequireAPIKey(private bool, apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !private {
			return c.Next()
		}
		if apiKey == "" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "private mode is enabled but API_KEY is not configured",
			})
		}

		provided := c.Get("X-API-Key")
		if subtle.ConstantTimeCompare([]byte(provided), []byte(apiKey)) != 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid api key",
			})
		}

		return c.Next()
	}
}
