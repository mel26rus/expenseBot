package flow

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/service"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type ReportFlow struct {
	reportService  *service.ReportService
	sessionService *service.SessionService
	userService    *service.UserService
}

func NewReportFlow(
	reportService *service.ReportService,
	sessionService *service.SessionService,
	userService *service.UserService,
) *ReportFlow {
	return &ReportFlow{
		reportService:  reportService,
		sessionService: sessionService,
		userService:    userService,
	}
}

func (r *ReportFlow) GetExistsTxUser(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.User, error) {
	return r.reportService.GetExistsTxUser(ctx, startDate, endDate)
}

func (r *ReportFlow) BuildUserReport(ctx context.Context, userID int64, start time.Time, end time.Time) (Response, error) {

	session, err := r.sessionService.GetUserSession(ctx, userID)
	slog.Debug("BuildUserReport_2", "session", session)
	report, err := r.reportService.BuildUserReport(
		ctx,
		userID,
		start,
		end,
	)
	if err != nil {
		slog.Error("reportService.BuildUserReport", "Error", err)
		return Response{}, err
	}

	text := r.BuildReportText(report, start, end)
	res := Response{
		Text:              text,
		IsSendMenuMessage: true,
		Keyboard:          nil,
		EditMessageId:     session.EditMessageId,
	}
	return res, nil
}

func (r *ReportFlow) BuildReportText(report []*model.AccountReport, start time.Time, end time.Time) string {

	var sb strings.Builder

	duration := end.Sub(start)
	days := int(duration.Hours() / 24)
	if days == 1 {
		sb.WriteString(
			fmt.Sprintf(
				"📅 <b>Отчет за %s</b>\n",
				start.Format("02.01.2006"),
			),
		)
	} else {
		startTitle := fmt.Sprintf("%s %s", getRussianMonthName(int(start.Month())), strconv.Itoa(start.Year()))

		sb.WriteString(
			fmt.Sprintf(
				"📅 <b>Отчет за %s</b>\n",
				startTitle,
			),
		)
	}

	var totRUBAmount float64 = 0.00
	var totUSDAmount float64 = 0.00
	var erDate time.Time
	var exRate float64 = 0.00
	for _, account := range report {

		sb.WriteString(
			fmt.Sprintf(
				"💳 <b>%s</b> %s\n",
				account.Title,
				account.CurrencyName,
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
						formatAmount(tx.Amount),
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
					formatAmount(-tx.Amount),
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
					"📈 Итого доход: <b>%s</b>\n",
					formatAmount(account.Income),
				),
			)
		}

		if account.Expense > 0 {

			sb.WriteString(
				fmt.Sprintf(
					"📉 Итого расход: <b>%s</b>\n",
					formatAmount(account.Expense),
				),
			)
		}

		sb.WriteString(
			fmt.Sprintf(
				"💰 Баланс: <code>%s</code> %s\n",
				formatAmount(account.Balance),
				account.CurrencyName,
			),
		)

		// if (account.CurrencyName != "USD") && (account.CurrencyName != "USDT") {
		// 	sb.WriteString(
		// 		fmt.Sprintf(
		// 			"💰 Баланс: <code>%s</code> USD\n",
		// 			formatAmount(account.USDBalance),
		// 		),
		// 	)
		// }

		// if account.CurrencyName != "RUB" {
		// 	sb.WriteString(
		// 		fmt.Sprintf(
		// 			"💰 Баланс: <code>%s</code> RUB\n",
		// 			formatAmount(account.RUBBalance),
		// 		),
		// 	)
		// }
		if account.CurrencyName == "RUB" {
			exRate = account.ExRate
		}
		totRUBAmount = totRUBAmount + account.RUBBalance
		totUSDAmount = totUSDAmount + account.USDBalance
		erDate = account.ExDate
		sb.WriteString("────────────────────\n")
	}

	sb.WriteString(
		fmt.Sprintf(
			"USD/RUB от %s: <code>%s</code> \n",
			erDate.Format("02.01.2006"),
			formatAmount(exRate),
		),
	)
	sb.WriteString(
		fmt.Sprintf(
			"💰 Общий: <code>%s</code> RUB\n",
			formatAmount(totRUBAmount),
		),
	)
	sb.WriteString(
		fmt.Sprintf(
			"💰 Общий: <code>%s</code> USD\n",
			formatAmount(totUSDAmount),
		),
	)
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

func getRussianMonthName(monNum int) string {
	switch monNum {
	case 1:
		return "Январь"
	case 2:
		return "Февраль"
	case 3:
		return "Март"
	case 4:
		return "Апрель"
	case 5:
		return "Май"
	case 6:
		return "Июнь"
	case 7:
		return "Июль"
	case 8:
		return "Август"
	case 9:
		return "Сентябрь"
	case 10:
		return "Октябрь"
	case 11:
		return "Ноябрь"
	case 12:
		return "Декабрь"
	default:
		return "неизвестный номер месяца"
	}
}
