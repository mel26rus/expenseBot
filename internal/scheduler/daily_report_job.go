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
	bot           *tgbotapi.BotAPI
	reportService *service.ReportService
}

func NewDailyReportJob(bot *tgbotapi.BotAPI, reportService *service.ReportService) *DailyReportJob {
	return &DailyReportJob{
		bot:           bot,
		reportService: reportService,
	}
}

func (j *DailyReportJob) Name() string {
	return "Daily Report Job"
}

func (j *DailyReportJob) NextRun(now time.Time) time.Time {

	return now.Add(time.Hour)
}

func (j *DailyReportJob) Run(ctx context.Context) error {
	slog.Debug("DailyReportJob.Run")
	users, _ := j.reportService.GetUsersForDailyReport(ctx)

	for _, userTGID := range users {

		report, err := j.reportService.BuildDailyReport(
			ctx,
			userTGID,
		)
		if err != nil {
			slog.Error("reportService.BuildDailyReport", "Error", err)
		}

		text := reportformatter.BuildDailyReportText(report)
		slog.Debug("Report", "userTGID", userTGID)
		slog.Debug("Report", "text", text)
		msg := tgbotapi.NewMessage(userTGID, text)
		msg.ParseMode = tgbotapi.ModeHTML

		_, err = j.bot.Send(msg)
		if err != nil {
			slog.Error("Daily report run", "Error", err)
		}
	}

	return nil
}
