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
				logger.LogHandlerError(loggerVar, fmt.Errorf("no token: %w", err), http.StatusBadRequest)
				send_err.SendError(w, "no token", http.StatusBadRequest)
				return
			}
			logger.LogHandlerError(loggerVar, fmt.Errorf("error while parsing cookie: %w", err), http.StatusBadRequest)
			send_err.SendError(w, "error while parsing cookie", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
