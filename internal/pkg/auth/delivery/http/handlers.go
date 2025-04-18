package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/K1tten2005/avito_pvz/internal/pkg/auth"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/send_err"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/validation"
	"github.com/mailru/easyjson"
)

type AuthHandler struct {
	uc     auth.AuthUsecase
	secret string
}

func CreateAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc, secret: os.Getenv("JWT_SECRET")}
}

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.DummyLoginReq

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка парсинга JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	req.Sanitize()

	if !validation.IsValidRole(req.Role) {
		logger.LogHandlerError(loggerVar, errors.New("невалидная роль"), http.StatusBadRequest)
		send_err.SendError(w, "невалидная роль", http.StatusBadRequest)
	}

	token, csrfToken, err := h.uc.DummyLogin(r.Context(), req)

	if err != nil {
		switch err {
		case auth.ErrGeneratingToken:
			logger.LogHandlerError(loggerVar, err, http.StatusUnauthorized)
			send_err.SendError(w, err.Error(), http.StatusUnauthorized)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("неизвестная ошибка: %w", err), http.StatusBadRequest)
			send_err.SendError(w, "неизвестная ошибка", http.StatusBadRequest)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "AvitoJWT",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "CSRF-Token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.Header().Set("X-CSRF-Token", csrfToken)
	w.Header().Set("Content-Type", "application/json")

	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusOK)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.LoginReq
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка парсинга JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	req.Sanitize()

	user, token, csrfToken, err := h.uc.Login(r.Context(), req)

	if err != nil {
		switch err {
		case auth.ErrInvalidEmail, auth.ErrUserNotFound:
			logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
			send_err.SendError(w, err.Error(), http.StatusBadRequest)
		case auth.ErrInvalidCredentials:
			logger.LogHandlerError(loggerVar, err, http.StatusUnauthorized)
			send_err.SendError(w, err.Error(), http.StatusUnauthorized)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("неизвестная ошибка: %w", err), http.StatusInternalServerError)
			send_err.SendError(w, "неизвестная ошибка", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "AvitoJWT",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "CSRF-Token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.Header().Set("X-CSRF-Token", csrfToken)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка формирования JSON: %w", err), http.StatusInternalServerError)
		send_err.SendError(w, "ошибка формирования JSON", http.StatusInternalServerError)
	}
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusOK)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.RegisterReq
	err := easyjson.UnmarshalFromReader(r.Body, &req)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка парсинга JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	req.Sanitize()

	if !validation.IsValidRole(req.Role) {
		logger.LogHandlerError(loggerVar, errors.New("невалидная роль"), http.StatusBadRequest)
		send_err.SendError(w, "невалидная роль", http.StatusBadRequest)
	}

	user, token, csrfToken, err := h.uc.Register(r.Context(), req)

	if err != nil {
		switch err {
		case auth.ErrInvalidEmail, auth.ErrInvalidPassword:
			logger.LogHandlerError(loggerVar, fmt.Errorf("неправильный логин или пароль: %w", err), http.StatusBadRequest)
			send_err.SendError(w, "неправильный логин или пароль", http.StatusBadRequest)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("неизвестная ошибка: %w", err), http.StatusInternalServerError)
			send_err.SendError(w, "неизвестная ошибка", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "AvitoJWT",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "CSRF-Token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.Header().Set("X-CSRF-Token", csrfToken)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка формирования JSON: %w", err), http.StatusInternalServerError)
		send_err.SendError(w, "ошибка формирования JSON", http.StatusInternalServerError)
	}
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
}
