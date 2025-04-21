package jwtUtils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

const secret = "test_secret"

func createTestJWT(t *testing.T, claims jwt.MapClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	return tokenStr
}

func TestGetRoleFromJWT(t *testing.T) {
	claims := jwt.MapClaims{
		"role": "employee",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	tokenStr := createTestJWT(t, claims, secret)

	role, ok := GetRoleFromJWT(tokenStr, jwt.MapClaims{}, secret)
	assert.True(t, ok)
	assert.Equal(t, "employee", role)
}

func TestGenerateJWTForTest(t *testing.T) {
	role := "employee"
	secret := "secret"

	tokenStr := GenerateJWTForTest(t, role, secret)
	assert.NotNil(t, tokenStr)
}
