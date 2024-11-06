package queue

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
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
		en, err := entry.New(filePath)
		if err != nil {
			return err
		}
		// Extract any inline tags, if any!
		if err := en.ExtractInlineTags(); err != nil {
			return err
		}
		if en.HasTag("ask") {
			// TODO: Handle inline tags
			if err := en.FileAction("Do you want to queue this content"); err != nil {
				return err
			}
		}
		destPath := fmt.Sprintf("%s/db/%s.%s.queued", args.GosDir, filepath.Base(en.Path), timestamp.Now())
		if args.DryRun {
			log.Println("Not queueing entry", en.Path, "to", destPath, "as dry-run mode enabled")
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
		for platform := range args.Platforms {
			excluded, err := en.PlatformExcluded(args, platform)
			if err != nil {
				return err
			}
			if excluded {
				log.Println("Not queueing entry", en, "to platform", platform, "as it is excluded")
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
		log.Printf("Trashing %s -> %s", en.Path, trashPath)
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
func queuePlatform(en entry.Entry, gosDir, platform string) error {
	destDir := filepath.Join(gosDir, "db/platforms", strings.ToLower(platform))
	destPath := filepath.Join(destDir, filepath.Base(en.Path))
	postedFile := fmt.Sprintf("%s.posted", strings.TrimSuffix(destPath, ".queued"))

	// Entry already posted platform?
	if oi.IsRegular(postedFile) {
		log.Println("Not re-queueing", destPath, "as", postedFile, "already exists")
		return nil
	}

	log.Println("Queuing", en.Path, "->", destPath)
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
			log.Println("Cleaning up", filePath)
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
