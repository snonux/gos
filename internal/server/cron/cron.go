package cron

import (
	"context"
	"log"
	"time"

	config "codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/server/repository"
)

func Start(ctx context.Context, conf config.ServerConfig) error {
	go func() {
		helloTicker := time.NewTicker(time.Hour)
		mergeTicker := time.NewTicker(time.Second * time.Duration(conf.CRONMergeIntervalS))

		for {
			select {
			case <-ctx.Done():
				return
			case <-helloTicker.C:
				log.Println("CRON hello ticker ticked")
			case <-mergeTicker.C:
				log.Println("CRON ticker initiating remote merge operation")
				if err := repository.Instance(conf).MergeRemotely(ctx); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return nil
}
