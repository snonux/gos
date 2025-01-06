package queue

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/platforms"
	"codeberg.org/snonux/gos/internal/tags"
	"codeberg.org/snonux/gos/internal/timestamp"
)

// Strictly, we only operate on .txt files, but we also accept .md as Obsidian creates only .md files.
var validExtensions = []string{".txt", ".md"}

func Run(args config.Args) error {
	if err := queueEntries(args); err != nil {
		return err
	}
	if err := queuePlatforms(args); err != nil {
		return err
	}
	return nil
}

// Queue all *.txt into ./db/*.txt.STAMP.queued
func queueEntries(args config.Args) error {
	ch, err := oi.ReadDirCh(args.GosDir, find(args.GosDir, validExtensions...))
	if err != nil {
		return err
	}

	for filePath := range ch {
		if filePath, err = tags.InlineExtract(filePath); err != nil {
			return err
		}
		en, err := entry.New(filePath)
		if err != nil {
			return err
		}

		hasHashtags, err := en.HasHashtags()
		if err != nil {
			return err
		}
		if !hasHashtags {
			colour.Warnln("The following entry has got no hashtags:")
		}
		if !hasHashtags || en.HasTag("ask") {
			if err := en.FileAction("Do you want to queue this"); err != nil {
				return err
			}
		}

		destPath := fmt.Sprintf("%s/db/%s.%s.queued", args.GosDir, filepath.Base(en.Path), timestamp.Now())
		if args.DryRun {
			colour.Infoln("Not queueing entry", en.Path, "to", destPath, "as dry-run mode enabled")
			continue
		}
		if err := oi.Rename(en.Path, destPath); err != nil {
			return err
		}
	}

	return nil
}

// Queue all ./db/queued/*.txt.STAMP.queued into ./db/platforms/PLATFORM/*.txt.STAMP.queued
// for each PLATFORM
func queuePlatforms(args config.Args) error {
	dbDir := filepath.Join(args.GosDir, "db")
	ch, err := oi.ReadDirCh(dbDir, find(dbDir, ".queued"))
	if err != nil {
		return err
	}

	trashDir := filepath.Join(args.GosDir, "db", "trashbin")
	for filePath := range ch {
		en, err := entry.New(filePath)
		if err != nil {
			return err
		}
		for platformStr := range args.Platforms {
			platform, err := platforms.New(platformStr)
			if err != nil {
				return err
			}
			// func NewShare(args config.Args, tags map[string]struct{}) (Share, error) {
			share, err := tags.NewShare(args, en.Tags)
			if err != nil {
				return err
			}
			if share.Excluded(platform.String()) {
				colour.Infoln("Not queueing entry", en, "to platform", platform, "as it is excluded")
				continue
			}
			if err := queuePlatform(en, args.GosDir, platform); err != nil {
				return err
			}
		}
		if args.DryRun {
			continue
		}

		// Keep queued items in trash for a while.
		trashPath := filepath.Join(trashDir, strings.TrimSuffix(filepath.Base(en.Path), ".queued")+".trash")
		colour.Infof("Trashing %s -> %s", en.Path, trashPath)
		if err := oi.EnsureParentDir(trashPath); err != nil {
			return err
		}
		if err := os.Rename(en.Path, trashPath); err != nil {
			return err
		}
	}

	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	return deleteFiles(trashDir, ".trash", sixMonthsAgo)
}

// Queue ./db/queued/*.txt.STAMP.queued to ./db/platforms/PLATFORM/*.txt.STAMP.queued
func queuePlatform(en entry.Entry, gosDir string, platform platforms.Platform) error {
	destDir := filepath.Join(gosDir, "db/platforms", platform.String())
	destPath := filepath.Join(destDir, filepath.Base(en.Path))
	postedFile := fmt.Sprintf("%s.posted", strings.TrimSuffix(destPath, ".queued"))

	// Entry already posted platform?
	if oi.IsRegular(postedFile) {
		colour.Infoln("Not re-queueing", destPath, "as", postedFile, "already exists")
		return nil
	}

	colour.Infoln("Queuing", en.Path, "->", destPath)
	return oi.CopyFile(en.Path, destPath)
}

func deleteFiles(path, suffix string, olderThan time.Time) error {
	ch, err := oi.ReadDirCh(path, find(path, suffix))
	if err != nil {
		return err
	}
	for filePath := range ch {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		if fileInfo.ModTime().Before(olderThan) {
			colour.Infoln("Cleaning up", filePath)
			err := os.Remove(filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func find(path string, suffixes ...string) func(os.DirEntry) (string, bool) {
	return func(file os.DirEntry) (string, bool) {
		filePath := filepath.Join(path, file.Name())
		for _, suffix := range suffixes {
			if strings.HasSuffix(file.Name(), suffix) {
				return filePath, true
			}
		}
		return filePath, false
	}
}
