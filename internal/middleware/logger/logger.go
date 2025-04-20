package logger

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/satori/uuid"
	"github.com/gorilla/mux"
)

type ctxKey string

const LoggerKey ctxKey = "logger"

func CreateLoggerMiddleware(loggerVar *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), LoggerKey, loggerVar.With(slog.String("ID", uuid.NewV4().String())))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
