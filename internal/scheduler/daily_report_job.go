package scheduler

import (
	"context"
	"expense-bot/internal/service"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DailyReportJob struct {
	logger        *slog.Logger
	bot           *tgbotapi.BotAPI
	reportService *service.ReportService
}

func NewDailyReportJob() *DailyReportJob {
	return &DailyReportJob{}
}

func (j *DailyReportJob) Name() string {
	return "Daily Report Job"
}

func (j *DailyReportJob) NextRun(now time.Time) time.Time {

	return now.Add(time.Hour)
}

func (j *DailyReportJob) Run(ctx context.Context) error {

	users, _ := j.reportService.GetUsersHasTransactionsDaily(ctx)

	for _, userID := range users {

		report, _ := j.reportService.BuildDailyReport(
			ctx,
			userID,
			start,
			end,
		)

		text := BuildDailyReportText(report)

		msg := tgbotapi.NewMessage(report.TgUserID, text)
		msg.ParseMode = tgbotapi.ModeHTML

		err := j.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
