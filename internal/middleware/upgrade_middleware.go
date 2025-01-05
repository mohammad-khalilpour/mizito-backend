package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func UpgradeMiddleware(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": "id parameter is not provided",
		})
	}
	return c.Next()
}
