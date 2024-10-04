package internal

import (
	"context"
	"errors"
	"log"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/platforms/mastodon"
	"codeberg.org/snonux/gos/internal/queue"
	"codeberg.org/snonux/gos/internal/schedule"
)

func Run(ctx context.Context, args config.Args) error {
	if err := queue.Run(args); err != nil {
		return err
	}

	for _, platform := range args.Platforms {
		ent, err := schedule.Run(args, platform)
		switch {
		case errors.Is(err, schedule.ErrNothingToSchedule):
			log.Println("Nothing to be scheduled for", platform)
			return nil
		case errors.Is(err, schedule.ErrNothingQueued):
			log.Println("Nothing queued for", platform)
			return nil
		case err != nil:
			return err
		}

		log.Println("Scheduling", ent)
		switch strings.ToLower(platform) {
		case "mastodon":
			if args.DryRun {
				log.Println("Not posting", ent, "to", platform, "as dry-run enabled")
				continue
			}
			if err := mastodon.Post(ctx, args, ent); err != nil {
				return err
			}
			if err := ent.MarkPosted(); err != nil {
				return err
			}
			log.Println("Posted", ent, "to", platform)
		default:
			log.Println("WARNING: Platform", platform, "not yet implemented")
		}
	}

	return nil
}
