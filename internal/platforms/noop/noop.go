package noop

import (
	"context"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/prompt"
)

// Psudo platform, not posting really anything.
func Post(ctx context.Context, args config.Args, sizeLimit int, en entry.Entry) error {
	content, _, err := en.ContentWithLimit(sizeLimit)
	if err != nil {
		return err
	}
	if args.DryRun {
		colour.Infoln("Not posting", en, "to Noop as dry-run enabled")
		return nil
	}
	if _, err = prompt.FileAction("Do you want to post this message to Noop?",
		content, en.Path, prompt.RandomOption); err != nil {
		return err
	}
	return nil
}
