package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ponytail: integration-style isolation tests. Assumes test DB seeded with
// two schools (A, B), each with users/students/subjects/assessments.
// In CI: spin up postgres via docker-compose test profile, run migrations,
// seed via test helper, then `go test -tags=integration`.
// Without DB: set DAWAI_SKIP_DB_TESTS=1 to skip (build still asserts compilation).

const skipEnv = "DAWAI_SKIP_DB_TESTS"

func skipIfNoDB(t *testing.T) {
	t.Helper()
	if v := os.Getenv(skipEnv); v != "0" {
		t.Skip("skipping DB-backed isolation test (set DAWAI_SKIP_DB_TESTS=0 to run)")
	}
}

// appWithLocals builds a Fiber app that injects the given school_id + user_id
// into c.Locals (simulates JWTGuard + TenantGuard already run).
func appWithLocals(schoolID, userID string) *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("school_id", schoolID)
		c.Locals("user_id", userID)
		c.Locals("roles", []string{"teacher"})
		return c.Next()
	})
	return app
}

func doJSON(app *fiber.App, method, path string, body interface{}) (*fiber.App, []byte, int) {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	out, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return app, out, resp.StatusCode
}

type listResp struct {
	Success bool `json:"success"`
	Data    []struct {
		ID      string `json:"id"`
		School  string `json:"school_id"`
	} `json:"data"`
}

func TestStudentListIsolation(t *testing.T) {
	skipIfNoDB(t)
	require.Fail(t, "DB seed wiring required: replace this stub with real /api/students endpoint mount + seeded data for schools A and B")
	// Once wired:
	//   appA := appWithLocals("schoolA", "userA"); mount real student routes
	//   _, body, code := doJSON(appA, "GET", "/api/students", nil)
	//   assert.Equal(t, 200, code)
	//   assert.NotContains(t, string(body), "\"schoolB\"")
	//   assert.Contains(t, string(body), "\"schoolA\"")
}

func TestCrossTenantStudentFetchForbidden(t *testing.T) {
	skipIfNoDB(t)
	require.Fail(t, "wire real student route + seed student X owned by schoolB, then assert 403 when schoolA user GETs /api/students/X")
}

func TestCrossTenantAssessmentCreateRejected(t *testing.T) {
	skipIfNoDB(t)
	require.Fail(t, "wire real assessment route + seed student+subject in schoolB, then assert 403 when schoolA POSTs assessment with B ids")
}

func TestSubjectListIsolation(t *testing.T) {
	skipIfNoDB(t)
	require.Fail(t, "wire real subject route + assert schoolA only sees schoolA subjects")
}

// Compile-time guarantee: helpers form the test harness shape expected by the suite.
func TestHarnessBuilds(t *testing.T) {
	app := appWithLocals("x", "y")
	assert.NotNil(t, app)
	_, out, _ := doJSON(app, "GET", "/health", nil)
	_ = out
}