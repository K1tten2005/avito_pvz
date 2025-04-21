package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/K1tten2005/avito_pvz/internal/pkg/pvz"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/send_err"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/validation"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/satori/uuid"
)

type PvzHandler struct {
	uc     pvz.PvzUsecase
	secret string
}

func CreatePvzHandler(uc pvz.PvzUsecase) *PvzHandler {
	return &PvzHandler{uc: uc, secret: os.Getenv("JWT_SECRET")}
}

func (h *PvzHandler) CreatePvz(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.PVZ

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка парсинга JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "ошибка парсинга JSON", http.StatusBadRequest)
		return
	}
	req.Sanitize()

	if req.Id == uuid.Nil {
		logger.LogHandlerError(loggerVar, errors.New("не передан UUID"), http.StatusBadRequest)
		send_err.SendError(w, "не передан UUID", http.StatusBadRequest)
		return
	}

	if _, err := uuid.FromString(req.Id.String()); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("невалидный UUID: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "невалидный UUID", http.StatusBadRequest)
		return
	}

	if !validation.IsValidCity(req.City) {
		logger.LogHandlerError(loggerVar, errors.New("не валидый город"), http.StatusBadRequest)
		send_err.SendError(w, "не валидый город", http.StatusBadRequest)
	}

	err := h.uc.CreatePvz(r.Context(), req)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка на уровне ниже (usecase): %w", err), http.StatusInternalServerError)
		send_err.SendError(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка формирования JSON: %w", err), http.StatusInternalServerError)
		send_err.SendError(w, "ошибка формирования JSON", http.StatusInternalServerError)
	}
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusOK)
}

func (h *PvzHandler) GetPvz(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil {
		limit = l
	}

	var startDate, endDate *time.Time
	var err error

	if startDateStr != "" {
		t, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			send_err.SendError(w, "wrong startDate format", http.StatusBadRequest)
			return
		}
		startDate = &t
	}
	if endDateStr != "" {
		t, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			send_err.SendError(w, "wrong endDate format", http.StatusBadRequest)
			return
		}
		endDate = &t
	}

	pvz, err := h.uc.GetPvz(r.Context(), startDate, endDate, page, limit)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка на уровне ниже (usecase): %w", err), http.StatusInternalServerError)
		send_err.SendError(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pvz); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка формирования JSON: %w", err), http.StatusInternalServerError)
		send_err.SendError(w, "ошибка формирования JSON", http.StatusInternalServerError)
	}
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusOK)
}

func (h *PvzHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.Reception

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка парсинга JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	if req.PvzId == uuid.Nil {
		logger.LogHandlerError(loggerVar, errors.New("PvzId not provided"), http.StatusBadRequest)
		send_err.SendError(w, "PvzId not provided", http.StatusBadRequest)
		return
	}

	reception, err := h.uc.CreateReception(r.Context(), req.PvzId)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error on a level below (usecase): %w", err), http.StatusBadRequest)
		send_err.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(reception); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка формирования JSON: %w", err), http.StatusInternalServerError)
		send_err.SendError(w, "ошибка формирования JSON", http.StatusInternalServerError)
		return
	}
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusCreated)
}

func (h *PvzHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	var req models.AddProductReq

	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка парсинга JSON: %w", err), http.StatusBadRequest)
		send_err.SendError(w, "ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	if !validation.IsValidProductType(req.Type) {
		logger.LogHandlerError(loggerVar, errors.New("wrong product type"), http.StatusBadRequest)
		send_err.SendError(w, "wrong product type", http.StatusBadRequest)
		return
	}

	product, err := h.uc.AddProduct(r.Context(), req.PvzId, req.Type)
	if err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("error on a level below (usecase): %w", err), http.StatusBadRequest)
		send_err.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка формирования JSON: %w", err), http.StatusInternalServerError)
		send_err.SendError(w, "ошибка формирования JSON", http.StatusInternalServerError)
		return
	}
	logger.LogHandlerInfo(loggerVar, "Successful", http.StatusCreated)
}

func (h *PvzHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	vars := mux.Vars(r)
	pvzIdStr := vars["pvzId"]
	pvzID, err := uuid.FromString(pvzIdStr)
	if err != nil {
		logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
		send_err.SendError(w, "wrong UUID format for pvzId query parameter", http.StatusBadRequest)
		return
	}

	err = h.uc.DeleteProduct(r.Context(), pvzID)
	if err != nil {
		logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
		send_err.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.LogHandlerInfo(loggerVar, "Success", http.StatusOK)
}

func (h *PvzHandler) CloseReception(w http.ResponseWriter, r *http.Request) {
	loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

	vars := mux.Vars(r)
	pvzIdStr := vars["pvzId"]
	pvzID, err := uuid.FromString(pvzIdStr)
	if err != nil {
		logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
		send_err.SendError(w, "wrong UUID format for pvzId query parameter", http.StatusBadRequest)
		return
	}

	reception, err := h.uc.CloseReception(r.Context(), pvzID)
	if err != nil {
		logger.LogHandlerError(loggerVar, err, http.StatusBadRequest)
		send_err.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(reception); err != nil {
		logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка кодирования JSON: %w", err), http.StatusInternalServerError)
		send_err.SendError(w, "ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	logger.LogHandlerInfo(loggerVar, "Success", http.StatusOK)
}
