package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/violin-assessment/dawai/internal/db"
)

type SubjectService struct {
	queries *db.Queries
}

func NewSubjectService(queries *db.Queries) *SubjectService {
	return &SubjectService{queries: queries}
}

func (s *SubjectService) ListBySchool(ctx context.Context, schoolID string) ([]db.Subject, error) {
	return s.queries.GetSubjectsBySchool(ctx, parseUUID(schoolID))
}

func (s *SubjectService) Create(ctx context.Context, schoolID, name, description string) (db.Subject, error) {
	return s.queries.CreateSubject(ctx, db.CreateSubjectParams{
		SchoolID:    parseUUID(schoolID),
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
	})
}

func (s *SubjectService) Delete(ctx context.Context, subjectID, schoolID string) error {
	sid := parseUUID(subjectID)
	schoolUUID := parseUUID(schoolID)

	// verify subject belongs to school
	subj, err := s.queries.GetSubjectByID(ctx, sid)
	if err != nil {
		return fmt.Errorf("subject not found")
	}
	if !uuidEqual(subj.SchoolID, schoolUUID) {
		return fmt.Errorf("subject not found in your school")
	}

	// check no assessments reference this subject
	count, err := s.queries.CountAssessmentsBySubject(ctx, sid)
	if err != nil {
		return fmt.Errorf("failed to check assessments: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete subject with existing assessments")
	}

	return s.queries.DeleteSubject(ctx, sid)
}

func (s *SubjectService) GetRubric(ctx context.Context, subjectID, schoolID string) ([]db.RubricComponent, error) {
	sid := parseUUID(subjectID)
	schoolUUID := parseUUID(schoolID)

	subj, err := s.queries.GetSubjectByID(ctx, sid)
	if err != nil {
		return nil, fmt.Errorf("subject not found")
	}
	if !uuidEqual(subj.SchoolID, schoolUUID) {
		return nil, fmt.Errorf("subject not found in your school")
	}

	return s.queries.GetRubricComponentsBySubject(ctx, sid)
}

func (s *SubjectService) AddRubricComponent(
	ctx context.Context, schoolID, subjectID, name, description string,
	scaleMin, scaleMax, weight int32,
) (db.RubricComponent, error) {
	schoolUUID := parseUUID(schoolID)
	subjUUID := parseUUID(subjectID)

	// cross-tenant: verify subject belongs to school
	subj, err := s.queries.GetSubjectByID(ctx, subjUUID)
	if err != nil {
		return db.RubricComponent{}, fmt.Errorf("subject not found")
	}
	if !uuidEqual(subj.SchoolID, schoolUUID) {
		return db.RubricComponent{}, fmt.Errorf("subject not found in your school")
	}

	return s.queries.CreateRubricComponent(ctx, db.CreateRubricComponentParams{
		SchoolID:    schoolUUID,
		SubjectID:   subjUUID,
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
		ScaleMin:    pgtype.Int4{Int32: scaleMin, Valid: true},
		ScaleMax:    pgtype.Int4{Int32: scaleMax, Valid: true},
		Weight:      pgtype.Int4{Int32: weight, Valid: true},
	})
}
