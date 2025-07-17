package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token valid for 1 day
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
