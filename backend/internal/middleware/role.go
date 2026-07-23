package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func RoleGuard(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rolesInterface := c.Locals("roles")
		if rolesInterface == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"code":    403,
				"error": fiber.Map{
					"message": "Roles not found in token",
					"type":    "authorization_error",
				},
			})
		}

		roles := rolesInterface.([]string)
		for _, role := range roles {
			if role == requiredRole || role == "super_admin" {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"code":    403,
			"error": fiber.Map{
				"message": "Insufficient permissions",
				"type":    "authorization_error",
			},
		})
	}
}
