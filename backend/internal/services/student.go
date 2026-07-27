package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/violin-assessment/dawai/internal/db"
)

type StudentService struct {
	queries *db.Queries
}

func NewStudentService(queries *db.Queries) *StudentService {
	return &StudentService{queries: queries}
}

type StudentResponse struct {
	ID        string `json:"id"`
	SchoolID  string `json:"school_id"`
	UserID    string `json:"user_id"`
	ClassName string `json:"class"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

type CreateStudentInput struct {
	UserID    string `json:"user_id"`
	ClassName string `json:"class"`
}

func (s *StudentService) ListBySchool(ctx context.Context, schoolID string, classFilter string) ([]StudentResponse, error) {
	schoolUUID := pgtype.UUID{}
	if err := schoolUUID.Scan(schoolID); err != nil {
		return nil, errors.New("invalid school_id")
	}

	rows, err := s.queries.ListStudentsBySchool(ctx, schoolUUID)
	if err != nil {
		return nil, err
	}

	out := make([]StudentResponse, 0, len(rows))
	for _, r := range rows {
		studentClass := ""
		if r.Class.Valid {
			studentClass = r.Class.String
		}
		// class filter: skip non-matching if filter active
		if classFilter != "" && studentClass != classFilter {
			continue
		}
		out = append(out, StudentResponse{
			ID:        uuidToString(r.ID),
			SchoolID:  uuidToString(r.SchoolID),
			UserID:    uuidToString(r.UserID),
			ClassName: studentClass,
			UserName:  r.UserName,
			UserEmail: r.UserEmail,
		})
	}
	return out, nil
}

func (s *StudentService) ListByParent(ctx context.Context, schoolID, parentID string) ([]StudentResponse, error) {
	schoolUUID := pgtype.UUID{}
	if err := schoolUUID.Scan(schoolID); err != nil {
		return nil, errors.New("invalid school_id")
	}
	parentUUID := pgtype.UUID{}
	if err := parentUUID.Scan(parentID); err != nil {
		return nil, errors.New("invalid parent_id")
	}

	rows, err := s.queries.GetStudentsByParent(ctx, db.GetStudentsByParentParams{
		SchoolID: schoolUUID,
		ParentID: parentUUID,
	})
	if err != nil {
		return nil, err
	}

	out := make([]StudentResponse, 0, len(rows))
	for _, r := range rows {
		studentClass := ""
		if r.Class.Valid {
			studentClass = r.Class.String
		}
		out = append(out, StudentResponse{
			ID:        uuidToString(r.ID),
			SchoolID:  uuidToString(r.SchoolID),
			UserID:    uuidToString(r.UserID),
			ClassName: studentClass,
		})
	}
	return out, nil
}

func (s *StudentService) GetByID(ctx context.Context, schoolID, studentID string) (*StudentResponse, error) {
	id := pgtype.UUID{}
	if err := id.Scan(studentID); err != nil {
		return nil, errors.New("invalid student id")
	}

	row, err := s.queries.GetStudentByID(ctx, id)
	if err != nil {
		return nil, errors.New("student not found")
	}

	if uuidToString(row.SchoolID) != schoolID {
		return nil, errors.New("student not found")
	}

	rowUser, err := s.queries.GetUserByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}

	className := ""
	if row.Class.Valid {
		className = row.Class.String
	}

	return &StudentResponse{
		ID:        uuidToString(row.ID),
		SchoolID:  uuidToString(row.SchoolID),
		UserID:    uuidToString(row.UserID),
		ClassName: className,
		UserName:  rowUser.Name,
		UserEmail: rowUser.Email,
	}, nil
}

func (s *StudentService) GetByUserID(ctx context.Context, schoolID, userID string) (*StudentResponse, error) {
	uid := pgtype.UUID{}
	if err := uid.Scan(userID); err != nil {
		return nil, errors.New("invalid user_id")
	}
	sid := pgtype.UUID{}
	if err := sid.Scan(schoolID); err != nil {
		return nil, errors.New("invalid school_id")
	}

	row, err := s.queries.GetStudentByUserID(ctx, db.GetStudentByUserIDParams{UserID: uid, SchoolID: sid})
	if err != nil {
		return nil, errors.New("student not found")
	}

	rowUser, err := s.queries.GetUserByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}

	className := ""
	if row.Class.Valid {
		className = row.Class.String
	}
	return &StudentResponse{
		ID:        uuidToString(row.ID),
		SchoolID:  uuidToString(row.SchoolID),
		UserID:    uuidToString(row.UserID),
		ClassName: className,
		UserName:  rowUser.Name,
		UserEmail: rowUser.Email,
	}, nil
}

func (s *StudentService) Create(ctx context.Context, schoolID string, input CreateStudentInput) (*StudentResponse, error) {
	schoolUUID := pgtype.UUID{}
	if err := schoolUUID.Scan(schoolID); err != nil {
		return nil, errors.New("invalid school_id")
	}

	userUUID := pgtype.UUID{}
	if err := userUUID.Scan(input.UserID); err != nil {
		return nil, errors.New("invalid user_id")
	}

	// FK check: user exists
	user, err := s.queries.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Cross-tenant check: user must belong to same school
	count, err := s.queries.GetUserBySchoolAndRole(ctx, db.GetUserBySchoolAndRoleParams{
		UserID:   userUUID,
		SchoolID: schoolUUID,
	})
	if err != nil || count == 0 {
		return nil, errors.New("user not found in your school")
	}

	classText := pgtype.Text{}
	if input.ClassName != "" {
		classText = pgtype.Text{String: input.ClassName, Valid: true}
	}

	student, err := s.queries.CreateStudent(ctx, db.CreateStudentParams{
		SchoolID: schoolUUID,
		UserID:   userUUID,
		Class:    classText,
	})
	if err != nil {
		return nil, err
	}

	className := ""
	if student.Class.Valid {
		className = student.Class.String
	}

	return &StudentResponse{
		ID:        uuidToString(student.ID),
		SchoolID:  uuidToString(student.SchoolID),
		UserID:    uuidToString(student.UserID),
		ClassName: className,
		UserName:  user.Name,
		UserEmail: user.Email,
	}, nil
}

// ponytail: UUID helpers below; move to shared pkg if >2 services need them.
func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return u.String()
}
