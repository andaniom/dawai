package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/violin-assessment/dawai/internal/db"
	"github.com/violin-assessment/dawai/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries *db.Queries
}

func NewAuthService(queries *db.Queries) *AuthService {
	return &AuthService{queries: queries}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func (s *AuthService) VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// LoginResult contains the JWT token and user info returned after successful login.
type LoginResult struct {
	Token string
	User  LoginUser
}

// LoginUser is the user shape returned by auth operations.
type LoginUser struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	SchoolID string   `json:"school_id"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.PasswordHash.Valid {
		return nil, errors.New("invalid email or password")
	}

	if err := s.VerifyPassword(user.PasswordHash.String, password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// ponytail: hardcoded roles + school; wire user_roles + user_schools queries when seed data ready
	rolesStr := []string{"teacher"}
	schoolID := "00000000-0000-0000-0000-000000000000"

	// Fetch first school from user_schools for login
	schools, err := s.queries.GetUserSchools(ctx, user.ID)
	if err == nil && len(schools) > 0 {
		schoolID = schools[0].String()
	}

	token, err := s.GenerateToken(ctx, user.ID.String(), schoolID, rolesStr)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User: LoginUser{
			ID:       user.ID.String(),
			Email:    user.Email,
			Name:     user.Name,
			Roles:    rolesStr,
			SchoolID: schoolID,
		},
	}, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, userID, schoolID string, roles []string) (string, error) {
	jti := uuid.New().String()
	now := time.Now()

	claims := &middleware.CustomClaims{
		SchoolID: schoolID,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) Logout(ctx context.Context, jti string) error {
	_, err := s.queries.BlacklistJWT(ctx, db.BlacklistJWTParams{
		Jti:       jti,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	return err
}

func (s *AuthService) CreateUser(ctx context.Context, email, password, name string) (*LoginUser, error) {
	hash, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: pgtype.Text{String: hash, Valid: true},
		Name:         name,
	})
	if err != nil {
		return nil, err
	}

	return &LoginUser{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
		Roles: []string{},
	}, nil
}

type UserWithRoles struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
	Schools   []SchoolInfo `json:"schools"`
}

type SchoolInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *AuthService) ListUsers(ctx context.Context, schoolID, roleFilter string) ([]UserWithRoles, error) {
	schoolUUID := pgtype.UUID{}
	if err := schoolUUID.Scan(schoolID); err != nil {
		return nil, errors.New("invalid school_id")
	}

	rows, err := s.queries.ListUsersBySchool(ctx, schoolUUID)
	if err != nil {
		return nil, err
	}

	// Aggregate users by ID (each user can have multiple roles)
	userMap := make(map[string]*UserWithRoles)
	for _, r := range rows {
		id := uuidValToString(r.ID)
		if _, ok := userMap[id]; !ok {
			userMap[id] = &UserWithRoles{
				ID:    id,
				Email: r.Email,
				Name:  r.Name,
				Roles: []string{},
			}
		}
		if r.RoleName != "" {
			userMap[id].Roles = append(userMap[id].Roles, r.RoleName)
		}
	}

	out := make([]UserWithRoles, 0, len(userMap))
	for _, u := range userMap {
		if roleFilter != "" {
			hasRole := false
			for _, r := range u.Roles {
				if r == roleFilter {
					hasRole = true
					break
				}
			}
			if !hasRole {
				continue
			}
		}
		out = append(out, *u)
	}

	return out, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID string) (*UserWithRoles, error) {
	userUUID := pgtype.UUID{}
	if err := userUUID.Scan(userID); err != nil {
		return nil, errors.New("invalid user_id")
	}

	user, err := s.queries.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get schools
	schools, err := s.queries.GetUserSchoolsWithName(ctx, userUUID)
	if err != nil {
		schools = nil
	}

	schoolInfos := make([]SchoolInfo, 0, len(schools))
	seenSchools := make(map[string]bool)
	for _, sch := range schools {
		if seenSchools[sch.ID] {
			continue
		}
		seenSchools[sch.ID] = true
		schoolInfos = append(schoolInfos, SchoolInfo{ID: sch.ID, Name: sch.Name})
	}

	// Get roles — ponytail: iterate all school memberships for roles, single school for now
	roles := []string{}

	// Get first school's roles
	schoolsRaw, _ := s.queries.GetUserSchools(ctx, userUUID)
	if len(schoolsRaw) > 0 {
		firstSchool := schoolsRaw[0]
		roleNames, _ := s.queries.GetRoleNamesByUserSchool(ctx, db.GetRoleNamesByUserSchoolParams{
			UserID:   userUUID,
			SchoolID: firstSchool,
		})
		for _, rn := range roleNames {
			roles = append(roles, rn)
		}
	}

	return &UserWithRoles{
		ID:      uuidValToString(user.ID),
		Email:   user.Email,
		Name:    user.Name,
		Roles:   roles,
		Schools: schoolInfos,
	}, nil
}

// uuidValToString converts a pgtype.UUID to string. ponytail: duplicate across services, move to shared pkg if >2.
func uuidValToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return u.String()
}
