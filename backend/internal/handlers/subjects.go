package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/violin-assessment/dawai/internal/models"
	"github.com/violin-assessment/dawai/internal/services"
)

type SubjectHandler struct {
	svc *services.SubjectService
}

func NewSubjectHandler(svc *services.SubjectService) *SubjectHandler {
	return &SubjectHandler{svc: svc}
}

type createSubjectReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *SubjectHandler) GetSubjects(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	subjects, err := h.svc.ListBySchool(c.Context(), schoolID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false, Code: 500,
			Error: &models.ErrorBody{Message: err.Error(), Type: "internal_error"},
		})
	}
	return c.JSON(models.Response{Success: true, Code: 200, Data: subjects})
}

func (h *SubjectHandler) CreateSubject(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	var req createSubjectReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false, Code: 400,
			Error: &models.ErrorBody{Message: "Invalid request body", Type: "validation_error"},
		})
	}
	subj, err := h.svc.Create(c.Context(), schoolID, req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false, Code: 400,
			Error: &models.ErrorBody{Message: err.Error(), Type: "validation_error"},
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true, Code: 201, Data: subj,
	})
}

func (h *SubjectHandler) DeleteSubject(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	subjectID := c.Params("id")
	if err := h.svc.Delete(c.Context(), subjectID, schoolID); err != nil {
		msg := err.Error()
		status := fiber.StatusBadRequest
		// ponytail: string match on known error messages, typed errors if this grows >3
		if msg == "cannot delete subject with existing assessments" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(models.Response{
			Success: false, Code: status,
			Error: &models.ErrorBody{Message: msg, Type: "validation_error"},
		})
	}
	return c.JSON(models.Response{Success: true, Code: 200, Data: fiber.Map{"deleted": true}})
}

func (h *SubjectHandler) GetRubric(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	subjectID := c.Params("id")
	components, err := h.svc.GetRubric(c.Context(), subjectID, schoolID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "subject not found" || err.Error() == "subject not found in your school" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(models.Response{
			Success: false, Code: status,
			Error: &models.ErrorBody{Message: err.Error(), Type: "not_found"},
		})
	}
	return c.JSON(models.Response{Success: true, Code: 200, Data: components})
}

type createRubricReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ScaleMin    int32  `json:"scale_min"`
	ScaleMax    int32  `json:"scale_max"`
	Weight      int32  `json:"weight"`
}

func (h *SubjectHandler) CreateRubricComponent(c *fiber.Ctx) error {
	schoolID := c.Locals("school_id").(string)
	subjectID := c.Params("id")
	var req createRubricReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false, Code: 400,
			Error: &models.ErrorBody{Message: "Invalid request body", Type: "validation_error"},
		})
	}
	comp, err := h.svc.AddRubricComponent(c.Context(), schoolID, subjectID,
		req.Name, req.Description, req.ScaleMin, req.ScaleMax, req.Weight)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "subject not found" || err.Error() == "subject not found in your school" {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(models.Response{
			Success: false, Code: status,
			Error: &models.ErrorBody{Message: err.Error(), Type: "validation_error"},
		})
	}
	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true, Code: 201, Data: comp,
	})
}
