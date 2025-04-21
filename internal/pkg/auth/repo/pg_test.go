package repo

import (
	"context"
	"testing"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthRepo_InsertUser(t *testing.T) {
	// Тестовые данные
	user := models.User{
		Id:           uuid.NewV4(),
		Email:        "test@example.com",
		Role:         "employee",
		PasswordHash: []byte("hashed_password"),
	}

	tests := []struct {
		name       string
		repoMocker func(*pgxpoolmock.MockPgxPool)
		err        error
	}{
		{
			name: "Success",
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(
					gomock.Any(), // context
					gomock.Any(), // SQL-запрос (insertUser)
					user.Id,
					user.Email,
					user.Role,
					user.PasswordHash,
				).Return(nil, nil)
			},
			err: nil,
		},
		{
			name: "Duplicate email",
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().Exec(
					gomock.Any(),
					gomock.Any(),
					user.Id,
					user.Email,
					user.Role,
					user.PasswordHash,
				).Return(nil, assert.AnError) // Или конкретная ошибка из pgx
			},
			err: assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			test.repoMocker(mockPool)

			authRepo := CreateAuthRepo(mockPool)
			err := authRepo.InsertUser(context.Background(), user)

			assert.Equal(t, test.err, err)
		})
	}
}

func TestAuthRepo_SelectUserByEmail(t *testing.T) {
	testEmail := "test@example.com"
	testUser := models.User{
		Id:           uuid.NewV4(),
		Email:        testEmail,
		Role:         "employee",
		PasswordHash: []byte("hashed_password"),
	}

	tests := []struct {
		name        string
		email       string
		repoMocker  func(*pgxpoolmock.MockPgxPool)
		expected    models.User
		expectedErr error
	}{
		{
			name:  "Success",
			email: testEmail,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"id", "role", "password_hash"}).
					AddRow(testUser.Id, testUser.Role, testUser.PasswordHash).ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), gomock.Any(), testEmail).
					Return(rows)
			},
			expected:    testUser,
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			test.repoMocker(mockPool)

			authRepo := CreateAuthRepo(mockPool)
			result, err := authRepo.SelectUserByEmail(context.Background(), test.email)

			assert.Equal(t, test.expectedErr, err)
			if test.expectedErr == nil {
				assert.Equal(t, test.expected.Id, result.Id)
				assert.Equal(t, test.expected.Email, result.Email)
				assert.Equal(t, test.expected.Role, result.Role)
				assert.Equal(t, test.expected.PasswordHash, result.PasswordHash)
			}
		})
	}
}
