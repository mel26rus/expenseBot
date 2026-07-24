package scheduler

import (
	"context"
	"expense-bot/internal/reportformatter"
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

	return now.Add(time.Minute)
}

func (j *DailyReportJob) Run(ctx context.Context) error {

	users, _ := j.reportService.GetUsersForDailyReport(ctx)

	for _, userID := range users {

		report, err := j.reportService.BuildDailyReport(
			ctx,
			userID,
		)
		if err != nil {
			slog.Error("reportService.BuildDailyReport", "Error", err)
		}

		text := reportformatter.BuildDailyReportText(report)

		msg := tgbotapi.NewMessage(userID, text)
		msg.ParseMode = tgbotapi.ModeHTML

		_, err = j.bot.Send(msg)
		if err != nil {
			slog.Error("Daily report run", "Error", err)
		}
	}

	return nil
}
