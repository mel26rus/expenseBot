package scheduler

import (
	"context"
	"expense-bot/internal/bot"
	"fmt"
	"log/slog"
	"time"
)

type DailyReportJob struct {
	Handler *bot.Handler
}

func NewDailyReportJob(handler *bot.Handler) *DailyReportJob {
	return &DailyReportJob{
		Handler: handler,
	}
}

func (j *DailyReportJob) Name() string {
	return "Daily Report Job"
}

func (j *DailyReportJob) NextRun(now time.Time) time.Time {

	next := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		6, 0, 0, 0,
		now.Location(),
	)

	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}

	return next
}

func (j *DailyReportJob) Run(ctx context.Context) error {
	fn_name := "DailyReportJob.Run"
	slog.Debug(fmt.Sprintf("+%s", fn_name))
	j.Handler.HandleDailyReports(ctx)
	slog.Debug(fmt.Sprintf("-%s", fn_name))
	return nil
}
