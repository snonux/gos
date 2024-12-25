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

type maybeEntry struct {
	en  entry.Entry
	err error
}

func matchingEntries(args config.Args) iter.Seq[maybeEntry] {
	return func(yield func(maybeEntry) bool) {
		for _, dateStr := range args.SummaryFor {
			glob := filepath.Join(args.GosDir, "db/platforms/*/", fmt.Sprintf("*%s*-??????.posted", dateStr))
			paths, err := filepath.Glob(glob)
			if err != nil && !yield(maybeEntry{err: err}) {
				return
			}
			for _, path := range paths {
				en, err := entry.New(path)
				if !yield(maybeEntry{en, err}) {
					return
				}
			}
		}
	}
}

func deduppedEntries(args config.Args) ([]entry.Entry, error) {
	dedup := make(map[string]entry.Entry)

	for maybe := range matchingEntries(args) {
		if maybe.err != nil {
			return []entry.Entry{}, maybe.err
		}
		if en, ok := dedup[maybe.en.Name()]; ok {
			// If two conflicting entries (e.g. same post for mastodon and linkedin)
			// select the one which was modified latest.
			after, err := en.After(maybe.en)
			if err != nil {
				return []entry.Entry{}, err
			}
			if after {
				continue
			}
		}
		dedup[maybe.en.Name()] = maybe.en
	}

	entries := make([]entry.Entry, len(dedup))
	for _, val := range dedup {
		entries = append(entries, val)
	}
	return entries, nil
}
