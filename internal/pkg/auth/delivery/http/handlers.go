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
	"github.com/satori/uuid"
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
		logger.LogHandlerError(loggerVar, errors.New("invalid role format"), http.StatusBadRequest)
		send_err.SendError(w, "invalid role format", http.StatusBadRequest)
	}

	token, err := h.uc.DummyLogin(r.Context(), req)

	if err != nil {
		switch err {
		case auth.ErrGeneratingToken:
			logger.LogHandlerError(loggerVar, err, http.StatusUnauthorized)
			send_err.SendError(w, err.Error(), http.StatusUnauthorized)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("unknkown error: %w", err), http.StatusBadRequest)
			send_err.SendError(w, "unknkown error", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(token))
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusOK)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.LoginReq
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error while parsing JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "error while parsing JSON", http.StatusBadRequest)
		return
	}
	req.Sanitize()

	_, token, err := h.uc.Login(r.Context(), req)

	if err != nil {
		switch err {
		case auth.ErrInvalidEmail, auth.ErrUserNotFound, auth.ErrInvalidCredentials:
			logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
			send_err.SendError(w, err.Error(), http.StatusBadRequest)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("unknkown error: %w", err), http.StatusInternalServerError)
			send_err.SendError(w, "unknkown error", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(token))
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusOK)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.RegisterReq
	err := easyjson.UnmarshalFromReader(r.Body, &req)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error while parsing JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "error while parsing JSON", http.StatusBadRequest)
		return
	}
	req.Sanitize()

	if !validation.IsValidRole(req.Role) {
		logger.LogHandlerError(loggerVar, errors.New("invalid role format"), http.StatusBadRequest)
		send_err.SendError(w, "invalid role format", http.StatusBadRequest)
	}

	user, token, err := h.uc.Register(r.Context(), req)

	if err != nil {
		switch err {
		case auth.ErrInvalidEmail, auth.ErrInvalidPassword:
			logger.LogHandlerError(loggerVar, fmt.Errorf("invalid mail or password: %w", err), http.StatusBadRequest)
			send_err.SendError(w, "invalid mail or password", http.StatusBadRequest)
		default:
			logger.LogHandlerError(loggerVar, fmt.Errorf("unknkown error: %w", err), http.StatusInternalServerError)
			send_err.SendError(w, "unknkown error", http.StatusBadRequest)
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

	resp := struct {
		Id    uuid.UUID `json:"id"`
		Email string    `json:"email"`
		Role  string    `json:"role"`
	}{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error while parsing JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "error while parsing JSON", http.StatusBadRequest)
	}
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
}
