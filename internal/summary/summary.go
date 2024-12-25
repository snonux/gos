package summary

import (
	"context"
	"fmt"
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

func deduppedEntries(args config.Args) ([]entry.Entry, error) {
	dedup := make(map[string]struct{})
	var entries []entry.Entry

	for _, dateStr := range args.SummaryFor {
		glob := filepath.Join(args.GosDir, "db/platforms/*/", fmt.Sprintf("*%s*-??????.posted", dateStr))
		paths, err := filepath.Glob(glob)
		if err != nil {
			return entries, err
		}

		for _, path := range paths {
			en, err := entry.New(path)
			if err != nil {
				return entries, err
			}
			if _, ok := dedup[en.Name()]; ok {
				continue
			}
			dedup[en.Name()] = struct{}{}
			entries = append(entries, en)
		}
	}

	return entries, nil
}
