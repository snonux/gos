package schedule

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"codeberg.org/snonux/gos/internal/config"
	"codeberg.org/snonux/gos/internal/entry"
	"codeberg.org/snonux/gos/internal/oi"
	"codeberg.org/snonux/gos/internal/prompt"
)

var (
	ErrNothingToSchedule = errors.New("nothing to schedule")
	ErrNothingQueued     = errors.New("nothing queued")
)

// TODO: Schedule more than N entries per week when the backlog of queued items is large enough.
func Run(args config.Args, platform string) (entry.Entry, error) {
	dir := fmt.Sprintf("%s/db/platforms/%s", args.GosDir, strings.ToLower(platform))
	stats, err := newStats(dir, args.Lookback, args.Target)
	if err != nil {
		return entry.Zero, err
	}
	stats.Render(platform)

	if stats.queued < args.MinQueued {
		_ = prompt.Acknowledge(
			fmt.Sprintf("There are only %d messages queued for %s - time to fill it up!",
				stats.queued, platform),
		)
	}

	en, err := selectEntry(dir)
	if err != nil && !errors.Is(err, oi.ErrNotFound) {
		return en, nil // Unknown error
	}
	if !en.HasTag("now") && stats.targetHit(args.PauseDays) {
		return entry.Zero, ErrNothingToSchedule
	}
	return en, nil
}

/**
 * Select a random entry, but in this order:
 * 1. Any antry with the now tag
 * 2. Any entry with the prio tag
 * 3. Any entry with the soon tag
 * 4. Any other entry
 */
func selectEntry(dir string) (en entry.Entry, err error) {
	tagsToTry := []string{"now", "prio", "soon", ""}
	for _, tag := range tagsToTry {
		if en, err = selectRandomEntry(dir, tag); err == nil {
			return
		}
		if !errors.Is(err, oi.ErrNotFound) {
			return
		}

	}
	return
}

// Select a random queed entry with a given tag. If the tag is the empty string,
// then select any random qeued entry.
func selectRandomEntry(dir, tag string) (entry.Entry, error) {
	return oi.ReadDirRandom(dir, func(file os.DirEntry) (entry.Entry, bool) {
		// Is there a ".TAG." in the file name?
		if tag != "" && !slices.Contains(strings.Split(file.Name(), "."), tag) {
			return entry.Zero, false
		}
		en, err := entry.New(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Println(err)
			return entry.Zero, false
		}
		return en, en.State == entry.Queued
	})
}
