package internal

import (
	"context"
	"errors"
	"log"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/platforms/linkedin"
	"codeberg.org/snonux/gos/internal/platforms/mastodon"
	"codeberg.org/snonux/gos/internal/prompt"
	"codeberg.org/snonux/gos/internal/queue"
	"codeberg.org/snonux/gos/internal/schedule"
)

func Run(ctx context.Context, args config.Args) error {
	if err := queue.Run(args); err != nil {
		return err
	}

	for _, platform := range args.Platforms {
		if err := runPlatform(ctx, args, platform); err != nil {
			if errors.Is(err, prompt.ErrAborted) {
				log.Println("Aborted posting to", platform)
				continue
			}
			return err
		}
	}

	return nil
}

func runPlatform(ctx context.Context, args config.Args, platform string) error {
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
	var postCB func(context.Context, config.Args, entry.Entry) error
	switch strings.ToLower(platform) {
	case "mastodon":
		postCB = mastodon.Post
	case "linkedin":
		postCB = linkedin.Post
	default:
		log.Fatal("Platform", platform, "(not yet) implemented")
	}

	if err := postCB(ctx, args, ent); err != nil {
		return err
	}

	log.Println("Posted", ent, "to", platform)
	return ent.MarkPosted()
}
