package repo

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
	"github.com/jackc/pgtype/pgxtype"
)

//go:embed sql/insertUser.sql
var insertUser string

//go:embed sql/selectUserByEmail.sql
var selectUserByEmail string

type AuthRepo struct {
	db pgxtype.Querier
}

func CreateAuthRepo(db pgxtype.Querier) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (repo *AuthRepo) InsertUser(ctx context.Context, user models.User) error {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	_, err := repo.db.Exec(ctx, insertUser, user.Id, user.Email, user.Role, user.PasswordHash)
	if err != nil {
		loggerVar.Error(err.Error())
		return err
	}
	loggerVar.Info("Successful")
	return nil
}

func (repo *AuthRepo) SelectUserByEmail(ctx context.Context, email string) (models.User, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	resultUser := models.User{Email: email}
	err := repo.db.QueryRow(ctx, selectUserByEmail, email).Scan(
		&resultUser.Id,
		&resultUser.Role,
		&resultUser.PasswordHash,
	)

	if err != nil {
		loggerVar.Error(err.Error())
		return models.User{}, err
	}
	resultUser.Sanitize()

	loggerVar.Info("Successful")
	return resultUser, nil
}
