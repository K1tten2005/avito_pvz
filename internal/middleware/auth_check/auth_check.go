package authcheck

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/logger"
	"github.com/K1tten2005/avito_pvz/internal/pkg/utils/send_err"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggerVar := logger.GetLoggerFromContext(r.Context()).With(slog.String("func", logger.GetFuncName()))

		_, err := r.Cookie("AvitoJWT")
		if err != nil {
			if err == http.ErrNoCookie {
				logger.LogHandlerError(loggerVar, fmt.Errorf("токен отсутствует: %w", err), http.StatusBadRequest)
				send_err.SendError(w, "токен отсутствует", http.StatusBadRequest)
				return
			}
			logger.LogHandlerError(loggerVar, fmt.Errorf("ошибка при чтении куки: %w", err), http.StatusBadRequest)
			send_err.SendError(w, "ошибка при чтении куки", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
