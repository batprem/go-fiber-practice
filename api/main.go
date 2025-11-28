package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	otel "gfp/api/middlewares/otel"
	"gfp/api/routes"
	"gfp/api/routes/nested"
	"gfp/lib"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
)

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

	// Define routes
	app.Get(
		"/",
		func(c *fiber.Ctx) error {
			ctx := c.UserContext()
			otel.LogInfo(ctx, "Handling root endpoint")
			return c.SendString("Hello, World!")
		},
	)

	app.Get(
		"/:user",
		func(c *fiber.Ctx) error {
			ctx := c.UserContext()
			user := c.Params("user")
			otel.LogInfo(ctx, "Handling user endpoint", "user", user)
			return c.SendString(lib.GetText(user))
		},
	)

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
