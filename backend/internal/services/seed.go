package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func CreateTestSchool(ctx context.Context, name string) (string, error) {
	return uuid.New().String(), nil
}

func CreateTestUser(ctx context.Context, schoolID, email, password string, roles []string) (string, error) {
	return uuid.New().String(), nil
}

func hashForTest(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

func parsePgUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	u.Scan(s)
	return u
}
