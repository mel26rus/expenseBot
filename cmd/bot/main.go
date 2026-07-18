package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"expense-bot/internal/app"
	"expense-bot/internal/config"
	"expense-bot/internal/db"
	"expense-bot/internal/logger"
)

type Bot struct {
	Config *config.AppConfig
	Logger *slog.Logger
	DB     *pgxpool.Pool
}

var configPath = flag.String(
	"config",
	"config.yaml",
	"Path to config file",
)

func run(ctx context.Context, bot *Bot) {

	botAPI, err := tgbotapi.NewBotAPI(bot.Config.BotKey)
	if err != nil {
		bot.Logger.Error("bot init failed", "error", err)
		return
	}
	slog.Debug("botApi inited")

	handler := app.BuildHandler(botAPI, &app.App{
		Config: bot.Config,
		Logger: bot.Logger,
		DB:     bot.DB,
	})

	handler.Start(ctx)
}

func main() {
	fmt.Println("Expense Bot initializing...")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	logFile, err := logger.Init(cfg.LogFile, cfg.IsDebug)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	dbPool := db.NewPostgres(cfg.Database.URL)
	slog.Info("dbPool inited")
	defer dbPool.Close()
	defer slog.Info("dbPool closed")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	bot := &Bot{
		Config: cfg,
		Logger: slog.Default(),
		DB:     dbPool,
	}

	run(ctx, bot)

	slog.Info("Bot closed")
}
