package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/violin-assessment/dawai/internal/models"
	"github.com/violin-assessment/dawai/internal/services"
)

type StudentHandler struct {
	svc *services.StudentService
}

func NewStudentHandler(svc *services.StudentService) *StudentHandler {
	return &StudentHandler{svc: svc}
}

// ListStudents godoc
// GET /api/students — list students by school, optional ?class filter.
// Students can only see their own record.
func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	userID := c.Locals("user_id").(string)
	roles := c.Locals("roles").([]string)

	// Student: scope to own record only
	for _, r := range roles {
		if r == "student" {
			student, err := h.svc.GetByUserID(c.Context(), schoolID, userID)
			if err != nil {
				return c.JSON(models.Response{Success: true, Code: 200, Data: []interface{}{}})
			}
			return c.JSON(models.Response{Success: true, Code: 200, Data: []interface{}{student}})
		}
	}

	classFilter := c.Query("class")
	students, err := h.svc.ListBySchool(c.Context(), schoolID, classFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false, Code: 500,
			Error: &models.ErrorBody{Message: err.Error(), Type: "server_error"},
		})
	}
	return c.JSON(models.Response{Success: true, Code: 200, Data: students})
}

// GetStudent godoc
// GET /api/students/:id — single student with cross-tenant check.
func (h *StudentHandler) GetStudent(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	studentID := c.Params("id")

	student, err := h.svc.GetByID(c.Context(), schoolID, studentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false, Code: 404,
			Error: &models.ErrorBody{Message: err.Error(), Type: "not_found"},
		})
	}
	return c.JSON(models.Response{Success: true, Code: 200, Data: student})
}

// CreateStudent godoc
// POST /api/students — create student record, user must exist in same school.
func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)

	var req services.CreateStudentInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false, Code: 400,
			Error: &models.ErrorBody{Message: "Invalid request body", Type: "validation_error"},
		})
	}

	student, err := h.svc.Create(c.Context(), schoolID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false, Code: 400,
			Error: &models.ErrorBody{Message: err.Error(), Type: "validation_error"},
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.Response{Success: true, Code: 201, Data: student})
}
