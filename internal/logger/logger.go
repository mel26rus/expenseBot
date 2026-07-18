package logger

import (
	"io"
	"log/slog"
	"os"
)

func Init(logFile string, isDebug bool) (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	logLevel := slog.LevelInfo
	if isDebug {
		logLevel = slog.LevelDebug
	}
	var logger *slog.Logger
	handler := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger = slog.New(handler)

	// делаем глобальным
	slog.SetDefault(logger)

	return file, nil
}
