package flow

import (
	"expense-bot/internal/model"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func buildCurrencyInline(list []*model.Currency) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, c := range list {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			c.Code,
			fmt.Sprintf("currency:%d", c.ID),
		)

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	ikb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &ikb
}

func buildTypeInline() *tgbotapi.InlineKeyboardMarkup {
	ikb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("💸 "+constTxTypeExpense, "txtype:-1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("💰 "+constTxTypeIncome, "txtype:1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔁 "+constTxTypeTransfer, "txtype:0")),
	)
	slog.Debug("Built type inline keyboard")
	return &ikb
}

func buildCommentInline(list []model.Comments) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, c := range list {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			c.Comment,
			"comment:"+fmt.Sprintf("%d", c.TransactionID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	ikb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &ikb
}

func buildUserAccountsInline(list []*model.Account) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range list {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			c.Name,
			fmt.Sprintf("account:%d", c.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	ikb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &ikb
}

func buildMenuSettingsInline(isDaily bool, isMonthly bool) *tgbotapi.InlineKeyboardMarkup {
	ikb := tgbotapi.NewInlineKeyboardMarkup()
	if isDaily {
		ikb.InlineKeyboard = append(ikb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ежедневный отчет 🟢", constDailyReportChange)))
	} else {
		ikb.InlineKeyboard = append(ikb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ежедневный отчет 🔴", constDailyReportChange)))
	}
	if isMonthly {
		ikb.InlineKeyboard = append(ikb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ежемесячный отчет 🟢", constMonthlyReportChange)))
	} else {
		ikb.InlineKeyboard = append(ikb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Ежемесячный отчет 🔴", constMonthlyReportChange)))
	}
	return &ikb
}

func buildMenuMainInline() *tgbotapi.InlineKeyboardMarkup {
	ikb := tgbotapi.NewInlineKeyboardMarkup()
	ikb.InlineKeyboard = append(ikb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Отчеты", constMenuReports)))
	ikb.InlineKeyboard = append(ikb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🧾 Счета", constMenuAccounts)))
	return &ikb
}

func buildMenuReportsInline() *tgbotapi.InlineKeyboardMarkup {
	ikb := tgbotapi.NewInlineKeyboardMarkup()
	ikb.InlineKeyboard = append(ikb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Потрачено сегодня", constMenuTodayReport)))
	ikb.InlineKeyboard = append(ikb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🧾 Потрачено текущий месяц", constMenuCurrMonthReport)))
	return &ikb
}

func buildMenuUserAccountsInline(list []*model.Account) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, c := range list {
		var ico string
		if c.Is_hidden {
			ico = "🙈"
		} else {
			ico = "👁️"
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", c.Name, ico),
			fmt.Sprintf("%s%d", constHideAccountChange, c.ID),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	ikb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &ikb
}
