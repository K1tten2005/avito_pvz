package acl

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/jwtUtils"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/send_err"
	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt"
)

var Enforcer *casbin.Enforcer

func InitACL(logger *slog.Logger) error {
	modelPath := "internal/middleware/acl/model.conf"
	policyPath := "internal/middleware/acl/policy.csv"

	e, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	Enforcer = e
	logger.Info("Успешный запуск ACL")
	return nil
}

func ACLMiddleware(next http.Handler) http.Handler {
	secret := os.Getenv("JWT_SECRET")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))
		cookieJWT, err := r.Cookie("AvitoJWT")
		if err != nil {
			if err == http.ErrNoCookie {
				logger.LogHandlerError(loggerVar, errors.New("нет куки jwt"), http.StatusForbidden)
				send_err.SendError(w, "нет куки jwt", http.StatusForbidden)
				return
			}
			send_err.SendError(w, "ошибка куки jwt", http.StatusBadRequest)
		}

		JWTStr := cookieJWT.Value
		claims := jwt.MapClaims{}

		role, ok := jwtUtils.GetRoleFromJWT(JWTStr, claims, secret)
		if !ok || role == "" {
			logger.LogHandlerError(loggerVar, errors.New("нет роли"), http.StatusForbidden)
			send_err.SendError(w, "нет роли", http.StatusForbidden)
			return
		}
		
		path := r.URL.Path
		method := r.Method

		allowed, err := Enforcer.Enforce(role, path, method)
		if err != nil {
			logger.LogHandlerError(loggerVar, errors.New("ошибка enforce"), http.StatusInternalServerError)
			send_err.SendError(w, "ошибка enforce", http.StatusInternalServerError)
			return
		}
		if !allowed {
			logger.LogHandlerError(loggerVar, errors.New("нет прав доступа"), http.StatusForbidden)
			send_err.SendError(w, "нет прав доступа", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
