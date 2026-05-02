package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jrecasens95/link-nest/backend/internal/auth"
	"github.com/jrecasens95/link-nest/backend/internal/config"
	"github.com/jrecasens95/link-nest/backend/internal/database"
	"github.com/jrecasens95/link-nest/backend/internal/handlers"
	"github.com/jrecasens95/link-nest/backend/internal/middleware"
	"github.com/jrecasens95/link-nest/backend/internal/security"
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
		BodyLimit:    1024 * 1024,
	})

	app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(cors.New())

	linkService := services.NewLinkService(db)
	authService := auth.NewService(db, cfg.JWTSecret)
	urlValidator := security.NewURLValidator(cfg.MaxURLLength, cfg.BlacklistedDomains)
	linkHandler := handlers.NewLinkHandler(cfg.BaseURL, linkService, urlValidator)
	authHandler := handlers.NewAuthHandler(authService, cfg.InviteOnly, cfg.APIKey)
	requireAuth := middleware.RequireAuth(authService)
	createLinkLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: time.Minute,
	})
	privateCreate := middleware.RequireAPIKey(cfg.Private, cfg.APIKey)

	app.Get("/api/health", handlers.Health)
	app.Post("/api/auth/register", authHandler.Register)
	app.Post("/api/auth/login", authHandler.Login)
	app.Get("/api/auth/me", requireAuth, authHandler.Me)
	app.Get("/api/links", requireAuth, linkHandler.List)
	app.Post("/api/links", requireAuth, createLinkLimiter, privateCreate, linkHandler.Create)
	app.Get("/api/links/:id", requireAuth, linkHandler.Get)
	app.Get("/api/links/:id/stats", requireAuth, linkHandler.Stats)
	app.Patch("/api/links/:id", requireAuth, linkHandler.Update)
	app.Delete("/api/links/:id", requireAuth, linkHandler.Delete)
	app.Get("/:code", linkHandler.Redirect)

	log.Printf("LinkNest API listening on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
