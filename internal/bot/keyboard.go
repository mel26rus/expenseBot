package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func TypeKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Расход", "type_expense"),
			tgbotapi.NewInlineKeyboardButtonData("Доход", "type_income"),
		//	tgbotapi.NewInlineKeyboardButtonData("Перевод", "type_transfer"),
		),
	)
}
