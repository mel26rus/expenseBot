package scheduler

import (
	"context"
	"expense-bot/internal/bot"
	"fmt"
	"log/slog"
	"time"
)

type MonthlyReportJob struct {
	Handler *bot.Handler
}

func NewMonthlyReportJob(handler *bot.Handler) *MonthlyReportJob {
	return &MonthlyReportJob{
		Handler: handler,
	}
}

func (j *MonthlyReportJob) Name() string {
	return "Monthly Report Job"
}

func (j *MonthlyReportJob) NextRun(now time.Time) time.Time {

	//для месячного

	next := time.Date(
		now.Year(),
		now.Month(),
		1,
		12, 0, 0, 0,
		now.Location(),
	)

	if !next.After(now) {
		next = time.Date(
			now.Year(),
			now.Month()+1,
			1,
			12, 0, 0, 0,
			now.Location(),
		)
	}
	return next

	//now.Add(time.Second * 10)
}

func (j *MonthlyReportJob) Run(ctx context.Context) error {
	fn_name := "MonthlyReportJob.Run"
	slog.Debug(fmt.Sprintf("+%s", fn_name))
	j.Handler.HandleMonthlyReports(ctx)
	slog.Debug(fmt.Sprintf("-%s", fn_name))
	return nil
}
