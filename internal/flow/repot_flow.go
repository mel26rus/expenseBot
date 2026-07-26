package flow

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/service"
	"fmt"
	"log/slog"
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

func (r *ReportFlow) GetYesterdayTxUserIDs(ctx context.Context) ([]int64, error) {
	return r.reportService.GetUsersIDYesterdayReport(ctx)
}

func (r *ReportFlow) BuildUserDailyReport(ctx context.Context, userTGID int64) (Response, int64, error) {

	user, err := r.userService.GetOrCreate(ctx, userTGID)
	session, err := r.sessionService.GetUserSession(ctx, user.ID)
	report, err := r.reportService.BuildDailyReport(
		ctx,
		user.TelegramID,
	)
	if err != nil {
		slog.Error("reportService.BuildDailyReport", "Error", err)
		return Response{}, user.TelegramID, err
	}

	text := r.BuildDailyReportText(report)
	res := Response{
		Text:              text,
		IsSendMenuMessage: true,
		Keyboard:          nil,
		EditMessageId:     session.EditMessageId,
	}
	return res, user.TelegramID, nil
}

func (r *ReportFlow) BuildDailyReports(ctx context.Context, session model.Session) ([]Response, error) {

	var resArr []Response
	usersIDs, _ := r.reportService.GetUsersIDYesterdayReport(ctx)

	for _, userID := range usersIDs {

		report, err := r.reportService.BuildDailyReport(
			ctx,
			userID,
		)
		if err != nil {
			slog.Error("reportService.BuildDailyReport", "Error", err)
			return nil, err
		}
		text := r.BuildDailyReportText(report)
		res := Response{
			Text:              text,
			IsSendMenuMessage: true,
			EditMessageId:     session.EditMessageId,
			Keyboard:          nil,
		}
		resArr = append(resArr, res)
	}

	return resArr, nil

}

func (r *ReportFlow) BuildDailyReportText(report []*model.AccountReport) string {

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
				formatAmount(account.Balance),
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
					"📈 Доход: <b>%s</b>\n",
					formatAmount(account.Income),
				),
			)
		}

		if account.Expense > 0 {

			sb.WriteString(
				fmt.Sprintf(
					"📉 Расход: <b>%s</b>\n",
					formatAmount(account.Expense),
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
