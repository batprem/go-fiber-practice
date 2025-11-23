package routes

import (
	"gfp/lib"
	"github.com/gofiber/fiber/v2"
)

func TextReturn2 (c *fiber.Ctx) error {
	return c.SendString(lib.GetText(c.Params("user")))
}