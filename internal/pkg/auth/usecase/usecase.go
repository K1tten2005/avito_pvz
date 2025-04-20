package usecase

import (
	"context"
	"crypto/rand"
	"log/slog"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"github.com/K1tten2005/avito_pvz/internal/pkg/auth"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/jwtUtils"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/validation"

	"github.com/satori/uuid"
)

type AuthUsecase struct {
	repo auth.AuthRepo
}

func CreateAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (uc *AuthUsecase) DummyLogin(ctx context.Context, data models.DummyLoginReq) (string, string, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	dummyUser := models.User{
		Id: uuid.NewV4(),
		Email: "dummy@user.com",
		Role: data.Role,
	}

	token, err := jwtUtils.GenerateToken(dummyUser)
	if err != nil {
		loggerVar.Error(auth.ErrGeneratingToken.Error())
		return "", "", auth.ErrGeneratingToken
	}

	csrfToken := uuid.NewV4().String()

	loggerVar.Info("Successful")
	return token, csrfToken, nil
}


func (uc *AuthUsecase) Login(ctx context.Context, data models.LoginReq) (models.User, string, string, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	if !validation.ValidEmail(data.Email) {
		loggerVar.Error(auth.ErrInvalidEmail.Error())
		return models.User{}, "", "", auth.ErrInvalidEmail
	}

	user, err := uc.repo.SelectUserByEmail(ctx, data.Email)
	if err != nil {
		loggerVar.Error(auth.ErrUserNotFound.Error())
		return models.User{}, "", "", auth.ErrUserNotFound
	}

	if !validation.CheckPassword(user.PasswordHash, data.Password) {
		loggerVar.Error(auth.ErrInvalidCredentials.Error())
		return models.User{}, "", "", auth.ErrInvalidCredentials
	}

	token, err := jwtUtils.GenerateToken(user)
	if err != nil {
		loggerVar.Error(auth.ErrGeneratingToken.Error())
		return models.User{}, "", "", auth.ErrGeneratingToken
	}

	csrfToken := uuid.NewV4().String()

	loggerVar.Info("Successful")
	return user, token, csrfToken, nil
}

func (uc *AuthUsecase) Register(ctx context.Context, data models.RegisterReq) (models.User, string, string, error) {
	loggerVar := logger.GetLoggerFromContext(ctx).With(slog.String("func", logger.GetFuncName()))

	if !validation.ValidEmail(data.Email) {
		loggerVar.Error(auth.ErrInvalidEmail.Error())
		return models.User{}, "", "", auth.ErrInvalidEmail
	}

	if !validation.ValidPassword(data.Password) {
		loggerVar.Error(auth.ErrInvalidPassword.Error())
		return models.User{}, "", "", auth.ErrInvalidPassword
	}

	salt := make([]byte, 8)
	rand.Read(salt)
	hashedPassword := validation.HashPassword(salt, data.Password)

	newUser := models.User{
		Id:           uuid.NewV4(),
		Email:        data.Email,
		Role:         data.Role,
		PasswordHash: hashedPassword,
	}

	err := uc.repo.InsertUser(ctx, newUser)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.User{}, "", "", auth.ErrCreatingUser
	}

	token, err := jwtUtils.GenerateToken(newUser)
	if err != nil {
		loggerVar.Error(err.Error())
		return models.User{}, "", "", auth.ErrGeneratingToken
	}

	csrfToken := uuid.NewV4().String()

	loggerVar.Info("Successful")
	return newUser, token, csrfToken, nil
}
