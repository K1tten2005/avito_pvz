package logger

import (
	"context"
	"errors"
	"os"
	"testing"

	"log/slog"

	"github.com/K1tten2005/avito_pvz/internal/middleware/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetFuncName(t *testing.T) {
	funcName := GetFuncName()
	assert.NotEmpty(t, funcName)
	assert.Equal(t, "logger.TestGetFuncName", funcName) 
}

func TestLogHandlerInfo(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	LogHandlerInfo(logger, "Test info message", 200)
	assert.NotNil(t, logger)
}

func TestLogHandlerError(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	err := errors.New("Test error message")
	LogHandlerError(logger, err, 500)

	assert.NotNil(t, logger)

	unwrappedErr := errors.Unwrap(err)
	assert.Nil(t, unwrappedErr) 
}

func TestGetLoggerFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), logger.LoggerKey, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))
	logger := GetLoggerFromContext(ctx)

	assert.NotNil(t, logger)

	ctxWithoutLogger := context.Background()
	loggerWithoutContext := GetLoggerFromContext(ctxWithoutLogger)

	assert.NotNil(t, loggerWithoutContext)
}
