package http

import (
	authdomain "link-nest/internal/features/auth/domain"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router) {
	service := authdomain.NewJWTService()

	app.Post("/login", func(c *fiber.Ctx) error {
		token, err := service.GenerateToken("starter-user")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "could not generate token",
			})
		}

		return c.JSON(fiber.Map{
			"token": token,
		})
	})
}
