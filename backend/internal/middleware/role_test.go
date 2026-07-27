package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func makeApp(role string) *fiber.App {
	app := fiber.New()
	app.Post("/test", RoleGuard(role), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	return app
}

func doReq(app *fiber.App, roles []string) int {
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	app2 := fiber.New()
	app2.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("roles", roles)
		return c.Next()
	}, RoleGuard("school_admin"), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	resp, _ := app2.Test(req)
	io.ReadAll(resp.Body)
	return resp.StatusCode
}

func TestRoleGuard_blocks_wrong_role(t *testing.T) {
	if got := doReq(nil, []string{"teacher"}); got != 403 {
		t.Errorf("teacher → want 403 got %d", got)
	}
}

func TestRoleGuard_blocks_student(t *testing.T) {
	if got := doReq(nil, []string{"student"}); got != 403 {
		t.Errorf("student → want 403 got %d", got)
	}
}

func TestRoleGuard_allows_school_admin(t *testing.T) {
	if got := doReq(nil, []string{"school_admin"}); got != 200 {
		t.Errorf("school_admin → want 200 got %d", got)
	}
}

func TestRoleGuard_allows_super_admin(t *testing.T) {
	if got := doReq(nil, []string{"super_admin"}); got != 200 {
		t.Errorf("super_admin → want 200 got %d", got)
	}
}

func TestRoleGuard_nil_roles(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	app := fiber.New()
	app.Post("/test", RoleGuard("school_admin"), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	resp, _ := app.Test(req)
	io.ReadAll(resp.Body)
	if resp.StatusCode != 403 {
		t.Errorf("nil roles → want 403 got %d", resp.StatusCode)
	}
}
