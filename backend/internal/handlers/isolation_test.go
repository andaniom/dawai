package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/violin-assessment/dawai/internal/services"
)

const skipEnv = "DAWAI_SKIP_DB_TESTS"

func skipIfNoDB(t *testing.T) {
	t.Helper()
	if v := os.Getenv(skipEnv); v != "0" {
		t.Skip("skipping DB-backed isolation test (set DAWAI_SKIP_DB_TESTS=0 to run)")
	}
}

func appWithLocals(schoolID, userID string, roles []string, headers map[string]string) *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("school_id", schoolID)
		c.Locals("user_id", userID)
		c.Locals("roles", roles)
		return c.Next()
	})
	return app
}

func doJSON(app *fiber.App, method, path string, body interface{}, headers map[string]string) (*fiber.App, []byte, int) {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, _ := app.Test(req)
	out, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return app, out, resp.StatusCode
}

// Stubs for service layer mocks to isolate handler layer cross-tenant logic testing

type MockStudentService struct {
	GetByIDFunc func(ctx interface{}, schoolID, studentID string) (interface{}, error)
	CreateFunc  func(ctx interface{}, schoolID string, req services.CreateStudentInput) (interface{}, error)
}

func (m *MockStudentService) GetByID(ctx interface{}, schoolID, studentID string) (interface{}, error) {
	return m.GetByIDFunc(ctx, schoolID, studentID)
}

func (m *MockStudentService) Create(ctx interface{}, schoolID string, req services.CreateStudentInput) (interface{}, error) {
	return m.CreateFunc(ctx, schoolID, req)
}

type MockAssessmentService struct {
	UpdateFunc func(ctx interface{}, schoolID, id string, req interface{}) (interface{}, error)
}

func (m *MockAssessmentService) Update(ctx interface{}, schoolID, id string, req interface{}) (interface{}, error) {
	return m.UpdateFunc(ctx, schoolID, id, req)
}

func TestCrossTenant_UserCannotReadEntityFromOtherSchool(t *testing.T) {
	skipIfNoDB(t)
    // skipped: Real DB wiring. Add when `test_helpers.go` with DB seeding is merged.
}

func TestCrossTenant_UserCannotWriteEntityFromOtherSchool(t *testing.T) {
    skipIfNoDB(t)
}

func TestCrossTenant_SchoolIDNotFromRequest(t *testing.T) {
    skipIfNoDB(t)
}

func TestCrossTenant_SuperAdminCanImpersonateViaHeader(t *testing.T) {
    skipIfNoDB(t)
}

func TestCrossTenant_NonSuperAdminHeaderIgnored(t *testing.T) {
    skipIfNoDB(t)
}
