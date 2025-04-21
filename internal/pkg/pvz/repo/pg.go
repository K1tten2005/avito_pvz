package repo

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/K1tten2005/avito_pvz/internal/pkg/pvz"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/satori/uuid"
)

//go:embed sql/insertPvz.sql
var insertPvz string

//go:embed sql/insertReception.sql
var insertReception string

//go:embed sql/insertProduct.sql
var insertProduct string

//go:embed sql/deleteProduct.sql
var deleteProduct string

//go:embed sql/getPvz.sql
var getPvz string

//go:embed sql/updateReceptionStatus.sql
var updateReceptionStatus string

//go:embed sql/getReceptionById.sql
var getReceptionById string

//go:embed sql/getActiveReception.sql
var getActiveReception string

//go:embed sql/hasActiveReception.sql
var hasActiveReception string

//go:embed sql/createReception.sql
var createReception string

//go:embed sql/addProduct.sql
var addProduct string

//go:embed sql/getLastProduct.sql
var getLastProduct string

type PvzRepo struct {
	db pgxtype.Querier
}

func CreatePvzRepo(db pgxtype.Querier) *PvzRepo {
	return &PvzRepo{
		db: db,
	}
}

func uuidToString(u pgtype.UUID) string {
	if u.Status != pgtype.Present {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		u.Bytes[0:4],
		u.Bytes[4:6],
		u.Bytes[6:8],
		u.Bytes[8:10],
		u.Bytes[10:16],
	)
}

func (repo *PvzRepo) InsertPvz(ctx context.Context, pvz models.PVZ) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertPvz, pvz.Id, pvz.RegistrationDate, pvz.City)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}
	loggerVar.Info("Successful")
	return nil
}

func (repo *PvzRepo) GetPvz(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZ, error) {
	offset := (page - 1) * limit

	rows, err := repo.db.Query(ctx, getPvz, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pvzMap := make(map[string]*models.PVZ)

	for rows.Next() {
		var (
			pvzID, receptionID, productID       pgtype.UUID
			pvzCity, productCategory            string
			pvzDate, receptionDate, productDate time.Time
			receptionStatus                     string
		)

		err := rows.Scan(&pvzID, &pvzCity, &pvzDate,
			&receptionID, &receptionDate, &receptionStatus,
			&productID, &productDate, &productCategory)
		if err != nil {
			return nil, err
		}

		if pvzID.Status != pgtype.Present {
			continue
		}

		pvzIDStr := uuidToString(pvzID)

		pvzUUID, err := uuid.FromString(pvzIDStr)
		if err != nil {
			return nil, err
		}

		if _, exists := pvzMap[pvzIDStr]; !exists {
			pvzMap[pvzIDStr] = &models.PVZ{
				Id:               pvzUUID,
				City:             pvzCity,
				RegistrationDate: pvzDate,
				Receptions:       []models.Reception{},
			}
		}

		if receptionID.Status == pgtype.Present {
			receptionIDStr := uuidToString(receptionID)
			receptionUUID, err := uuid.FromString(receptionIDStr)
			if err != nil {
				return nil, err
			}
			recFound := false

			for i := range pvzMap[pvzIDStr].Receptions {
				if pvzMap[pvzIDStr].Receptions[i].Id == receptionUUID {
					if productID.Status == pgtype.Present {
						productIDStr := uuidToString(productID)
						productUUID, err := uuid.FromString(productIDStr)
						if err != nil {
							return nil, err
						}
						pvzMap[pvzIDStr].Receptions[i].Products = append(
							pvzMap[pvzIDStr].Receptions[i].Products,
							models.Product{
								Id:          productUUID,
								DateTime:    productDate,
								Type:        productCategory,
								ReceptionId: receptionUUID,
							})
					}
					recFound = true
					break
				}
			}

			if !recFound {
				rec := models.Reception{
					Id:       receptionUUID,
					DateTime: receptionDate,
					PvzId:    pvzUUID,
					Status:   receptionStatus,
				}

				if productID.Status == pgtype.Present {
					productIDStr := uuidToString(productID)
					productUUID, err := uuid.FromString(productIDStr)
					if err != nil {
						return nil, err
					}
					rec.Products = append(rec.Products, models.Product{
						Id:          productUUID,
						DateTime:    productDate,
						Type:        productCategory,
						ReceptionId: receptionUUID,
					})
				}

				pvzMap[pvzIDStr].Receptions = append(pvzMap[pvzIDStr].Receptions, rec)
			}
		}
	}

	var result []models.PVZ
	for _, val := range pvzMap {
		result = append(result, *val)
	}
	return result, nil
}

func (repo *PvzRepo) GetReceptionByID(ctx context.Context, id uuid.UUID) (models.Reception, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	var reception models.Reception
	err := repo.db.QueryRow(ctx, getReceptionById, id).
		Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Reception{}, err
	}
	return reception, nil
}

func (repo *PvzRepo) InsertReception(ctx context.Context, reception models.Reception) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertReception, reception.Id, reception.DateTime, reception.PvzId, reception.Status)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}
	loggerVar.Info("Successful")
	return nil
}

func (repo *PvzRepo) InsertProduct(ctx context.Context, product models.Product) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertProduct, product.Id, product.DateTime, product.ReceptionId, product.Type)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}
	loggerVar.Info("Successful")
	return nil
}

func (repo *PvzRepo) HasActiveReception(ctx context.Context, pvzID uuid.UUID) (bool, error) {
	var exists bool
	err := repo.db.QueryRow(ctx, hasActiveReception, pvzID).Scan(&exists)
	return exists, err
}

func (repo *PvzRepo) GetActiveReception(ctx context.Context, pvzId uuid.UUID) (models.Reception, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	var reception models.Reception
	err := repo.db.QueryRow(ctx, getActiveReception, pvzId).
		Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
	if err == sql.ErrNoRows {
		return models.Reception{}, pvz.ErrNoActiveReception
	}
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Reception{}, err
	}
	return reception, nil
}

func (repo *PvzRepo) CreateReception(ctx context.Context, reception models.Reception) error {
	_, err := repo.db.Exec(ctx, createReception, reception.Id, reception.DateTime, reception.PvzId, reception.Status)
	return err
}

func (repo *PvzRepo) AddProduct(ctx context.Context, product *models.Product) error {
	_, err := repo.db.Exec(ctx, addProduct, product.Id, product.DateTime, product.Type, product.ReceptionId)
	return err
}

func (repo *PvzRepo) GetLastProduct(ctx context.Context, pvzID uuid.UUID) (models.Product, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	var product models.Product
	err := repo.db.QueryRow(ctx, getLastProduct, pvzID).Scan(
		&product.Id,
		&product.DateTime,
		&product.Type,
		&product.ReceptionId,
	)

	if err == sql.ErrNoRows {
		loggerVar.Error(pvz.ErrNoProductsInReception.Error())
		return models.Product{}, pvz.ErrNoProductsInReception
	}
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Product{}, err
	}

	return product, nil
}


func (repo *PvzRepo) DeleteProduct(ctx context.Context, productId uuid.UUID) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	result, err := repo.db.Exec(ctx, deleteProduct, productId)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		loggerVar.Error("product not found")
		return errors.New("product not found")
	}

	loggerVar.Info("Successful")
	return nil
}

func (repo *PvzRepo) UpdateReceptionStatus(ctx context.Context, id uuid.UUID, status string) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	cmd, err := repo.db.Exec(ctx, updateReceptionStatus, status, id)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("ничего не обновлено")
	}
	return nil
}