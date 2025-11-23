package nested

import (
	"gfp/lib"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	GreetingText string `json:"greetingText"`
}

func NestedTextReturn (c *fiber.Ctx) error {
	userName := c.Params("user")
	user := User { 
		ID: 1,
		Name: userName,
		GreetingText: lib.GetText(userName),
	}
	return c.JSON(user)
}