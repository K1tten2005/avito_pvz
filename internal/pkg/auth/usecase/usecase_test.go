package usecase

import (
	"context"
	"os"
	"testing"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuthUsecase_DummyLogin(t *testing.T) {
	oldSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Setenv("JWT_SECRET", oldSecret)

	uc := CreateAuthUsecase(nil)

	tests := []struct {
		name        string
		input       models.DummyLoginReq
		validate    func(t *testing.T, token string)
		expectedErr error
	}{
		{
			name: "success employee role",
			input: models.DummyLoginReq{
				Role: "employee",
			},
			validate: func(t *testing.T, token string) {
				claims := jwt.MapClaims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				assert.NoError(t, err)
				assert.Equal(t, "employee", claims["role"])
			},
			expectedErr: nil,
		},
		{
			name: "success moderator role",
			input: models.DummyLoginReq{
				Role: "moderator",
			},
			validate: func(t *testing.T, token string) {
				claims := jwt.MapClaims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test_secret"), nil
				})
				assert.NoError(t, err)
				assert.Equal(t, "moderator", claims["role"])
			},
			expectedErr: nil,
		},
		{
			name: "invalid role",
			input: models.DummyLoginReq{
				Role: "invalid_role",
			},
			validate: func(t *testing.T, token string) {
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := uc.DummyLogin(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				tt.validate(t, token)
			}
		})
	}
}
