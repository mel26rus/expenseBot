package scheduler

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"
)

type scheduledJob struct {
	job     Job
	nextRun time.Time
	running atomic.Bool
}

type Scheduler struct {
	jobs []*scheduledJob
}

func New() *Scheduler {
	return &Scheduler{
		jobs: make([]*scheduledJob, 0),
	}
}

func (s *Scheduler) Add(job Job) {

	s.jobs = append(s.jobs, &scheduledJob{
		job:     job,
		nextRun: job.NextRun(time.Now()),
	})
}

func (s *Scheduler) Run(ctx context.Context) {

	slog.Debug("Scheduler Run")

	for {

		var next *scheduledJob

		for _, job := range s.jobs {

			if next == nil || job.nextRun.Before(next.nextRun) {
				next = job
			}
		}

		if !next.running.CompareAndSwap(false, true) {
			continue
		}

		timer := time.NewTimer(time.Until(next.nextRun))

		select {

		case <-ctx.Done():

			timer.Stop()
			slog.Info("Scheduler stopped")
			return

		case <-timer.C:

			slog.Debug(
				"Running job",
				"name", next.job.Name(),
			)

			next.nextRun = next.job.NextRun(time.Now())

			go func(job *scheduledJob) {

				if err := job.job.Run(ctx); err != nil {
					slog.Error(
						"Job failed",
						"name", job.job.Name(),
						"error", err,
					)
				}

			}(next)
		}
	}
}
