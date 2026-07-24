package middleware

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func BenchmarkTokenParse(b *testing.B) {
	// Setup JWT secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test_secret_key_for_benchmark_only_12345"
		os.Setenv("JWT_SECRET", secret)
	}

	// Create sample token
	claims := &CustomClaims{
		SchoolID: "school-123",
		Roles:    []string{"teacher"},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-456",
			ID:        "jti-789",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))

	// Benchmark parsing
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	}
}
