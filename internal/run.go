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
	"github.com/fatih/color"
)

func Run(ctx context.Context, args config.Args) error {
	if err := queue.Run(args); err != nil {
		if !softError(err) {
			return err
		}
		log.Println(err)
	}

	for platform, sizeLimit := range args.Platforms {
		if err := runPlatform(ctx, args, platform, sizeLimit); err != nil {
			if softError(err) {
				log.Println(err)
				continue
			}
			return err
		}
	}

	return nil
}

func runPlatform(ctx context.Context, args config.Args, platform string, sizeLimit int) error {
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

	log.Println("Posting", ent)
	var postCB func(context.Context, config.Args, int, entry.Entry) error
	switch strings.ToLower(platform) {
	case "mastodon":
		postCB = mastodon.Post
	case "linkedin":
		postCB = linkedin.Post
	default:
		log.Fatal("Platform", platform, "(not yet) implemented")
	}

	if err := postCB(ctx, args, sizeLimit, ent); err != nil {
		return err
	}
	if err := ent.MarkPosted(); err != nil {
		return err
	}

	// TODO: Put all color definitions into ints own package
	color.New(color.FgWhite, color.BgGreen).Println("Successfully posted message to ", platform)
	return nil
}

func softError(err error) bool {
	return errors.Is(err, prompt.ErrAborted) || errors.Is(err, prompt.ErrDeleted)
}
