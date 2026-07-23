package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/violin-assessment/dawai/internal/models"
	"github.com/violin-assessment/dawai/internal/services"
)

type UserHandler struct {
	authService *services.AuthService
}

func NewUserHandler(authService *services.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
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

	user, err := h.authService.CreateUser(c.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Code:    400,
			Error: &models.ErrorBody{
				Message: err.Error(),
				Type:    "validation_error",
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Code:    201,
		Data:    user,
	})
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	roleFilter := c.Query("role")

	users, err := h.authService.ListUsers(c.Context(), schoolID, roleFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Code:    500,
			Error: &models.ErrorBody{
				Message: err.Error(),
				Type:    "server_error",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success: true,
		Code:    200,
		Data:    users,
	})
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := h.authService.GetMe(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Code:    404,
			Error: &models.ErrorBody{
				Message: err.Error(),
				Type:    "not_found",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success: true,
		Code:    200,
		Data:    user,
	})
}
