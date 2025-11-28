package routes

import (
	otel "gfp/api/middlewares/otel"
	"gfp/lib"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TextReturn(c *fiber.Ctx) error {
	// Get context from Fiber
	ctx := c.UserContext()

	// Create a span for this handler
	tracer := otel.Tracer()
	ctx, span := tracer.Start(ctx, "routes.TextReturn")
	defer span.End()

	// Get user parameter
	user := c.Params("user")

	// Add span attributes
	span.SetAttributes(
		attribute.String("user.name", user),
		attribute.String("http.route", "/simple-return/:user"),
	)

	// Log with trace context
	otel.LogInfo(ctx, "Processing TextReturn request", "user", user)

	// Get text from lib
	text := lib.GetText(user)

	// Set span status
	span.SetStatus(codes.Ok, "Request processed successfully")

	return c.SendString(text)
}
