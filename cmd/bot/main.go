package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kardianos/service"
	"golang.org/x/sys/windows/svc"

	"expense-bot/internal/app"
	"expense-bot/internal/config"
	"expense-bot/internal/db"
	"expense-bot/internal/logger"
	"expense-bot/internal/servicehost"
)

var configPath = flag.String(
	"config",
	"config.yaml",
	"Path to config file",
)

func main() {
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	err = logger.New(cfg)
	if err != nil {
		panic(err)
	}
	slog.Info("Logger inited")

	dbPool := db.NewPostgres(cfg.Database.URL)
	slog.Info("dbPool inited")
	defer dbPool.Close()
	defer slog.Info("dbPool closed")

	application := &app.App{
		Config: cfg,
		DB:     dbPool,
		Logger: slog.Default(),
	}

	application.Boot()

	svcConfig := &service.Config{
		Name:        "ExpenseBot",
		DisplayName: "Expense Bot",
		Description: "Telegram Expense Bot",
	}

	program := &servicehost.Program{
		App: application,
	}

	isService, err := svc.IsWindowsService()
	if err != nil {
		slog.Error("failed to detect service mode: ", "Error", err)
		return
	}

	svc, err := service.New(program, svcConfig)
	if err != nil {
		panic(err)
	}

	if len(flag.Args()) > 0 {
		err = nil
		switch flag.Args()[0] {

		case "install":
			err = svc.Install()

		case "uninstall":
			err = svc.Uninstall()
		}
		if err != nil {
			panic(err)
		}
		return
	}

	if !isService {

		ctx, stop := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)
		defer stop()

		application.Run(ctx)

		return
	}

	err = svc.Run()
	if err != nil {
		panic(err)
	}

	slog.Info("Bot closed")
}
