package app

import (
	"expense-bot/internal/bot"
	"expense-bot/internal/flow"
	"expense-bot/internal/repository"
	"expense-bot/internal/service"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BuildHandler(botAPI *tgbotapi.BotAPI, app *App) *bot.Handler {
	slog.Debug("Building bot handler")
	// repos
	userRepo := repository.NewUserRepo(app.DB)
	accountRepo := repository.NewAccountRepo(app.DB)
	sessionRepo := repository.NewSessionRepo(app.DB)
	currencyRepo := repository.NewCurrencyRepo(app.DB)
	transactionRepo := repository.NewTransactionRepo(app.DB)

	// services
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	sessionService := service.NewSessionService(sessionRepo)
	currencyService := service.NewCurrencyService(currencyRepo)
	transactionService := service.NewTransactionService(transactionRepo)

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
	)

	// handler
	slog.Debug("Bot handler built")
	return bot.NewHandler(botAPI, mainFlow)
}
