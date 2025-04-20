package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/K1tten2005/avito_pvz/internal/pkg/pvz"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
	"github.com/satori/uuid"
)

type PvzUsecase struct {
	repo pvz.PvzRepo
}

func CreatePvzUsecase(repo pvz.PvzRepo) *PvzUsecase {
	return &PvzUsecase{repo: repo}
}

func (uc *PvzUsecase) CreatePvz(ctx context.Context, pvz models.PVZ) error {
	return uc.repo.InsertPvz(ctx, pvz)
}

func (u *PvzUsecase) GetPvz(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	return u.repo.GetPvz(ctx, startDate, endDate, page, limit)
}

func (uc *PvzUsecase) CloseReception(ctx context.Context, receptionID uuid.UUID) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	reception, err := uc.repo.GetReceptionByID(ctx, receptionID)
	if err != nil {
		loggerVar.Error(err.Error())
		return errors.New("приемка не найдена")
	}

	if reception.Status == "close" {
		return errors.New("приемка уже закрыта")
	}

	return uc.repo.UpdateReceptionStatus(ctx, receptionID, "close")
}
