package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/violin-assessment/dawai/internal/models"
	"github.com/violin-assessment/dawai/internal/services"
)

type AssessmentHandler struct {
	service *services.AssessmentService
}

func NewAssessmentHandler(service *services.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{service: service}
}

type CreateAssessmentRequest struct {
	StudentID string                 `json:"student_id"`
	SubjectID string                 `json:"subject_id"`
	Scores    []services.ScoreInput  `json:"scores"`
	Feedback  string                 `json:"feedback"`
}

type UpdateAssessmentRequest struct {
	Scores   []services.ScoreInput `json:"scores"`
	Feedback string                `json:"feedback"`
}

// Create godoc
// POST /api/assessments — teacher submits an assessment.
func (h *AssessmentHandler) Create(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	teacherID := c.Locals("user_id").(string)

	var req CreateAssessmentRequest
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

	result, err := h.service.Create(c.Context(), schoolID, teacherID, req.StudentID, req.SubjectID, req.Feedback, req.Scores)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(models.Response{
			Success: false,
			Code:    403,
			Error: &models.ErrorBody{
				Message: err.Error(),
				Type:    "validation_error",
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Code:    201,
		Data:    result,
	})
}

// List godoc
// GET /api/assessments — list assessments scoped by role.
func (h *AssessmentHandler) List(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	userID := c.Locals("user_id").(string)
	roles := c.Locals("roles").([]string)

	filter := services.AssessmentFilter{
		StudentID: c.Query("student_id"),
	}

	// Student: only see own assessments.
	for _, r := range roles {
		if r == "student" {
			filter.StudentID = userIDDir(userID)
			break
		}
	}

	assessments, err := h.service.List(c.Context(), schoolID, filter)
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
		Data:    assessments,
	})
}

// GetByID godoc
// GET /api/assessments/:id — single assessment with component breakdown.
func (h *AssessmentHandler) GetByID(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)

	result, err := h.service.GetByID(c.Context(), c.Params("id"), schoolID)
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
		Data:    result,
	})
}

// Update godoc
// PATCH /api/assessments/:id — edit own assessment within 24h window.
func (h *AssessmentHandler) Update(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	teacherID := c.Locals("user_id").(string)

	var req UpdateAssessmentRequest
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

	result, err := h.service.Update(c.Context(), c.Params("id"), teacherID, req.Scores, req.Feedback, schoolID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(models.Response{
			Success: false,
			Code:    403,
			Error: &models.ErrorBody{
				Message: err.Error(),
				Type:    "validation_error",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success: true,
		Code:    200,
		Data:    result,
	})
}

// Delete godoc
// DELETE /api/assessments/:id — teacher deletes own assessment.
func (h *AssessmentHandler) Delete(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	teacherID := c.Locals("user_id").(string)

	if err := h.service.Delete(c.Context(), c.Params("id"), teacherID, schoolID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(models.Response{
			Success: false,
			Code:    403,
			Error: &models.ErrorBody{
				Message: err.Error(),
				Type:    "validation_error",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Success: true,
		Code:    200,
		Data:    fiber.Map{"deleted": true},
	})
}

// userIDDir returns the user_id that corresponds to a student lookup.
// ponytail: stub — student user_id lookup needs student table. Replace when GetStudentUserID query exists.
func userIDDir(_ string) string {
	return ""
}
