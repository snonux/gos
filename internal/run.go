package internal

import (
	"context"
	"errors"
	"log"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/queue"
	"codeberg.org/snonux/gos/internal/schedule"
)

func Run(ctx context.Context, args config.Args) error {
	if err := queue.Run(args); err != nil {
		return err
	}

	for _, platform := range args.Platforms {
		path, err := schedule.Run(args, platform)
		switch err {
		case nil:
			log.Println("Scheduling", path)
			// TODO: Implement action here to post it
		case schedule.ErrNothingToSchedule:
			log.Println("Nothing to be scheduled for", platform)
		case schedule.ErrNothingQueued
			log.Println("Nothing queued for", platform)
		default:
			return err
		}
	}

	return nil
}
