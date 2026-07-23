package middleware

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

type CustomClaims struct {
	SchoolID string   `json:"school_id"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

func JWTGuard(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    401,
			"error": fiber.Map{
				"message": "Missing authorization header",
				"type":    "auth_error",
			},
		})
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    401,
			"error": fiber.Map{
				"message": "Invalid authorization format",
				"type":    "auth_error",
			},
		})
	}

	token, err := jwt.ParseWithClaims(parts[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    401,
			"error": fiber.Map{
				"message": "Invalid token",
				"type":    "auth_error",
			},
		})
	}

	claims := token.Claims.(*CustomClaims)
	c.Locals("user_id", claims.Subject)
	c.Locals("school_id", claims.SchoolID)
	c.Locals("roles", claims.Roles)

	return c.Next()
}
