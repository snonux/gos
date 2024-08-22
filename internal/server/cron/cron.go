package cron

import (
	"context"
	"log"
	"time"

	config "codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/server/health"
	"codeberg.org/snonux/gos/internal/server/repository"
	"codeberg.org/snonux/gos/internal/server/scheduler"
)

func Run(ctx context.Context, conf config.ServerConfig, status health.Status) {
	helloTicker := time.NewTicker(time.Hour)
	mergeTicker := time.NewTicker(time.Second * time.Duration(conf.MergeIntervalS))
	scheduleTicker := time.NewTicker(time.Second * time.Duration(conf.ScheduleIntervalS))

	for {
		select {
		case <-ctx.Done():
			return
		case <-helloTicker.C:
			run(ctx, "cron->Hello", status, func(ctx context.Context) error {
				log.Println("hello world")
				return nil
			})
		case <-mergeTicker.C:
			run(ctx, "cron->repository.Merge", status, repository.Instance(conf).MergeRemotely)
		case <-scheduleTicker.C:
			run(ctx, "cron->scheduler.Run", status, scheduler.Run)
		}
	}
}

func run(ctx context.Context, what string, status health.Status, cb func(ctx context.Context) error) {
	log.Println("CRON ticker initiating", what)
	if err := cb(ctx); err != nil {
		status.Set(health.Critical, what, err)
		return
	}
	status.Clear(what)
}
