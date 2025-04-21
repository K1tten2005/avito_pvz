package usecase

import (
	"context"
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

func (uc *PvzUsecase) GetPvz(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	pvz, err := uc.repo.GetPvz(ctx, startDate, endDate, page, limit)
	if err != nil {
		loggerVar.Error(err.Error())
		return nil, err
	}
	return pvz, nil
}

func (uc *PvzUsecase) CreateReception(ctx context.Context, PvzId uuid.UUID) (models.Reception, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	active, err := uc.repo.HasActiveReception(ctx, PvzId)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Reception{}, err
	}
	if active {
		loggerVar.Error(pvz.ErrActiveReceptionExists.Error())
		return models.Reception{}, pvz.ErrActiveReceptionExists
	}

	reception := models.Reception{
		Id:       uuid.NewV4(),
		DateTime: time.Now(),
		PvzId:    PvzId,
		Status:   models.StatusInProgress,
	}

	err = uc.repo.CreateReception(ctx, reception)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Reception{}, err
	}

	loggerVar.Info("Success")
	return reception, nil
}

func (uc *PvzUsecase) AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (*models.Product, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	reception, err := uc.repo.GetActiveReception(ctx, pvzID)
	if err != nil {
		loggerVar.Error(err.Error())
		return nil, err
	}

	product := &models.Product{
		Id:          uuid.NewV4(),
		DateTime:    time.Now(),
		Type:        productType,
		ReceptionId: reception.Id,
	}

	if err := uc.repo.AddProduct(ctx, product); err != nil {
		return nil, err
	}

	loggerVar.Info("Success")
	return product, nil
}

func (uc *PvzUsecase) DeleteProduct(ctx context.Context, pvzID uuid.UUID) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	active, err := uc.repo.HasActiveReception(ctx, pvzID)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}
	if !active {
		loggerVar.Error(pvz.ErrNoActiveReception.Error())
		return pvz.ErrNoActiveReception
	}

	product, err := uc.repo.GetLastProduct(ctx, pvzID)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}

	if err := uc.repo.DeleteProduct(ctx, product.Id); err != nil {
		loggerVar.Error(err.Error())
		return err
	}

	loggerVar.Info("Success")
	return nil
}

func (uc *PvzUsecase) CloseReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	reception, err := uc.repo.GetActiveReception(ctx, pvzID)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Reception{}, err
	}

	_, err = uc.repo.GetLastProduct(ctx, pvzID)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Reception{}, err
	}

	reception.Status = models.StatusClose
	if err := uc.repo.UpdateReceptionStatus(ctx, reception.Id, reception.Status); err != nil {
		loggerVar.Error(err.Error())
		return models.Reception{}, err
	}

	loggerVar.Info("Success")
	return reception, nil
}
