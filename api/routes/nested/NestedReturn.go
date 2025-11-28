package nested

import (
	otel "gfp/api/middlewares/otel"
	"gfp/lib"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	GreetingText string `json:"greetingText"`
}

func NestedTextReturn(c *fiber.Ctx) error {
	// Get context from Fiber
	ctx := c.UserContext()

	// Create a span for this handler
	tracer := otel.Tracer()
	ctx, span := tracer.Start(ctx, "routes.nested.NestedTextReturn")
	defer span.End()

	// Get user parameter
	userName := c.Params("user")

	// Add span attributes
	span.SetAttributes(
		attribute.String("user.name", userName),
		attribute.String("http.route", "/simple-return/nested/:user"),
	)

	// Log with trace context
	otel.LogInfo(ctx, "Processing NestedTextReturn request", "user", userName)

	// Create user object
	user := User{
		ID:           1,
		Name:         userName,
		GreetingText: lib.GetText(userName),
	}

	// Add more span attributes
	span.SetAttributes(
		attribute.Int("user.id", user.ID),
	)

	// Set span status
	span.SetStatus(codes.Ok, "Request processed successfully")

	return c.JSON(user)
}
