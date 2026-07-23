package scheduler

import (
	"context"
	"log/slog"
	"time"
)

type HourlyLogJob struct {
}

func NewHourlyLogJob() *HourlyLogJob {
	return &HourlyLogJob{}
}

func (j *HourlyLogJob) Name() string {
	return "Hourly Log Job"
}

func (j *HourlyLogJob) NextRun(now time.Time) time.Time {

	return now.Add(time.Hour)
}

func (j *HourlyLogJob) Run(ctx context.Context) error {

	slog.Info(
		"Scheduler is alive",
		"time",
		time.Now().Format("15:04:05"),
	)

	return nil
}
