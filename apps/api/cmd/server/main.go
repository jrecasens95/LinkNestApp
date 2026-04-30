package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jrecasens95/link-nest/backend/internal/config"
	"github.com/jrecasens95/link-nest/backend/internal/database"
	"github.com/jrecasens95/link-nest/backend/internal/handlers"
	"github.com/jrecasens95/link-nest/backend/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "LinkNest API",
		ServerHeader: "LinkNest",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	linkService := services.NewLinkService(db)
	linkHandler := handlers.NewLinkHandler(cfg.BaseURL, linkService)

	app.Get("/api/health", handlers.Health)
	app.Get("/api/links", linkHandler.List)
	app.Post("/api/links", linkHandler.Create)
	app.Get("/api/links/:id", linkHandler.Get)
	app.Patch("/api/links/:id", linkHandler.Update)
	app.Delete("/api/links/:id", linkHandler.Delete)
	app.Get("/:code", linkHandler.Redirect)

	log.Printf("LinkNest API listening on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
