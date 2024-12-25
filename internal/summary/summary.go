package summary

import (
	"context"
	"fmt"
	"iter"
	"path/filepath"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
)

func Run(ctx context.Context, args config.Args) error {
	colour.Infoln("Generating summary for", args.SummaryFor)

	entries, err := deduppedEntries(args)
	if err != nil {
		return err
	}
	fmt.Println(entries)

	return nil
}

func matchingEntries(args config.Args) iter.Seq2[entry.Entry, error] {
	return func(yield func(entry.Entry, error) bool) {
		for _, dateStr := range args.SummaryFor {
			glob := filepath.Join(args.GosDir, "db/platforms/*/", fmt.Sprintf("*%s*-??????.posted", dateStr))
			paths, err := filepath.Glob(glob)
			if err != nil && !yield(entry.Zero, err) {
				return
			}
			for _, path := range paths {
				en, err := entry.New(path)
				if !yield(en, err) {
					return
				}
			}
		}
	}
}

func deduppedEntries(args config.Args) ([]entry.Entry, error) {
	dedup := make(map[string]entry.Entry)

	for en, err := range matchingEntries(args) {
		if err != nil {
			return entry.Zeroes, err
		}
		if other, ok := dedup[en.Name()]; ok {
			// If two conflicting entries (e.g. same post for mastodon and linkedin)
			// select the one which was modified latest.
			after, err := other.After(en)
			if err != nil {
				return entry.Zeroes, err
			}
			if after {
				continue
			}
		}
		dedup[en.Name()] = en
	}

	entries := make([]entry.Entry, len(dedup))
	for _, val := range dedup {
		entries = append(entries, val)
	}
	return entries, nil
}
