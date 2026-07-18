package app

import (
	"log/slog"

	"expense-bot/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config *config.AppConfig
	Logger *slog.Logger
	DB     *pgxpool.Pool
}
