package handlers

import (
	"net/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/violin-assessment/dawai/internal/services"
)

type mockAuthServiceForLogout struct{}

func (m *mockAuthServiceForLogout) Logout(ctx interface{}, jti string) error {
	return nil // Success mock
}

func TestLogoutWithValidJTI(t *testing.T) {
	app := fiber.New()
	handler := &AuthHandler{}

	app.Get("/logout", func(c *fiber.Ctx) error {
		// Simulate middleware injection of valid jti
		c.Locals("jti", "valid-jti-123")
		return handler.Logout(c)
	})

	req := httptest.NewRequest("GET", "/logout", nil)
	resp, _ := app.Test(req)

	// Should NOT panic even though mock doesn't blacklist
	if resp.StatusCode == 500 {
		t.Errorf("Unexpected 500 error during logout with valid jti")
	}
}

func TestLogoutWithMissingJTI(t *testing.T) {
	app := fiber.New()
	handler := &AuthHandler{}

	app.Get("/logout", func(c *fiber.Ctx) error {
		// Simulate middleware NOT injecting jti (should not panic)
		return handler.Logout(c)
	})

	req := httptest.NewRequest("GET", "/logout", nil)
	resp, _ := app.Test(req)

	// Should return 401, not panic
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected 401 for missing jti, got %d", resp.StatusCode)
	}
}

func TestLogoutWithWrongJTIType(t *testing.T) {
	app := fiber.New()
	handler := &AuthHandler{}

	app.Get("/logout", func(c *fiber.Ctx) error {
		// Simulate middleware injecting wrong type (should not panic)
		c.Locals("jti", 12345) // int instead of string
		return handler.Logout(c)
	})

	req := httptest.NewRequest("GET", "/logout", nil)
	resp, _ := app.Test(req)

	// Should return 401, not panic
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong jti type, got %d", resp.StatusCode)
	}
}
