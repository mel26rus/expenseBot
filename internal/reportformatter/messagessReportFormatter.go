package reportformatter

import (
	"expense-bot/internal/flow"
	"expense-bot/internal/model"
	"fmt"
	"strings"
	"time"
	"unicode"
)

func BuildDailyReportText(report []*model.AccountReport) string {

	var sb strings.Builder

	sb.WriteString(
		fmt.Sprintf(
			"📅 <b>Отчет за %s</b>\n",
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
				"💰 Баланс: <code>%s</code>\n",
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
					sb.WriteString("📉 <b>Расходы</b>\n")
					hasExpense = true
				}

				text := fmt.Sprintf(
					"• %-18s <b>%s</b>\n",
					emoji(tx.Category)+capitalize(tx.Category),
					flow.FormatAmount(-tx.Amount),
				)

				sb.WriteString(
					text,
				)
			}
		}

		//	sb.WriteString("\n")

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

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return fmt.Sprintf("%s", string(r))
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
