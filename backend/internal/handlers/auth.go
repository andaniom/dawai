package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/violin-assessment/dawai/internal/models"
	"github.com/violin-assessment/dawai/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Code:    400,
			Error: &models.ErrorBody{
				Message: "Invalid request body",
				Type:    "validation_error",
			},
		})
	}

	// TODO: Fetch user from database, verify password, issue token
	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success: true,
		Code:    200,
		Data: models.LoginResponse{
			AccessToken: "placeholder",
			User: models.User{
				ID:    "user_id",
				Email: req.Email,
				Name:  "User Name",
				Roles: []string{"teacher"},
			},
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// TODO: Blacklist JWT
	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success: true,
		Code:    200,
		Data:    fiber.Map{"message": "logged out"},
	})
}
