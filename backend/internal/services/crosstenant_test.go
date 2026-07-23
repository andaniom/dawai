package services_test

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/violin-assessment/dawai/internal/services"
)

const skipEnv = "DAWAI_SKIP_DB_TESTS"

func skipIfNoDB(t *testing.T) {
	t.Helper()
	if v := os.Getenv(skipEnv); v != "0" {
		t.Skip("skipping DB-backed cross-tenant test (set DAWAI_SKIP_DB_TESTS=0 to run)")
	}
}

// Compile-time harness: asserts services instantiate against real signatures.
func TestHarnessBuilds(t *testing.T) {
	pool, err := pgxpool.New(t.Context(), "")
	if err != nil {
		t.Skip("no DB connection")
	}
	defer pool.Close()

	_ = services.NewAssessmentService(pool)
	_ = services.NewStudentService(nil)
	_ = services.NewSubjectService(nil)
}

// AssessmentService.Create rejects student belonging to a different school.
func TestAssessmentCreateRejectsStudentWrongSchool(t *testing.T) {
	skipIfNoDB(t)
	require.Fail(t, "DB seed wiring required: seed studentB in schoolB, then assert AssessmentService.Create(ctx, schoolA, ...) with studentB.ID returns 'student does not belong to your school'")
}

// AssessmentService.Create rejects subject belonging to a different school.
func TestAssessmentCreateRejectsSubjectWrongSchool(t *testing.T) {
	skipIfNoDB(t)
	require.Fail(t, "DB seed wiring required: seed subjectB in schoolB, then assert AssessmentService.Create(ctx, schoolA, ..., subjectB.ID, ...) returns 'subject does not belong to your school'")
}

// SubjectService.AddRubricComponent rejects subject belonging to a different school.
func TestRubricComponentCreateRejectsSubjectWrongSchool(t *testing.T) {
	skipIfNoDB(t)
	require.Fail(t, "DB seed wiring required: seed subjectB in schoolB, then assert SubjectService.AddRubricComponent(ctx, schoolA, subjectB.ID, ...) returns 'subject not found in your school'")
}

// StudentService.Create rejects user belonging to a different school.
func TestStudentCreateRejectsUserWrongSchool(t *testing.T) {
	skipIfNoDB(t)
	require.Fail(t, "DB seed wiring required: seed userB in schoolB (via user_roles), then assert StudentService.Create(ctx, schoolA, {UserID: userB.ID}) returns 'user not found in your school'")
}
