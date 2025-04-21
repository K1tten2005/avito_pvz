package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/K1tten2005/avito_pvz/internal/pkg/auth"
	"github.com/K1tten2005/avito_pvz/internal/pkg/auth/mocks"
	"github.com/golang/mock/gomock"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockAuthUsecase(ctrl)

	tests := []struct {
		name             string
		reqBody          string
		mockBehavior     func()
		expectedStatus   int
		expectedResponse string
		expectedCookie   string
	}{
		{
			name:             "Error parsing JSON",
			reqBody:          `{"role": "employee"`, // Ошибка в JSON
			mockBehavior:     func() {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "ошибка парсинга JSON",
			expectedCookie:   "",
		},
		{
			name:    "Error generating token",
			reqBody: `{"role": "employee"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().DummyLogin(gomock.Any(), gomock.Any()).Return("", auth.ErrGeneratingToken)
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "ошибка генерации токена",
			expectedCookie:   "",
		},
		{
			name:    "Successful login",
			reqBody: `{"role": "employee"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().DummyLogin(gomock.Any(), gomock.Any()).Return("valid_token", nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "valid_token",
			expectedCookie:   "valid_token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос
			req := httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Устанавливаем мок-бихевиор для каждого теста
			tt.mockBehavior()

			// Создаем новый респондер
			rr := httptest.NewRecorder()
			handler := &AuthHandler{uc: mockUsecase}

			// Вызываем обработчик
			handler.DummyLogin(rr, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Проверяем тело ответа
			body := rr.Body.String()
			assert.Contains(t, body, tt.expectedResponse)

			// Проверяем куки
			cookies := rr.Result().Cookies()
			if tt.expectedCookie != "" {
				assert.Len(t, cookies, 1)
				assert.Equal(t, "AvitoJWT", cookies[0].Name)
				assert.Equal(t, tt.expectedCookie, cookies[0].Value)
				assert.True(t, cookies[0].HttpOnly)
			} else {
				assert.Len(t, cookies, 0)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockAuthUsecase(ctrl)

	tests := []struct {
		name             string
		reqBody          string
		mockBehavior     func()
		expectedStatus   int
		expectedResponse string
		expectedCookie   string
	}{
		{
			name:             "Invalid JSON format",
			reqBody:          `{"email": "test@example.com"`, 
			mockBehavior:     func() {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "ошибка парсинга JSON",
			expectedCookie:   "",
		},
		{
			name:    "Invalid credentials",
			reqBody: `{"email": "wrong@example.com", "password": "wrongpassword"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(models.User{}, "", auth.ErrInvalidCredentials)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "неверный логин или пароль",
			expectedCookie:   "",
		},
		{
			name:    "User not found",
			reqBody: `{"email": "nonexistent@example.com", "password": "password123"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(models.User{}, "", auth.ErrUserNotFound)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "пользователь не найден",
			expectedCookie:   "",
		},
		{
			name:    "Successful login",
			reqBody: `{"email": "test@example.com", "password": "password123"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(models.User{}, "valid_token", nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "valid_token",
			expectedCookie:   "valid_token",
		},
		{
			name:    "Unknown error",
			reqBody: `{"email": "test@example.com", "password": "password123"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(models.User{}, "", fmt.Errorf("неизвестная ошибка"))
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "неизвестная ошибка",
			expectedCookie:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			tt.mockBehavior()

			rr := httptest.NewRecorder()
			handler := &AuthHandler{uc: mockUsecase}

			handler.Login(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			body := rr.Body.String()
			assert.Contains(t, body, tt.expectedResponse)

			cookies := rr.Result().Cookies()
			if tt.expectedCookie != "" {
				assert.Len(t, cookies, 1)
				assert.Equal(t, "AvitoJWT", cookies[0].Name)
				assert.Equal(t, tt.expectedCookie, cookies[0].Value)
				assert.True(t, cookies[0].HttpOnly)
			} else {
				assert.Len(t, cookies, 0)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := uuid.NewV4()

	mockUsecase := mocks.NewMockAuthUsecase(ctrl)

	tests := []struct {
		name             string
		reqBody          string
		mockBehavior     func()
		expectedStatus   int
		expectedResponse string
		expectedCookie   string
	}{
		{
			name:             "Invalid JSON format",
			reqBody:          `{"email": "test@example.com"`, 
			mockBehavior:     func() {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "ошибка парсинга JSON",
			expectedCookie:   "",
		},
		{
			name:             "Invalid email or password",
			reqBody:          `{"email": "invalidemail", "password": "password123", "role": "employee"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().Register(gomock.Any(), gomock.Any()).Return(models.User{}, "", auth.ErrInvalidEmail)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "неправильный логин или пароль",
			expectedCookie:   "",
		},
		{
			name:             "Successful registration",
			reqBody:          `{"email": "test@example.com", "password": "password123", "role": "employee"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().Register(gomock.Any(), gomock.Any()).Return(models.User{Id: id, Email: "test@example.com", Role: "employee"}, "valid_token", nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: fmt.Sprintf(`{"id":"%s","email":"test@example.com","role":"employee"}`, id),
			expectedCookie:   "valid_token",
		},
		{
			name:             "Unknown error",
			reqBody:          `{"email": "test@example.com", "password": "password123", "role": "employee"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().Register(gomock.Any(), gomock.Any()).Return(models.User{}, "", fmt.Errorf("неизвестная ошибка"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "неизвестная ошибка",
			expectedCookie:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			tt.mockBehavior()

			// Создаем новый респондер
			rr := httptest.NewRecorder()
			handler := &AuthHandler{uc: mockUsecase}

			// Вызываем обработчик
			handler.Register(rr, req)

			// Проверяем статус код
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Проверяем тело ответа
			body := rr.Body.String()
			assert.Contains(t, body, tt.expectedResponse)

			// Проверяем куки
			cookies := rr.Result().Cookies()
			if tt.expectedCookie != "" {
				assert.Len(t, cookies, 1)
				assert.Equal(t, "AvitoJWT", cookies[0].Name)
				assert.Equal(t, tt.expectedCookie, cookies[0].Value)
				assert.True(t, cookies[0].HttpOnly)
			} else {
				assert.Len(t, cookies, 0)
			}
		})
	}
}
