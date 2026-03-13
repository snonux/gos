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
	if err := checkPauseStatus(args); err != nil {
		return err
	}

	// Handle compose mode
	if args.ComposeMode {
		if err := handleComposeMode(args); err != nil {
			return err
		}
	}

	// Run queue operations
	if err := runQueueOperations(args); err != nil {
		return err
	}

	// Check run interval
	if err := checkRunInterval(args); err != nil {
		return err
	}

	// Post to platforms
	if err := postToPlatforms(ctx, args); err != nil {
		return err
	}

	// Update last run time
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

func checkPauseStatus(args config.Args) error {
	// Check if posting is paused
	paused, err := args.Config.IsPaused()
	if err != nil {
		return fmt.Errorf("error checking pause status: %w", err)
	}
	if paused {
		colour.Infoln("Posting is paused until", args.Config.PauseEnd, "- skipping all posts")
		return nil
	}
	return nil
}

func handleComposeMode(args config.Args) error {
	entryPath := fmt.Sprintf("%s/%d.ask.txt", args.GosDir, time.Now().Unix())
	if err := prompt.EditFile(entryPath); err != nil {
		return err
	}
	return nil
}

func runQueueOperations(args config.Args) error {
	if err := queue.Run(args); err != nil {
		if !softError(err) {
			return err
		}
		colour.Infoln(err)
	}
	return nil
}

func checkRunInterval(args config.Args) error {
	now := time.Now().Unix()
	sinceLastRun := time.Duration(now-args.Config.LastRunEpoch) * time.Second
	if sinceLastRun < args.RunInterval {
		colour.Infoln("Run interval of", args.RunInterval, "with", sinceLastRun, "not yet reached. Not posting anything!")
		return nil
	}
	return nil
}

func postToPlatforms(ctx context.Context, args config.Args) error {
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
	return nil
}

func softError(err error) bool {
	return errors.Is(err, prompt.ErrAborted)
}
