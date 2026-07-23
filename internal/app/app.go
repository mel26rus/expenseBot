package app

import (
	"context"
	"log/slog"

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
}

func (a *App) Run(ctx context.Context) {
	a.Logger.Debug("+App.run")
	botAPI, err := tgbotapi.NewBotAPI(a.Config.BotKey)
	if err != nil {
		a.Logger.Error("bot init failed", "error", err)
		return
	}
	slog.Debug("botApi inited")
	go a.Scheduler.Run(ctx)
	slog.Debug("Scheduler started")

	handler := BuildHandler(botAPI, a)
	a.Logger.Debug("+handler.Start(ctx)")
	handler.Start(ctx)
	a.Logger.Debug("-handler.Start(ctx)")
	a.Logger.Debug("-App.run")
}
