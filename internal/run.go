package internal

import (
	"context"
	"errors"
	"log"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/platforms"
	"codeberg.org/snonux/gos/internal/platforms/linkedin"
	"codeberg.org/snonux/gos/internal/platforms/mastodon"
	"codeberg.org/snonux/gos/internal/prompt"
	"codeberg.org/snonux/gos/internal/queue"
	"codeberg.org/snonux/gos/internal/schedule"
)

func Run(ctx context.Context, args config.Args) error {
	if err := queue.Run(args); err != nil {
		if !softError(err) {
			return err
		}
		colour.Infoln(err)
	}

	for platform, sizeLimit := range args.Platforms {
		if err := runPlatform(ctx, args, platform, sizeLimit); err != nil {
			if softError(err) {
				colour.Infoln(err)
				continue
			}
			return err
		}
	}

	return nil
}

func runPlatform(ctx context.Context, args config.Args, platform platforms.Platform, sizeLimit int) error {
	en, err := schedule.Run(args, platform)
	switch {
	case errors.Is(err, schedule.ErrNothingToSchedule):
		colour.Infoln("Nothing to be scheduled for", platform)
		return nil
	case errors.Is(err, schedule.ErrNothingQueued):
		colour.Infoln("Nothing queued for", platform)
		return nil
	case err != nil:
		return err
	}

	colour.Infoln("Posting", en)
	var postCB func(context.Context, config.Args, int, entry.Entry) error
	switch platform.String() {
	case "mastodon":
		postCB = mastodon.Post
	case "linkedin":
		postCB = linkedin.Post
	default:
		log.Fatal("Platform", platform, "(not yet) implemented")
	}

	if err := postCB(ctx, args, sizeLimit, en); err != nil {
		return err
	}
	if err := en.MarkPosted(); err != nil {
		return err
	}

	colour.Successf("Successfully posted message to %s", platform)
	return nil
}

func softError(err error) bool {
	return errors.Is(err, prompt.ErrAborted) || errors.Is(err, prompt.ErrDeleted)
}
