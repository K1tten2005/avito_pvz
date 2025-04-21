package pvz

import (
	"context"
	"errors"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/satori/uuid"
)

var (
	ErrActiveReceptionExists = errors.New("active reception already exsists")
	ErrNoActiveReception     = errors.New("no active reception")
	ErrNoProductsInReception = errors.New("no products in reception")
)

type PvzRepo interface {
	InsertPvz(ctx context.Context, pvz models.PVZ) error
	InsertReception(ctx context.Context, reception models.Reception) error
	InsertProduct(ctx context.Context, product models.Product) error
	GetPvz(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error)
	GetActiveReception(ctx context.Context, pvzId uuid.UUID) (models.Reception, error)
	GetReceptionByID(ctx context.Context, id uuid.UUID) (models.Reception, error)
	UpdateReceptionStatus(ctx context.Context, id uuid.UUID, status string) error
	HasActiveReception(ctx context.Context, pvzID uuid.UUID) (bool, error)
	CreateReception(ctx context.Context, reception models.Reception) error
	AddProduct(ctx context.Context, product *models.Product) error 
	GetLastProduct(ctx context.Context, pvzID uuid.UUID) (models.Product, error)
	DeleteProduct(ctx context.Context, productId uuid.UUID) error
}

type PvzUsecase interface {
	CreatePvz(ctx context.Context, pvz models.PVZ) error
	GetPvz(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error)
	CloseReception(ctx context.Context, receptionID uuid.UUID) (models.Reception, error)
	CreateReception(ctx context.Context, PvzId uuid.UUID) (models.Reception, error)
	AddProduct(ctx context.Context, pvzID uuid.UUID, productType string) (*models.Product, error)
	DeleteProduct(ctx context.Context, pvzID uuid.UUID) error
}
