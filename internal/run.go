package internal

import (
	"context"
	"errors"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/platforms"
	"codeberg.org/snonux/gos/internal/prompt"
	"codeberg.org/snonux/gos/internal/queue"
	"codeberg.org/snonux/gos/internal/schedule"
)

func Run(ctx context.Context, args config.Args) error {
	if err := queue.Run(args); err != nil {
		if !softError(err) {
			return err
		}
		_, _ = colour.Infoln(err)
	}

	for platformStr, sizeLimit := range args.Platforms {
		platform, err := platforms.New(platformStr)
		if err != nil {
			return err
		}
		if err := runPlatform(ctx, args, platform, sizeLimit); err != nil {
			if softError(err) {
				_, _ = colour.Infoln(err)
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
		_, _ = colour.Infoln("Nothing to be scheduled for", platform)
		return nil
	case errors.Is(err, schedule.ErrNothingQueued):
		_, _ = colour.Infoln("Nothing queued for", platform)
		return nil
	case err != nil:
		return err
	}
	err = platform.Post(ctx, args, sizeLimit, en)
	if errors.Is(err, prompt.ErrRamdomOther) {
		return runPlatform(ctx, args, platform, sizeLimit)
	}
	return err
}

func softError(err error) bool {
	return errors.Is(err, prompt.ErrAborted) || errors.Is(err, prompt.ErrDeleted)
}
