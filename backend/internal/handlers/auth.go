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

	result, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success: false,
			Code:    401,
			Error: &models.ErrorBody{
				Message: err.Error(),
				Type:    "auth_error",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success: true,
		Code:    200,
		Data: models.LoginResponse{
			AccessToken: result.Token,
			User: models.User{
				ID:       result.User.ID,
				Email:    result.User.Email,
				Name:     result.User.Name,
				Roles:    result.User.Roles,
				SchoolID: result.User.SchoolID,
			},
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Safe extraction of jti from context locals with type check
	jtiVal := c.Locals("jti")
	if jtiVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success: false,
			Code:    401,
			Error: &models.ErrorBody{
				Message: "Invalid token: missing jti claim",
				Type:    "auth_error",
			},
		})
	}

	jti, ok := jtiVal.(string)
	if !ok || jti == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success: false,
			Code:    401,
			Error: &models.ErrorBody{
				Message: "Invalid token: malformed jti claim",
				Type:    "auth_error",
			},
		})
	}

	if err := h.authService.Logout(c.Context(), jti); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Code:    500,
			Error: &models.ErrorBody{
				Message: "Failed to blacklist token",
				Type:    "server_error",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success: true,
		Code:    200,
		Data:    fiber.Map{"message": "logged out"},
	})
}
