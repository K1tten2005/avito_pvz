package auth

import (
	"context"
	"errors"

	"github.com/K1tten2005/avito_pvz/internal/models"
)

var (
	ErrCreatingUser       = errors.New("ошибка в создании пользователя")
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInvalidEmail       = errors.New("неверный формат логина")
	ErrInvalidPassword    = errors.New("неверный формат пароля")
	ErrInvalidCredentials = errors.New("неверный логин или пароль")
	ErrGeneratingToken    = errors.New("ошибка генерации токена")
	ErrUUID               = errors.New("ошибка создания UUID")
	ErrSamePassword       = errors.New("новый пароль совпадает со старым")
	ErrDBError            = errors.New("ошибка БД")
	ErrAddressNotFound    = errors.New("ошибка поиска адреса")
)

type AuthRepo interface {
	InsertUser(ctx context.Context, user models.User) error
	SelectUserByEmail(ctx context.Context, email string) (models.User, error)
}

type AuthUsecase interface {
	Login(ctx context.Context, data models.LoginReq) (models.User, string, string, error)
	Register(ctx context.Context, data models.RegisterReq) (models.User, string, string, error)
	DummyLogin(ctx context.Context, data models.DummyLoginReq) (string, string, error)
}
