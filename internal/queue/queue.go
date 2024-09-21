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
	"codeberg.org/snonux/gos/internal/oi"
)

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
	// Strictly, we only operate on .txt files, but we also accept .md as Obsidian creates only .md files.
	var validExtensions = []string{".txt", ".md"}

	ch, err := oi.ReadDirFilter(args.GosDir, func(entry os.DirEntry) bool {
		return slices.Contains(validExtensions, filepath.Ext(entry.Name())) &&
			entry.Type().IsRegular()
	})
	if err != err {
		return err
	}

	now := time.Now()
	for filePath := range ch {
		destPath := fmt.Sprintf("%s/db/%s.%s.queued", args.GosDir,
			filepath.Base(filePath), now.Format("20060102-150405"))
		if err := os.Rename(filePath, destPath); err != nil {
			return err
		}
	}

	return nil
}

// Queue all ./db/queued/*.txt.STAMP.queued into ./db/platforms/PLATFORM/*.txt.STAMP.queued
// for each PLATFORM
func queuePlatforms(args config.Args) error {
	dbDir := fmt.Sprintf("%s/db", args.GosDir)
	ch, err := oi.ReadDirFilter(dbDir, func(entry os.DirEntry) bool {
		return strings.HasSuffix(entry.Name(), ".queued")
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
		log.Println("Removing", filePath)
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return nil
}

// Queue ./db/queued/*.txt.STAMP.queued to ./db/platforms/PLATFORM/*.txt.STAMP.queued
func queuePlatform(entryPath, gosDir, platform string) error {
	destDir := fmt.Sprintf("%s/db/platforms/%s/", gosDir, strings.ToLower(platform))
	destPath := fmt.Sprintf("%s/%s", destDir, filepath.Base(entryPath))
	postedFile := fmt.Sprintf("%s.posted", strings.TrimSuffix(destPath, ".queued"))

	// Entry already posted platform?
	if oi.IsRegular(postedFile) {
		log.Println("Not re-queueing", destPath, "as", postedFile, "already exists")
		return nil
	}

	log.Println("Queuing", entryPath, "->", destPath)
	return oi.CopyFile(entryPath, destPath)
}
