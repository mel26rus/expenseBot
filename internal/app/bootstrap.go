package app

import (
	"expense-bot/internal/bot"
	"expense-bot/internal/flow"
	"expense-bot/internal/repository"
	"expense-bot/internal/scheduler"
	"expense-bot/internal/service"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *App) Boot() {
	slog.Info("Booting application")

	var err error
	a.botAPI, err = tgbotapi.NewBotAPI(a.Config.BotKey)
	if err != nil {
		a.Logger.Error("bot init failed", "error", err)
		return
	}
	slog.Info("botApi inited")
	// repos
	userRepo := repository.NewUserRepo(a.DB)
	accountRepo := repository.NewAccountRepo(a.DB)
	sessionRepo := repository.NewSessionRepo(a.DB)
	currencyRepo := repository.NewCurrencyRepo(a.DB)
	transactionRepo := repository.NewTransactionRepo(a.DB)
	reportRepo := repository.NewReportRepo(a.DB)
	slog.Info("Rpositories inited")
	// services
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	sessionService := service.NewSessionService(sessionRepo)
	currencyService := service.NewCurrencyService(currencyRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	reportService := service.NewReportService(reportRepo)
	slog.Info("Services inited")
	// flow

	accountFlow := flow.NewAccountFlow(
		userService,
		accountService,
		currencyService,
		sessionService,
	)

	transactionFlow := flow.NewTransactionFlow(
		userService,
		accountService,
		transactionService,
		sessionService,
	)

	mainFlow := flow.NewMainFlow(
		sessionService,
		accountFlow,
		transactionFlow,
		flow.NewMenuFlow(userService),
		flow.NewReportFlow(reportService, sessionService, userService),
	)
	slog.Info("Flows inited")

	a.handler = bot.NewHandler(a.botAPI, mainFlow)
	reportHandler := bot.NewHandler(a.botAPI, mainFlow)

	sch := scheduler.New()

	sch.Add(
		scheduler.NewHourlyLogJob(),
	)
	sch.Add(
		scheduler.NewDailyReportJob(reportHandler),
	)

	sch.Add(
		scheduler.NewMonthlyReportJob(reportHandler),
	)

	a.Scheduler = sch

	slog.Info("Boot aplication completed")
}
