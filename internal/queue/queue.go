package queue

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/prompt"
	"codeberg.org/snonux/gos/internal/timestamp"
)

// Strictly, we only operate on .txt files, but we also accept .md as Obsidian creates only .md files.
var validExtensions = []string{".txt", ".md"}

// TODO: Red alert when there are no messages to schedule, or less than N
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
	ch, err := oi.ReadDirCh(args.GosDir, func(file os.DirEntry) (string, bool) {
		filePath := filepath.Join(args.GosDir, file.Name())
		return filePath, slices.Contains(validExtensions, filepath.Ext(file.Name())) &&
			file.Type().IsRegular()
	})
	if err != nil {
		return err
	}

	for filePath := range ch {
		ent, err := entry.New(filePath)
		if err != nil {
			return err
		}
		if ent.HasTag("ask") {
			content, _, err := ent.Content()
			if err != nil {
				return err
			}
			// TODO Refactor
			err = prompt.DoYouWantThis("Do you want to queue this content", content)
			switch {
			case errors.Is(err, prompt.ErrEditContent):
				err = ent.Edit()
			case errors.Is(err, prompt.ErrDeleteFile):
				if err = ent.Remove(); err == nil {
					continue
				}
			}
			if err != nil {
				return err
			}
		}
		destPath := fmt.Sprintf("%s/db/%s.%s.queued", args.GosDir, filepath.Base(filePath), timestamp.Now())
		if args.DryRun {
			log.Println("Not queueing entry", filePath, "to", destPath, "as dry-run mode enabled")
			continue
		}
		if err := oi.Rename(filePath, destPath); err != nil {
			return err
		}
	}

	return nil
}

// Queue all ./db/queued/*.txt.STAMP.queued into ./db/platforms/PLATFORM/*.txt.STAMP.queued
// for each PLATFORM
func queuePlatforms(args config.Args) error {
	dbDir := filepath.Join(args.GosDir, "db")
	ch, err := oi.ReadDirCh(dbDir, func(file os.DirEntry) (string, bool) {
		filePath := filepath.Join(dbDir, file.Name())
		return filePath, strings.HasSuffix(file.Name(), ".queued")
	})
	if err != nil {
		return err
	}

	trashDir := filepath.Join(args.GosDir, "db", "trashbin")
	for filePath := range ch {
		for platform := range args.Platforms {
			excluded, err := excludedByTags(args, filePath, platform)
			if err != nil {
				return err
			}
			if excluded {
				log.Println("Not queueing entry", filePath, "to platform", platform, "as it is excluded")
				continue
			}
			if err := queuePlatform(filePath, args.GosDir, platform); err != nil {
				return err
			}
		}
		if args.DryRun {
			continue
		}

		// Keep queued items in trash for a while.
		trashPath := filepath.Join(trashDir, strings.TrimSuffix(filepath.Base(filePath), ".queued")+".trash")
		log.Printf("Trashing %s -> %s", filePath, trashPath)
		if err := oi.EnsureParentDir(trashPath); err != nil {
			return err
		}
		if err := os.Rename(filePath, trashPath); err != nil {
			return err
		}
	}

	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	return deleteFiles(trashDir, ".trash", sixMonthsAgo)
}

// Queue ./db/queued/*.txt.STAMP.queued to ./db/platforms/PLATFORM/*.txt.STAMP.queued
func queuePlatform(entryPath, gosDir, platform string) error {
	destDir := filepath.Join(gosDir, "db/platforms", strings.ToLower(platform))
	destPath := filepath.Join(destDir, filepath.Base(entryPath))
	postedFile := fmt.Sprintf("%s.posted", strings.TrimSuffix(destPath, ".queued"))

	// Entry already posted platform?
	if oi.IsRegular(postedFile) {
		log.Println("Not re-queueing", destPath, "as", postedFile, "already exists")
		return nil
	}

	log.Println("Queuing", entryPath, "->", destPath)

	return oi.CopyFile(entryPath, destPath)
}

func deleteFiles(path, suffix string, olderThan time.Time) error {
	ch, err := oi.ReadDirCh(path, func(file os.DirEntry) (string, bool) {
		filePath := filepath.Join(path, file.Name())
		return filePath, strings.HasSuffix(filePath, suffix) && file.Type().IsRegular()
	})
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
