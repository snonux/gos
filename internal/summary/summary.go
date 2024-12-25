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
	entries := make(map[string]entry.Entry)

	for _, dateStr := range args.SummaryFor {
		glob := filepath.Join(args.GosDir, "db/platforms/*/", fmt.Sprintf("*%s*-??????.posted", dateStr))
		paths, err := filepath.Glob(glob)
		if err != nil {
			return err
		}

		for _, path := range paths {
			en, err := entry.New(path)
			if err != nil {
				return err
			}
			// TODO: Dedup here by name
			entries[en.Name()] = en
		}
	}
	return nil
}
