package server

import (
	"github.com/gofiber/fiber/v2"

	healthhttp "link-nest/internal/features/health/http"
	authhttp "link-nest/internal/features/auth/http"
	database "link-nest/internal/platform/database"
)

func New() (*fiber.App, error) {
	if err := database.Connect(); err != nil {
		return nil, err
	}
	app := fiber.New()

	healthhttp.RegisterRoutes(app)
	authhttp.RegisterRoutes(app)


	return app, nil
}
