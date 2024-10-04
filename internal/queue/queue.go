package queue

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/format"
	"codeberg.org/snonux/gos/internal/oi"
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
	ch, err := oi.ReadDirCh(args.GosDir, func(file os.DirEntry) (string, bool) {
		filePath := filepath.Join(args.GosDir, file.Name())
		return filePath, slices.Contains(validExtensions, filepath.Ext(file.Name())) &&
			file.Type().IsRegular()
	})
	if err != err {
		return err
	}

	now := time.Now()
	for filePath := range ch {
		destPath := fmt.Sprintf("%s/db/%s.%s.queued", args.GosDir,
			filepath.Base(filePath), now.Format(format.Time))
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
	if err != err {
		return err
	}

	for filePath := range ch {
		for _, platform := range args.Platforms {
			if err := queuePlatform(filePath, args.GosDir, platform); err != nil {
				return err
			}
		}
		if args.DryRun {
			continue
		}
		log.Println("Removing", filePath)
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return nil
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
