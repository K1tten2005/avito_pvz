package repo

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
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
								Id:            productUUID,
								ReceptionTime: productDate,
								Category:      productCategory,
								ReceptionId:   receptionUUID,
							})
					}
					recFound = true
					break
				}
			}

			if !recFound {
				rec := models.Reception{
					Id:            receptionUUID,
					ReceptionTime: receptionDate,
					PvzId:         pvzUUID,
					Status:        receptionStatus,
				}

				if productID.Status == pgtype.Present {
					productIDStr := uuidToString(productID)
					productUUID, err := uuid.FromString(productIDStr)
					if err != nil {
						return nil, err
					}
					rec.Products = append(rec.Products, models.Product{
						Id:            productUUID,
						ReceptionTime: productDate,
						Category:      productCategory,
						ReceptionId:   receptionUUID,
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

func (repo *PvzRepo) GetReceptionByID(ctx context.Context, id uuid.UUID) (*models.Reception, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	var reception models.Reception
	err := repo.db.QueryRow(ctx, getReceptionById, id).
		Scan(&reception.Id, &reception.ReceptionTime, &reception.PvzId, &reception.Status)
	if err != nil {
		loggerVar.Error(err.Error())
		return nil, err
	}
	return &reception, nil
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

func (repo *PvzRepo) InsertReception(ctx context.Context, reception models.Reception) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertReception, reception.Id, reception.ReceptionTime, reception.PvzId, reception.Status)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}
	loggerVar.Info("Successful")
	return nil
}

func (repo *PvzRepo) InsertProduct(ctx context.Context, product models.Product) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertProduct, product.Id, product.ReceptionTime, product.ReceptionId, product.Category)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}
	loggerVar.Info("Successful")
	return nil
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
		loggerVar.Error("продукт не найден")
		return errors.New("продукт не найден")
	}

	loggerVar.Info("Successful")
	return nil
}
