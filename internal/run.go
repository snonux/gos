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
		switch {
		case err == nil:
			log.Println("Scheduling", path)
			// TODO: Implement action here to post it
		case errors.Is(err, schedule.ErrNothingToSchedule):
			log.Println("Nothing to be scheduled for", platform)
		case errors.Is(err, schedule.ErrNothingQueued):
			log.Println("Nothing queued for", platform)
		default:
			return err
		}
	}

	return nil
}
