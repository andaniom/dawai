package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/violin-assessment/dawai/internal/handlers"
	"github.com/violin-assessment/dawai/internal/middleware"
	"github.com/violin-assessment/dawai/internal/services"
)

func main() {
	app := fiber.New()

	// Public routes
	authService := services.NewAuthService()
	authHandler := handlers.NewAuthHandler(authService)

	app.Post("/api/auth/login", authHandler.Login)
	app.Post("/api/auth/logout", middleware.JWTGuard, authHandler.Logout)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.Listen(":" + port)
}
