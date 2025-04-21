package repo

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"log/slog"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/K1tten2005/avito_pvz/internal/pkg/pvz"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
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
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	offset := (page - 1) * limit

	rows, err := repo.db.Query(ctx, getPvz, startDate, endDate, limit, offset)
	if err != nil {
		loggerVar.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	pvzMap := make(map[string]*models.PVZ)

	for rows.Next() {
		var (
			pvzID   uuid.NullUUID
			pvzCity sql.NullString
			pvzDate sql.NullTime

			receptionID     uuid.NullUUID
			receptionDate   sql.NullTime
			receptionStatus sql.NullString

			productID       uuid.NullUUID
			productDate     sql.NullTime
			productCategory sql.NullString
		)

		err := rows.Scan(
			&pvzID, &pvzCity, &pvzDate,
			&receptionID, &receptionDate, &receptionStatus,
			&productID, &productDate, &productCategory,
		)
		if err != nil {
			loggerVar.Error(err.Error())
			return nil, err
		}

		if !pvzID.Valid {
			continue
		}

		pvzUUID := pvzID.UUID
		pvzKey := pvzUUID.String()

		if _, exists := pvzMap[pvzKey]; !exists {
			pvzMap[pvzKey] = &models.PVZ{
				Id:               pvzUUID,
				City:             pvzCity.String,
				RegistrationDate: pvzDate.Time,
				Receptions:       []models.Reception{},
			}
		}

		if receptionID.Valid {
			receptionUUID := receptionID.UUID
			recFound := false

			for i := range pvzMap[pvzKey].Receptions {
				if pvzMap[pvzKey].Receptions[i].Id == receptionUUID {
					if productID.Valid {
						pvzMap[pvzKey].Receptions[i].Products = append(
							pvzMap[pvzKey].Receptions[i].Products,
							models.Product{
								Id:          productID.UUID,
								DateTime:    productDate.Time,
								Type:        productCategory.String,
								ReceptionId: receptionUUID,
							},
						)
					}
					recFound = true
					break
				}
			}

			if !recFound {
				rec := models.Reception{
					Id:       receptionUUID,
					DateTime: receptionDate.Time,
					PvzId:    pvzUUID,
					Status:   receptionStatus.String,
					Products: []models.Product{},
				}

				if productID.Valid {
					rec.Products = append(rec.Products, models.Product{
						Id:          productID.UUID,
						DateTime:    productDate.Time,
						Type:        productCategory.String,
						ReceptionId: receptionUUID,
					})
				}

				pvzMap[pvzKey].Receptions = append(pvzMap[pvzKey].Receptions, rec)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []models.PVZ
	for _, val := range pvzMap {
		result = append(result, *val)
	}

	loggerVar.Info("Successful")
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

	loggerVar.Info("Successful")
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
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	var exists bool
	err := repo.db.QueryRow(ctx, hasActiveReception, pvzID).Scan(&exists)
	if err != nil {
		loggerVar.Error(err.Error())
	}

	loggerVar.Info("Successful")
	return exists, err
}

func (repo *PvzRepo) GetActiveReception(ctx context.Context, pvzId uuid.UUID) (models.Reception, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	var reception models.Reception
	err := repo.db.QueryRow(ctx, getActiveReception, pvzId).
		Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
	if err == sql.ErrNoRows {
		loggerVar.Error(pvz.ErrNoActiveReception.Error())
		return models.Reception{}, pvz.ErrNoActiveReception
	}
	if err != nil {
		loggerVar.Error(err.Error())
		return models.Reception{}, err
	}

	loggerVar.Info("Successful")
	return reception, nil
}

func (repo *PvzRepo) CreateReception(ctx context.Context, reception models.Reception) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, createReception, reception.Id, reception.DateTime, reception.PvzId, reception.Status)
	if err != nil {
		loggerVar.Error(err.Error())
	}

	loggerVar.Info("Successful")
	return err
}

func (repo *PvzRepo) AddProduct(ctx context.Context, product *models.Product) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertProduct, product.Id, product.DateTime, product.ReceptionId, product.Type)
	if err != nil {
		loggerVar.Error(err.Error())
	}

	loggerVar.Info("Successful")
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

	loggerVar.Info("Successful")
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
		return errors.New("nothing was updated")
	}

	loggerVar.Info("Successful")
	return nil
}
