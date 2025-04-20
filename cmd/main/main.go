package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authCheck "github.com/K1tten2005/avito_pvz/internal/middleware/auth_check"
	"github.com/K1tten2005/avito_pvz/internal/middleware/cors"
	"github.com/K1tten2005/avito_pvz/internal/middleware/csp"
	"github.com/K1tten2005/avito_pvz/internal/middleware/logger"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	authHandler "github.com/K1tten2005/avito_pvz/internal/pkg/auth/delivery/http"
	authRepo "github.com/K1tten2005/avito_pvz/internal/pkg/auth/repo"
	authUsecase "github.com/K1tten2005/avito_pvz/internal/pkg/auth/usecase"
	pvzHandler "github.com/K1tten2005/avito_pvz/internal/pkg/pvz/delivery/http"
	pvzRepo "github.com/K1tten2005/avito_pvz/internal/pkg/pvz/repo"
	pvzUsecase "github.com/K1tten2005/avito_pvz/internal/pkg/pvz/usecase"
)

func initDB(logger *slog.Logger) (*pgxpool.Pool, error) {
	connStr := os.Getenv("POSTGRES_CONNECTION")

	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("Успешное подключение к PostgreSQL")
	return pool, nil
}

func main() {
	logFile, err := os.OpenFile(os.Getenv("MAIN_LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("error opening log file: " + err.Error())
		return
	}
	defer logFile.Close()

	loggerVar := slog.New(slog.NewJSONHandler(io.MultiWriter(logFile, os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))

	// acl.InitACL(loggerVar)

	pool, err := initDB(loggerVar)
	if err != nil {
		loggerVar.Error("Ошибка при подключении к PostgreSQL: " + err.Error())
		return
	}
	defer pool.Close()

	logMW := logger.CreateLoggerMiddleware(loggerVar)

	authRepo := authRepo.CreateAuthRepo(pool)
	authUsecase := authUsecase.CreateAuthUsecase(authRepo)
	authHandler := authHandler.CreateAuthHandler(authUsecase)

	pvzRepo := pvzRepo.CreatePvzRepo(pool)
	pvzUsecase := pvzUsecase.CreatePvzUsecase(pvzRepo)
	pvzHandler := pvzHandler.CreatePvzHandler(pvzUsecase)

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Не найдено", http.StatusNotFound)
	})
	r.Use(
		logMW,
		cors.CorsMiddleware,
		csp.CspMiddleware,
	)

	// Публичные маршруты
	publicRoutes := r.PathPrefix("/").Subrouter()

	publicRoutes.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	publicRoutes.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	publicRoutes.HandleFunc("/dummyLogin", authHandler.DummyLogin).Methods(http.MethodPost)

	// Защищенные маршруты
	protectedRoutes := r.NewRoute().Subrouter()
	protectedRoutes.Use(
		logMW,
		cors.CorsMiddleware,
		csp.CspMiddleware,
		authCheck.AuthMiddleware,
		//acl.ACLMiddleware, // Оставим отключенным пока
	)

	protectedRoutes.HandleFunc("/pvz", pvzHandler.CreatePvz).Methods(http.MethodPost)
	protectedRoutes.HandleFunc("/pvz", pvzHandler.GetPvz).Methods(http.MethodGet)

	srv := http.Server{
		Handler:           r,
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			loggerVar.Error("Ошибка при запуске сервера: " + err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	loggerVar.Info("Получен сигнал остановки")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		loggerVar.Error("Ошибка при остановке сервера: " + err.Error())
	} else {
		loggerVar.Info("Сервер успешно остановлен")
	}
}
