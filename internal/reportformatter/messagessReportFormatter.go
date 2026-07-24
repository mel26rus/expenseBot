package reportformatter

import (
	"expense-bot/internal/flow"
	"expense-bot/internal/model"
	"fmt"
	"strings"
	"time"
)

func BuildDailyReportText(report []*model.AccountReport) string {

	var sb strings.Builder

	sb.WriteString(
		fmt.Sprintf(
			"📅 <b>Отчет за %s</b>\n\n",
			time.Now().AddDate(0, 0, -1).Format("02.01.2006"),
		),
	)

	for _, account := range report {

		sb.WriteString(
			fmt.Sprintf(
				"💳 <b>%s</b>\n",
				account.Title,
			),
		)

		sb.WriteString(
			fmt.Sprintf(
				"💰 Баланс: <code>%s</code>\n\n",
				flow.FormatAmount(account.Balance),
			),
		)

		hasIncome := false
		hasExpense := false

		for _, tx := range account.TransactionsReport {

			if tx.Amount > 0 {

				if !hasIncome {
					sb.WriteString("📈 <b>Доходы</b>\n")
					hasIncome = true
				}

				sb.WriteString(
					fmt.Sprintf(
						"• %-18s <b>+%s</b>\n",
						emoji(tx.Category)+capitalize(tx.Category),
						flow.FormatAmount(tx.Amount),
					),
				)

			} else {

				if !hasExpense {
					sb.WriteString("\n📉 <b>Расходы</b>\n")
					hasExpense = true
				}

				sb.WriteString(
					fmt.Sprintf(
						"• %-18s <b>%s</b>\n",
						emoji(tx.Category)+capitalize(tx.Category),
						flow.FormatAmount(-tx.Amount),
					),
				)
			}
		}

		sb.WriteString("\n")

		if account.Income > 0 {

			sb.WriteString(
				fmt.Sprintf(
					"📈 Доход: <b>%s</b>\n",
					flow.FormatAmount(account.Income),
				),
			)
		}

		if account.Expense > 0 {

			sb.WriteString(
				fmt.Sprintf(
					"📉 Расход: <b>%s</b>\n",
					flow.FormatAmount(account.Expense),
				),
			)
		}

		sb.WriteString("\n────────────────────\n\n")
	}

	return sb.String()
}

func capitalize(s string) string {

	if s == "" {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func emoji(category string) string {

	switch strings.ToLower(category) {

	case "обед", "еда":
		return "🍔 "

	case "кофе":
		return "☕ "

	case "бензин":
		return "⛽ "

	case "зарплата":
		return "💼 "

	case "подарок":
		return "🎁 "

	case "такси":
		return "🚕 "

	case "магазин":
		return "🛒 "

	case "квартира":
		return "🏠 "

	case "интернет":
		return "🌐 "

	case "телефон":
		return "📱 "

	case "развлечения":
		return "🎮 "

	default:
		return "• "
	}
}
