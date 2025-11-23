package main

import (
	"gfp/api/routes"
	"gfp/api/routes/nested"
	"gfp/lib"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get(
		"/",
		func(c *fiber.Ctx) error {
			return c.SendString("Hello, World!")
		},
	)
	app.Get(
		"/:user",
		func(c *fiber.Ctx) error {
			return c.SendString(lib.GetText(c.Params("user")))
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

	log.Fatal(app.Listen(":3002"))
}
