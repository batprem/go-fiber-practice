package main

// @title GFP API
// @version 1.0
// @description Fiber Practice API with OpenTelemetry and Swagger documentation
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3002
// @BasePath /
// @schemes http

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "gfp/api/docs" // Import generated docs
	otel "gfp/api/middlewares/otel"
	"gfp/api/routes"
	"gfp/api/routes/nested"
	"gfp/lib"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// HealthCheck godoc
// @Summary Health check
// @Description Returns a simple greeting message
// @Tags health
// @Produce plain
// @Success 200 {string} string "Hello, World!"
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	ctx := c.UserContext()
	otel.LogInfo(ctx, "Handling root endpoint")
	return c.SendString("Hello, World!")
}

// GetUserGreeting godoc
// @Summary Get user greeting
// @Description Returns a personalized greeting for the specified user
// @Tags users
// @Produce plain
// @Param user path string true "Username"
// @Success 200 {string} string "Personalized greeting"
// @Router /{user} [get]
func GetUserGreeting(c *fiber.Ctx) error {
	ctx := c.UserContext()
	user := c.Params("user")
	otel.LogInfo(ctx, "Handling user endpoint", "user", user)
	return c.SendString(lib.GetText(user))
}

func main() {
	// Initialize OpenTelemetry
	if err := otel.InitOpenTelemetry(); err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// Add OpenTelemetry middleware for automatic tracing
	app.Use(otelfiber.Middleware())

	// Swagger documentation route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Define routes
	app.Get("/", HealthCheck)
	app.Get("/:user", GetUserGreeting)

	app.Get(
		"simple-return/:user",
		routes.TextReturn,
	)

	app.Get(
		"simple-return/nested/:user",
		nested.NestedTextReturn,
	)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("🚀 Starting server on :3002")
		if err := app.Listen(":3002"); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("🛑 Shutting down server...")

	// Shutdown Fiber app
	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	// Shutdown OpenTelemetry
	if err := otel.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down OpenTelemetry: %v", err)
	}

	log.Println("👋 Server stopped")
}
