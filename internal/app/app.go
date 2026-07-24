package app

import (
	"context"
	"log/slog"

	"expense-bot/internal/bot"
	"expense-bot/internal/config"
	"expense-bot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config    *config.AppConfig
	Logger    *slog.Logger
	DB        *pgxpool.Pool
	Scheduler *scheduler.Scheduler
	botAPI    *tgbotapi.BotAPI
	handler   *bot.Handler
}

// тут всё запускаем впринципе оно и так есть
func (a *App) Run(ctx context.Context) {
	a.Logger.Debug("+App.run")
	go a.Scheduler.Run(ctx)
	slog.Info("Scheduler started")

	a.Logger.Info("+handler.Start(ctx)")
	a.handler.Start(ctx)
	a.Logger.Info("-handler.Start(ctx)")
	a.Logger.Debug("-App.run")
}
