package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/violin-assessment/dawai/internal/db"
	"github.com/violin-assessment/dawai/internal/handlers"
	"github.com/violin-assessment/dawai/internal/middleware"
	"github.com/violin-assessment/dawai/internal/services"
)

func main() {
	// Validate critical env vars before initialization
	if secret := os.Getenv("JWT_SECRET"); secret == "" {
		panic("JWT_SECRET env var not set — tokens will forge trivially")
	}

	dbURL := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	queries := db.New(pool)

	// Services
	authService := services.NewAuthService(queries)
	assessmentService := services.NewAssessmentService(pool)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)
	assessmentHandler := handlers.NewAssessmentHandler(assessmentService)
	studentService := services.NewStudentService(queries)
	studentHandler := handlers.NewStudentHandler(studentService)
	subjectService := services.NewSubjectService(queries)
	subjectHandler := handlers.NewSubjectHandler(subjectService)

	app := fiber.New()

	// Public routes
	app.Post("/api/auth/login", middleware.RateLimiter(), authHandler.Login)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Protected routes
	authGroup := app.Group("/api/auth", middleware.NewJWTGuard(queries))
	authGroup.Post("/logout", authHandler.Logout)

	api := app.Group("/api", middleware.NewJWTGuard(queries), middleware.TenantGuard)

	// /api/me — current user
	api.Get("/me", userHandler.GetMe)

	// Users (admin only)
	usersGroup := api.Group("/users", middleware.RoleGuard("school_admin"))
	usersGroup.Post("", userHandler.CreateUser)
	usersGroup.Get("", userHandler.ListUsers)

	// Students
	studentsGroup := api.Group("/students")
	studentsGroup.Get("", studentHandler.ListStudents)
	studentsGroup.Get("/:id", studentHandler.GetStudent)
	studentsGroup.Post("", studentHandler.CreateStudent)

	// Subjects
	subjectsGroup := api.Group("/subjects")
	subjectsGroup.Get("", subjectHandler.GetSubjects)
	subjectsGroup.Post("", subjectHandler.CreateSubject)
	subjectsGroup.Delete("/:id", subjectHandler.DeleteSubject)
	subjectsGroup.Get("/:id/rubric", subjectHandler.GetRubric)
	subjectsGroup.Post("/:id/rubric", subjectHandler.CreateRubricComponent)

	// Assessments
	assessmentsGroup := api.Group("/assessments")
	assessmentsGroup.Post("", assessmentHandler.Create)
	assessmentsGroup.Get("", assessmentHandler.List)
	assessmentsGroup.Get("/:id", assessmentHandler.GetByID)
	assessmentsGroup.Patch("/:id", assessmentHandler.Update)
	assessmentsGroup.Delete("/:id", assessmentHandler.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.Listen(":" + port)
}
