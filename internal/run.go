package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/platforms"
	"codeberg.org/snonux/gos/internal/prompt"
	"codeberg.org/snonux/gos/internal/queue"
	"codeberg.org/snonux/gos/internal/schedule"
	"codeberg.org/snonux/gos/internal/summary"
)

func run(ctx context.Context, args config.Args) error {
	if len(args.GeminiSummaryFor) > 0 {
		return summary.Run(ctx, args)
	}
	now := time.Now().Unix()
	printLogo()

	// Check if posting is paused
	paused, err := args.Config.IsPaused()
	if err != nil {
		return fmt.Errorf("error checking pause status: %w", err)
	}
	if paused {
		colour.Infoln("Posting is paused until", args.Config.PauseEnd, "- skipping all posts")
		return nil
	}

	if args.ComposeMode {
		entryPath := fmt.Sprintf("%s/%d.ask.txt", args.GosDir, now)
		if err := prompt.EditFile(entryPath); err != nil {
			return err
		}
	}

	if err := queue.Run(args); err != nil {
		if !softError(err) {
			return err
		}
		colour.Infoln(err)
	}

	sinceLastRun := time.Duration(now-args.Config.LastRunEpoch) * time.Second
	if sinceLastRun < args.RunInterval {
		colour.Infoln("Run interval of", args.RunInterval, "with", sinceLastRun, "not yet reached. Not posting anything!")
		return nil
	}

	for platformStr, sizeLimit := range args.Platforms {
		platform, err := platforms.New(platformStr)
		if err != nil {
			return err
		}
		if err := runPlatform(ctx, args, platform, sizeLimit); err != nil {
			if softError(err) {
				colour.Infoln(err)
				continue
			}
			return err
		}
	}

	args.Config.LastRunEpoch = now
	return args.Config.WriteToDisk(args.ConfigPath)
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

	if args.ComposeMode {
		colour.Infoln("Not posting any entry in compose mode!")
		return nil
	}

	err = platform.Post(ctx, args, sizeLimit, en)
	if errors.Is(err, prompt.ErrRamdomOther) || errors.Is(err, prompt.ErrDeleted) {
		return runPlatform(ctx, args, platform, sizeLimit)
	}
	return err
}

func softError(err error) bool {
	return errors.Is(err, prompt.ErrAborted)
}
