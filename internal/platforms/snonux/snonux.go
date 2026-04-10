package snonux

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/prompt"
	"codeberg.org/snonux/gos/internal/timestamp"
)

// Post writes the entry body as plain text under <gosDir>/snonux/inbox (filename .md).
func Post(ctx context.Context, args config.Args, sizeLimit int, en entry.Entry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	content, _, err := en.ContentWithLimit(sizeLimit)
	if err != nil {
		return err
	}
	if args.DryRun {
		colour.Infoln("Not writing", en, "to Snonux inbox as dry-run enabled")
		return nil
	}
	if _, err = prompt.FileAction("Do you want to post this message to Snonux (save to inbox)?",
		content, en.Path, prompt.RandomOption); err != nil {
		return err
	}
	path, err := inboxMarkdownPath(args.GosDir, en)
	if err != nil {
		return err
	}
	if err := oi.EnsureParentDir(path); err != nil {
		return err
	}
	out := strings.TrimSuffix(content, "\n") + "\n"
	if err := oi.WriteFile(path, out); err != nil {
		return err
	}
	colour.Infoln("Wrote Snonux inbox file", path)
	return nil
}

func inboxMarkdownPath(gosDir string, en entry.Entry) (string, error) {
	inbox := filepath.Join(gosDir, "snonux", "inbox")
	base := strings.TrimSuffix(filepath.Base(en.Path), ".queued")
	base = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_', r == '.':
			return r
		default:
			return '-'
		}
	}, base)
	base = strings.Trim(base, "-.")
	if base == "" {
		base = "post"
	}
	if len(base) > 120 {
		base = base[:120]
	}
	name := fmt.Sprintf("%s-%09d-%s.md", timestamp.Now(), time.Now().Nanosecond(), base)
	return filepath.Join(inbox, name), nil
}
