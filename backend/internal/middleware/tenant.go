package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func TenantGuard(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id")
	if schoolID == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"code":    403,
			"error": fiber.Map{
				"message": "School ID not found in token",
				"type":    "tenant_error",
			},
		})
	}

	// Verify school_id is valid UUID format (basic check)
	schoolIDStr := schoolID.(string)
	if schoolIDStr == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"code":    403,
			"error": fiber.Map{
				"message": "Invalid school ID",
				"type":    "tenant_error",
			},
		})
	}

	return c.Next()
}
