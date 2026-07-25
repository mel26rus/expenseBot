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

	return now.Add(time.Minute)
}

func (j *DailyReportJob) Run(ctx context.Context) error {
	fn_name := "DailyReportJob.Run"
	slog.Debug(fmt.Sprintf("+%s", fn_name))
	j.Handler.HandleDailyReports(ctx)
	slog.Debug(fmt.Sprintf("-%s", fn_name))
	return nil
}
