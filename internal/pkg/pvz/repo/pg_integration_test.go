package repo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	dsn := os.Getenv("POSTGRES_CONNECTION")
	pool, err := pgxpool.Connect(context.Background(), dsn)
	require.NoError(t, err)
	return pool
}

func TestIntegration_PvzReceptionFlow(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := CreatePvzRepo(db)

	pvzID := uuid.NewV4()
	pvz := models.PVZ{
		Id:               pvzID,
		City:             "Москва",
		RegistrationDate: time.Now(),
	}
	err := repo.InsertPvz(ctx, pvz)
	require.NoError(t, err)

	receptionID := uuid.NewV4()
	reception := models.Reception{
		Id:       receptionID,
		DateTime: time.Now(),
		PvzId:    pvzID,
		Status:   "active",
	}
	err = repo.InsertReception(ctx, reception)
	require.NoError(t, err)

	for i := 0; i < 50; i++ {
		product := models.Product{
			Id:          uuid.NewV4(),
			DateTime:    time.Now(),
			ReceptionId: receptionID,
			Type:        "одежда",
		}
		err := repo.InsertProduct(ctx, product)
		require.NoError(t, err)
	}

	err = repo.UpdateReceptionStatus(ctx, receptionID, "closed")
	require.NoError(t, err)

	receptionFromDB, err := repo.GetReceptionByID(ctx, receptionID)
	require.NoError(t, err)
	require.Equal(t, "closed", receptionFromDB.Status)

	pvzs, err := repo.GetPvz(ctx, nil, nil, 1, 10)
	require.NoError(t, err)

	var found bool
	for _, p := range pvzs {
		if p.Id == pvzID {
			for _, r := range p.Receptions {
				if r.Id == receptionID {
					require.Len(t, r.Products, 50)
					found = true
				}
			}
		}
	}
	require.True(t, found, "Reception with products is not found")
}
