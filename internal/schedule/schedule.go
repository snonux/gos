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
)

var (
	ErrNothingToSchedule = errors.New("nothing to schedule")
	ErrNothingQueued     = errors.New("nothing queued")
)

func Run(args config.Args, platform string) (entry.Entry, error) {
	dir := fmt.Sprintf("%s/db/platforms/%s", args.GosDir, strings.ToLower(platform))
	stats, err := newStats(dir, args.Lookback, args.Target)
	if err != nil {
		return entry.Zero, err
	}

	log.Println("For", platform, "stats:", stats)
	// Schedule random queued entry with "now" tag, ignoring the target hit stats.
	// TODO: Document .now. tag
	ent, err := selectRandomEntry(dir, "now")
	if err != nil && !errors.Is(err, oi.ErrNotFound) {
		// Unknown error
		return ent, nil
	}
	if err == nil {
		return ent, nil
	}

	if stats.targetHit(args.PauseDays) {
		return entry.Zero, ErrNothingToSchedule
	}

	// Schedule random qeued entry for platform. Find one with prio tag.
	ent, err = selectRandomEntry(dir, "prio")
	if errors.Is(err, oi.ErrNotFound) {
		// No entry with priority tag found, select another one.
		ent, err = selectRandomEntry(dir, "")
	}
	if err != nil {
		return entry.Zero, fmt.Errorf("%w: %w", ErrNothingQueued, err)
	}
	return ent, nil
}

// Select a random queed entry with a given tag. If the tag is the empty string,
// then select any random qeued entry.
func selectRandomEntry(dir, tag string) (entry.Entry, error) {
	return oi.ReadDirRandom(dir, func(file os.DirEntry) (entry.Entry, bool) {
		// Is there a ".TAG." in the file name?
		if tag != "" && !slices.Contains(strings.Split(file.Name(), "."), tag) {
			return entry.Zero, false
		}
		ent, err := entry.New(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Println(err)
			return entry.Zero, false
		}
		return ent, ent.State == entry.Queued
	})
}
