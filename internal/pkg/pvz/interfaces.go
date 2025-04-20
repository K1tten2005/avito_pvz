package pvz

import (
	"context"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/satori/uuid"
)

type PvzRepo interface {
	InsertPvz(ctx context.Context, pvz models.PVZ) error
	InsertReception(ctx context.Context, reception models.Reception) error
	InsertProduct(ctx context.Context, product models.Product) error
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
	GetPvz(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error)
	GetReceptionByID(ctx context.Context, id uuid.UUID) (*models.Reception, error)
	UpdateReceptionStatus(ctx context.Context, id uuid.UUID, status string) error
}

type PvzUsecase interface {
	CreatePvz(ctx context.Context, pvz models.PVZ) error
	GetPvz(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error)
	CloseReception(ctx context.Context, receptionID uuid.UUID) error
}
