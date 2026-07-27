package middleware

import (
	"fmt"
	"os"
	"strings"
	"github.com/violin-assessment/dawai/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

type CustomClaims struct {
	SchoolID string   `json:"school_id"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

func NewJWTGuard(queries *db.Queries) fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	secretBytes := []byte(secret)

	return func(c *fiber.Ctx) error {
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
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretBytes, nil
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

	// Validate required claims: subject (user_id) must be present and non-empty
	if claims.Subject == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    401,
			"error": fiber.Map{
				"message": "Invalid token: missing or empty subject",
				"type":    "auth_error",
			},
		})
	}

	// Validate required claims: school_id must be present and non-empty
	if claims.SchoolID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"code":    401,
			"error": fiber.Map{
				"message": "Invalid token: missing or empty school_id",
				"type":    "auth_error",
			},
		})
	}

	// Inject validated claims into context locals
	c.Locals("user_id", claims.Subject)
	c.Locals("school_id", claims.SchoolID)
	c.Locals("roles", claims.Roles)
	c.Locals("jti", claims.ID)

	// Check if token is blacklisted (jti is required from token generation)
	if claims.ID != "" {
		_, err := queries.IsJWTBlacklisted(c.Context(), claims.ID)
		if err == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"code":    401,
				"error": fiber.Map{
					"message": "Token revoked",
					"type":    "auth_error",
				},
			})
		}
	}

	return c.Next()
	}
}
