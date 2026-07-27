package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/violin-assessment/dawai/internal/db"
)

type AssessmentService struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewAssessmentService(pool *pgxpool.Pool) *AssessmentService {
	return &AssessmentService{
		queries: db.New(pool),
		pool:    pool,
	}
}

type ScoreInput struct {
	RubricComponentID string `json:"rubric_component_id"`
	Score             int32  `json:"score"`
}

type AssessmentWithComponents struct {
	Assessment db.Assessment            `json:"assessment"`
	Components []db.AssessmentComponent `json:"components"`
}

type AssessmentFilter struct {
	StudentID string
	TeacherID string
}

// Create creates an assessment with components in a transaction.
// Cross-tenant FK validation: student.SchoolID == subject.SchoolID == schoolID.
func (s *AssessmentService) Create(
	ctx context.Context,
	schoolID, teacherID, studentID, subjectID string,
	feedback string,
	scores []ScoreInput,
) (*AssessmentWithComponents, error) {
	schoolUUID := parseUUID(schoolID)
	studentUUID := parseUUID(studentID)
	subjectUUID := parseUUID(subjectID)
	teacherUUID := parseUUID(teacherID)

	// Cross-tenant FK validation: student belongs to this school.
	student, err := s.queries.GetStudentByID(ctx, studentUUID)
	if err != nil {
		return nil, fmt.Errorf("student not found: %w", err)
	}
	if !uuidEqual(student.SchoolID, schoolUUID) {
		return nil, errors.New("student does not belong to your school")
	}

	// Cross-tenant FK validation: subject belongs to this school.
	subject, err := s.queries.GetSubjectByID(ctx, subjectUUID)
	if err != nil {
		return nil, fmt.Errorf("subject not found: %w", err)
	}
	if !uuidEqual(subject.SchoolID, schoolUUID) {
		return nil, errors.New("subject does not belong to your school")
	}

	// Validate rubric components exist and belong to this school.
	for _, sc := range scores {
		rcUUID := parseUUID(sc.RubricComponentID)
		rc, err := s.queries.GetRubricComponentByID(ctx, rcUUID)
		if err != nil {
			return nil, fmt.Errorf("rubric component %s not found: %w", sc.RubricComponentID, err)
		}
		if !uuidEqual(rc.SchoolID, schoolUUID) {
			return nil, fmt.Errorf("rubric component %s does not belong to your school", sc.RubricComponentID)
		}
	}

	// Transaction: insert assessment + all components.
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	assessment, err := qtx.CreateAssessment(ctx, db.CreateAssessmentParams{
		SchoolID:  schoolUUID,
		SubjectID: subjectUUID,
		StudentID: studentUUID,
		TeacherID: teacherUUID,
		Feedback:  pgtype.Text{String: feedback, Valid: feedback != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create assessment: %w", err)
	}

	// Create each assessment component.
	components := make([]db.AssessmentComponent, 0, len(scores))
	for _, sc := range scores {
		comp, err := qtx.CreateAssessmentComponent(ctx, db.CreateAssessmentComponentParams{
			AssessmentID:      assessment.ID,
			RubricComponentID: parseUUID(sc.RubricComponentID),
			Score:             sc.Score,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create assessment component: %w", err)
		}
		components = append(components, comp)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &AssessmentWithComponents{
		Assessment: assessment,
		Components: components,
	}, nil
}

// List returns assessments filtered by school and optional student/teacher.
func (s *AssessmentService) List(
	ctx context.Context,
	schoolID string,
	filter AssessmentFilter,
) ([]db.Assessment, error) {
	schoolUUID := parseUUID(schoolID)

	// ponytail: scan all then filter in Go instead of dynamic SQL. Add query builder when >1000 assessments per school.
	all, err := s.queries.GetAssessmentsBySchool(ctx, schoolUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to list assessments: %w", err)
	}

	if filter.StudentID == "" && filter.TeacherID == "" {
		return all, nil
	}

	filtered := make([]db.Assessment, 0, len(all))
	for _, a := range all {
		if filter.StudentID != "" && !uuidEqual(a.StudentID, parseUUID(filter.StudentID)) {
			continue
		}
		if filter.TeacherID != "" && !uuidEqual(a.TeacherID, parseUUID(filter.TeacherID)) {
			continue
		}
		filtered = append(filtered, a)
	}
	return filtered, nil
}

// GetByID returns an assessment with all components.
// Cross-tenant check: assessment.SchoolID must match request schoolID.
func (s *AssessmentService) GetByID(
	ctx context.Context,
	assessmentID, schoolID string,
) (*AssessmentWithComponents, error) {
	aUUID := parseUUID(assessmentID)

	assessment, err := s.queries.GetAssessmentByID(ctx, aUUID)
	if err != nil {
		return nil, fmt.Errorf("assessment not found: %w", err)
	}

	// Cross-tenant check.
	if !uuidEqual(assessment.SchoolID, parseUUID(schoolID)) {
		return nil, errors.New("assessment not found in your school")
	}

	components, err := s.queries.GetAssessmentComponents(ctx, aUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assessment components: %w", err)
	}

	return &AssessmentWithComponents{
		Assessment: assessment,
		Components: components,
	}, nil
}

// Update modifies assessment scores and feedback. Only the teacher who created it,
// within 24h of submission.
func (s *AssessmentService) Update(
	ctx context.Context,
	assessmentID, teacherID string,
	scores []ScoreInput,
	feedback string,
	teacherSchoolID string,
) (*AssessmentWithComponents, error) {
	aUUID := parseUUID(assessmentID)
	tUUID := parseUUID(teacherID)

	assessment, err := s.queries.GetAssessmentByID(ctx, aUUID)
	if err != nil {
		return nil, fmt.Errorf("assessment not found: %w", err)
	}

	// Ownership check.
	if !uuidEqual(assessment.TeacherID, tUUID) {
		return nil, errors.New("you can only edit your own assessments")
	}

	// School match.
	if !uuidEqual(assessment.SchoolID, parseUUID(teacherSchoolID)) {
		return nil, errors.New("assessment not found in your school")
	}

	// 24h edit window (ponytail: simple time comparison; add configurable window when needed).
	if assessment.SubmittedAt.Valid {
		deadline := assessment.SubmittedAt.Time.Add(24 * time.Hour)
		if time.Now().After(deadline) {
			return nil, errors.New("assessment can only be edited within 24 hours of submission")
		}
	}

	// Validate rubric components belong to same school.
	for _, sc := range scores {
		rcUUID := parseUUID(sc.RubricComponentID)
		rc, err := s.queries.GetRubricComponentByID(ctx, rcUUID)
		if err != nil {
			return nil, fmt.Errorf("rubric component %s not found: %w", sc.RubricComponentID, err)
		}
		if !uuidEqual(rc.SchoolID, parseUUID(teacherSchoolID)) {
			return nil, fmt.Errorf("rubric component %s does not belong to your school", sc.RubricComponentID)
		}
	}

	// Transaction: delete old components, insert new ones, update feedback.
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	if err := qtx.DeleteAssessmentComponents(ctx, aUUID); err != nil {
		return nil, fmt.Errorf("failed to delete old components: %w", err)
	}

	for _, sc := range scores {
		_, err := qtx.CreateAssessmentComponent(ctx, db.CreateAssessmentComponentParams{
			AssessmentID:      aUUID,
			RubricComponentID: parseUUID(sc.RubricComponentID),
			Score:             sc.Score,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create assessment component: %w", err)
		}
	}

	if err := qtx.UpdateAssessment(ctx, db.UpdateAssessmentParams{
		ID:       aUUID,
		Feedback: pgtype.Text{String: feedback, Valid: feedback != ""},
	}); err != nil {
		return nil, fmt.Errorf("failed to update assessment: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Re-fetch the updated assessment with components.
	return s.GetByID(ctx, assessmentID, teacherSchoolID)
}

// Delete removes an assessment. Only the teacher who created it (or admin) can delete.
func (s *AssessmentService) Delete(
	ctx context.Context,
	assessmentID, teacherID, teacherSchoolID string,
) error {
	aUUID := parseUUID(assessmentID)
	tUUID := parseUUID(teacherID)

	assessment, err := s.queries.GetAssessmentByID(ctx, aUUID)
	if err != nil {
		return fmt.Errorf("assessment not found: %w", err)
	}

	// School match.
	if !uuidEqual(assessment.SchoolID, parseUUID(teacherSchoolID)) {
		return errors.New("assessment not found in your school")
	}

	// Ownership check.
	if !uuidEqual(assessment.TeacherID, tUUID) {
		return errors.New("you can only delete your own assessments")
	}

	// Delete assessment (cascade deletes components via FK).
	return s.queries.DeleteAssessment(ctx, aUUID)
}

// -- helpers --

func parseUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	_ = u.Scan(s)
	return u
}

func uuidEqual(a, b pgtype.UUID) bool {
	if !a.Valid || !b.Valid {
		return false
	}
	if len(a.Bytes) != len(b.Bytes) {
		return false
	}
	for i := range a.Bytes {
		if a.Bytes[i] != b.Bytes[i] {
			return false
		}
	}
	return true
}

func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	// Format: 8-4-4-4-12
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
