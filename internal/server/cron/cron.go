package cron

import (
	"context"
	"fmt"
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
			log.Println("CRON hello ticker ticked")
		case <-mergeTicker.C:
			log.Println("CRON ticker initiating remote merge operation")
			if err := repository.Instance(conf).MergeRemotely(ctx); err != nil {
				status.Set(health.Critical, "cron", fmt.Errorf("unable to merge remote repository: %w", err))
			}
		case <-scheduleTicker.C:
			log.Println("CRON ticker initiating schedule operation")
			if err := scheduler.Run(ctx); err != nil {
				status.Set(health.Critical, "cron", fmt.Errorf("unable to schedule post(s): %w", err))
			}
		}
	}
}
