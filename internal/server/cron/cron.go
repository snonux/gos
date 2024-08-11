package cron

import (
	"context"
	"log"
	"time"

	config "codeberg.org/snonux/gos/internal/config/server"
	"codeberg.org/snonux/gos/internal/server/handler"
)

func Start(ctx context.Context, conf config.ServerConfig, hand handler.Handler) error {
	go func() {
		helloTicker := time.NewTicker(10 * time.Second)
		mergeTicker := time.NewTicker(time.Duration(conf.CRONMergeIntervalS) * time.Second)

		for {
			select {
			case <-ctx.Done():
				return
			case <-helloTicker.C:
				log.Println("Hello ticker ticked")
			case <-mergeTicker.C:
				log.Println("CRON merge ticker ticked")
			}
		}
	}()

	return nil
}
